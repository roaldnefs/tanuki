// Copyright © 2018 Roald Nefs <info@roaldnefs.com>

package cmd

import (
	"fmt"
	"log"
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var repository string

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit members, branches, hooks, deploy keys etc.",
	Long: `Audit members, branches, hooks, deploy keys etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient()
		handleAudit(client, repository)
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// auditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// auditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
		projectMembers := []string{}
		for _, m := range members {
			projectMembers = append(projectMembers, fmt.Sprintf("\t\t\t%s", m.Username))
		}
		output += fmt.Sprintf("\tMembers (%d):\n", len(members))
		output += fmt.Sprintf("\t\tRole (%d):\n%s\n", len(projectMembers), strings.Join(projectMembers, "\n"))
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
		m, resp, err := client.ProjectMembers.ListAllProjectMembers(pid, opt)
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

	return nil, errors.New(fmt.Sprintf("Project %s not found!", repository))
}
