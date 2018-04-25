package main

import (
	"github.com/andygrunwald/go-jira"
	"strings"
	"fmt"
)

const TEAM_LABEL = "team-mcp-core"

func buildEpic(story Story, sampleEpic *jira.Issue) jira.Issue{
	return jira.Issue{
		Fields: &jira.IssueFields{
			Reporter:    &jira.User{Name: "jello"},
			Type: sampleEpic.Fields.Type,
			Project: sampleEpic.Fields.Project,
			Description: buildDescription(story),
			Summary: story.Name,
			Labels: []string{TEAM_LABEL},
		},
	}
}

func buildDescription(story Story) string {
	var desc []string
	desc = append(desc, fmt.Sprintf("%s\n\n", story.Name))
	desc = append(desc, fmt.Sprintf("%s\n\n\n", story.Description))
	ac := "Acceptance Criteria"
	desc = append(desc, fmt.Sprintf("%s:\n", ac))

	for _, checkItem := range story.Checklists[ac] {
		desc = append(desc, fmt.Sprintf("- %s\n", checkItem))
	}

	attachments := "Attachments"
	desc = append(desc, fmt.Sprintf("\n\n%s:\n", attachments))

	for _, att := range story.Attachments {
		desc = append(desc, fmt.Sprintf("- %s\n", att.Name))
		desc = append(desc, fmt.Sprintf("  %s\n", att.URL))
	}


	labels := "Labels"
	desc = append(desc, fmt.Sprintf("\n\n%s:\n", labels))

	for _, label := range story.Labels {
		desc = append(desc, fmt.Sprintf("- %s\n", label))
	}

	tasks := "Tasks"
	desc = append(desc, fmt.Sprintf("\n\n%s:\n", tasks))

	for _, checkItem := range story.Checklists[tasks] {
		desc = append(desc,
			fmt.Sprintf("- %s / labels: \"%s\" issueType:\"task\" \n", checkItem, TEAM_LABEL))
	}

	descStr := strings.Join(desc, " ")

	return descStr
}
