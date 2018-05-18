package main

import (
	"github.com/adlio/trello"
	"strings"
)

type Attachment struct{
	Name string `json:"name"`
	URL string `json:"url"`
}
type Story struct{
	Name string `json:"story_name"`
	StoryPoints string `json:"story_points"`
	Labels []string `json:"labels"`
	Checklists map[string][]string `json:"checklists"`
	CardUrl string `json:"card_url"`
	Attachments []Attachment `json:"attachments"`
	Description string `json:"description"`
}

func extractAttachments(trelloAttachments []*trello.Attachment) []Attachment{
	var attachments = make([]Attachment, 0, len(trelloAttachments))
	for _, a := range trelloAttachments{
		attachments = append(attachments, Attachment{a.Name, a.URL})
	}
	return attachments
}

func extractChecklists(trelloChecklists []*trello.Checklist) map[string][]string {
	checklist := make(map[string][]string)

	for _,tc := range trelloChecklists{
		checkItems := make([]string, 0, len(tc.CheckItems))
		for _, checkItem := range tc.CheckItems{
			checkItems = append(checkItems, checkItem.Name)
		}
		checklist[tc.Name] = checkItems
	}
	return checklist
}

func extractLabels(trelloLabels []*trello.Label)[]string  {
	labels := make([]string, 0, len(trelloLabels))
	for _, tl :=range trelloLabels{
		labels = append(labels, tl.Name)
	}
	return labels
}

func extractStoryPoints(name string) string  {
	storyPoints := strings.SplitAfter(name, "]")[0]
	storyPoints = strings.Replace(storyPoints, "[", "", 1)
	storyPoints = strings.Replace(storyPoints, "]", "", 1)
	return storyPoints
}
