package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/adlio/trello"
	"github.com/andygrunwald/go-jira"
)

var (
	flagTrelloToken = flag.String("trello_token", "", "set your trello token")
	flagJiraPass    = flag.String("jira-pass", "", "your jira password")
	flagJiraUser    = flag.String("jira-user", "", "your jira username")
	flagJiraHost    = flag.String("jira-host", "https://issues.jboss.org", "your jira installation host")
	jiraClient      *jira.Client
)

func createJiraClient() {
	authURL := fmt.Sprintf("%s/rest/auth/1/session", *flagJiraHost)
	tp := jira.CookieAuthTransport{
		Username: *flagJiraUser,
		Password: *flagJiraPass,
		AuthURL:  authURL,
	}

	fmt.Println(tp)

	c, err := jira.NewClient(tp.Client(), *flagJiraHost)
	if err != nil {
		panic(err)
	}
	if _, res, err := c.User.Get(*flagJiraUser); err != nil {
		b, _ := ioutil.ReadAll(res.Body)
		log.Fatal(string(b), err)
	}

	jiraClient = c
}

func AddToJira() {
	if jiraClient == nil {
		log.Fatal("no jira client")
		return
	}

}

func main() {
	flag.Parse()
	createJiraClient()

	//get sample epic
	sampleEpic, _, _ := jiraClient.Issue.Get("AEROGEAR-2363", nil)

	stories := getStoriesFromTrello()

	for _, story := range stories{
		epic := buildEpic(story, sampleEpic)
		fmt.Println(epic)
		//jiraClient.Issue.Create(&epic)
	}

}

func getStoriesFromTrello()[]Story{
	appKey := os.Getenv("TRELLO_APP_KEY")
	token := os.Getenv("TRELLO_TOKEN")

	client := trello.NewClient(appKey, token)
	board, err := client.GetBoard("VrnvGm7P", trello.Defaults())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	stories := make([]Story, 0)
	for _, l := range lists {
		if strings.TrimSpace(strings.ToLower(l.Name)) == "backlog" {
			fmt.Println("found backlog", l)
			args := trello.Defaults()
			args["checklists"] = "all"
			args["attachments"] = "true"
			cards, _ := l.GetCards(args)
			stories = processCards(cards)
		}
	}
	return stories
}

func processCards(cards []*trello.Card) []Story{
	stories := make([]Story, 0, len(cards))
	for _, c := range cards {
		name := c.Name
		storyPoints := extractStoryPoints(name)
		labels := extractLabels(c.Labels)
		checklists := extractChecklists(c.Checklists)
		cardUrl := c.ShortUrl
		attachments := extractAttachments(c.Attachments)

		stories = append(stories, Story{
			name,
			storyPoints,
			labels,
			checklists,
			cardUrl,
			attachments,
			c.Desc,
		})
	}
	return stories
}
