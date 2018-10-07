// Copyright © 2018 Roald Nefs <info@roaldnefs.com>

package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// List of patched available access levels.
// Should be removed once github.com/xanzy/go-gitlab/pull/486 is merged.
const (
	NoPermissions         gitlab.AccessLevelValue = 0
	GuestPermissions      gitlab.AccessLevelValue = 10
	ReporterPermissions   gitlab.AccessLevelValue = 20
	DeveloperPermissions  gitlab.AccessLevelValue = 30
	MaintainerPermissions gitlab.AccessLevelValue = 40
	OwnerPermissions      gitlab.AccessLevelValue = 50
)

var repository string

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit members, branches, hooks, deploy keys etc.",
	Long:  `Audit members, branches, hooks, deploy keys etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		handleAudit(client, repository)
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)

	// Flags and configuration settings.
	auditCmd.Flags().StringVarP(&repository, "repository", "r", "", "specific repo (e.g. 'roaldnefs/tanuki')")
	auditCmd.MarkFlagRequired("repository")
}

// handleAudit will return nil error if the user does not habe access to
// something.
func handleAudit(client *gitlab.Client, repository string) error {
	project, err := getProject(client, repository)
	if err != nil {
		log.Fatal(err)
	}

	members, err := getAllProjectMembers(client, project.ID)
	if err != nil {
		log.Fatal(err)
	}

	output := fmt.Sprintf("%s -> \n", project.NameWithNamespace)

	if len(members) >= 1 {
		owners := []string{}
		maintainers := []string{}
		developers := []string{}
		reporters := []string{}
		guests := []string{}

		for _, m := range members {
			switch m.AccessLevel {
			case OwnerPermissions:
				owners = append(owners, fmt.Sprintf("\t\t\t%s", m.Username))
			case MaintainerPermissions:
				maintainers = append(maintainers, fmt.Sprintf("\t\t\t%s", m.Username))
			case DeveloperPermissions:
				developers = append(developers, fmt.Sprintf("\t\t\t%s", m.Username))
			case ReporterPermissions:
				reporters = append(reporters, fmt.Sprintf("\t\t\t%s", m.Username))
			case GuestPermissions:
				guests = append(guests, fmt.Sprintf("\t\t\t%s", m.Username))
			}
		}

		output += fmt.Sprintf("\tMembers (%d):\n", len(members))
		output += fmt.Sprintf("\t\tOwner (%d):\n%s\n", len(owners), strings.Join(owners, "\n"))
		output += fmt.Sprintf("\t\tMaintainer (%d):\n%s\n", len(maintainers), strings.Join(maintainers, "\n"))
		output += fmt.Sprintf("\t\tDeveloper (%d):\n%s\n", len(developers), strings.Join(developers, "\n"))
		output += fmt.Sprintf("\t\tReporter (%d):\n%s\n", len(reporters), strings.Join(reporters, "\n"))
		output += fmt.Sprintf("\t\tGuest (%d):\n%s\n", len(guests), strings.Join(guests, "\n"))
	}

	visibility := fmt.Sprintf("\tVisibility: %s", project.Visibility)
	output += visibility + "\n"

	mergeMethod := fmt.Sprintf("\tMerge Method: %s", project.MergeMethod)
	output += mergeMethod + "\n"

	fmt.Printf("%s--\n\n", output)

	return nil
}

func getAllProjectMembers(client *gitlab.Client, pid interface{}) ([]*gitlab.ProjectMember, error) {
	opt := &gitlab.ListProjectMembersOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	var members []*gitlab.ProjectMember

	for {
		// Get the current page with members.
		//
		// TODO Should be replaced by github.com/xanzy/go-gitlab/pull/485 the
		// following code once merged:
		// m, resp, err := client.ProjectMembers.ListAllProjectMembers(pid, opt)
		m, resp, err := ListAllProjectMembers(client, pid, opt)
		if err != nil {
			log.Fatal(err)
		}

		// Add newly found members to the list.
		members = append(members, m...)

		// Exit loop when we've seen all the pages.
		if opt.Page >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return members, nil
}

// getProject returns the GitLab project based on the repository name by looping
// over all the projects.
func getProject(client *gitlab.Client, repository string) (*gitlab.Project, error) {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	for {
		// Get the current page with projects.
		projects, resp, err := client.Projects.ListProjects(opt)
		if err != nil {
			log.Fatal(err)
		}

		// List all the projects we've found so far.
		for _, p := range projects {
			// Return the project if the PathWithNamespace equals the repository
			// string.
			if p.PathWithNamespace == repository {
				return p, nil
			}
		}

		// Exit loop when we've seen all the pages.
		if opt.Page >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return nil, errors.New(fmt.Errorf("project %s not found", repository))
}

// Helper function to accept and format both the project ID or name as project
// identifier for all API calls.
//
// TODO Should be removed once github.com/xanzy/go-gitlab/pull/485 is merged.
func parseID(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid ID type %#v, the ID must be an int or a string", id)
	}
}

// ListAllProjectMembers gets a list of a project's team members viewable by the
// authenticated user. Returns a list including inherited members through
// ancestor groups.
//
// TODO Should be removed once github.com/xanzy/go-gitlab/pull/485 is merged.
func ListAllProjectMembers(client *gitlab.Client, pid interface{}, opt *gitlab.ListProjectMembersOptions, options ...gitlab.OptionFunc) ([]*gitlab.ProjectMember, *gitlab.Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/members/all", url.QueryEscape(project))

	req, err := client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var pm []*gitlab.ProjectMember
	resp, err := client.Do(req, &pm)
	if err != nil {
		return nil, resp, err
	}

	return pm, resp, err
}
