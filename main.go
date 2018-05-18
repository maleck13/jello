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
	"time"
	json2 "encoding/json"
)

var (
	flagTrelloToken = flag.String("trello-token", "", "set your trello token")
	flagJiraHost    = flag.String("jira-host", "https://issues.jboss.org", "your jira installation host")
	flagJiraUser    = flag.String("jira-user", "", "your jira username")
	flagJiraPass    = flag.String("jira-pass", "", "your jira password")
	writeStories	= flag.Bool("out", false, "whether you want output stories to json file")
	jiraClient      *jira.Client
)

func createJiraClient() {
	authURL := fmt.Sprintf("%s/rest/auth/1/session", *flagJiraHost)
	if *flagJiraUser == "" {
		*flagJiraUser = os.Getenv("JIRA_USER")
	}
	if *flagJiraPass == "" {
		*flagJiraPass = os.Getenv("JIRA_PASS")
	}

 	tp := jira.CookieAuthTransport{
		Username: *flagJiraUser,
		Password: *flagJiraPass,
		AuthURL:  authURL,
	}


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

func getStoriesFromTrello()[]Story{
	//system env vars in your
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

func main() {
	flag.Parse()
	createJiraClient()
	
	stories := getStoriesFromTrello()

	if *writeStories {
		storiesJson, _ := json2.MarshalIndent(stories, "", "  ")
		ioutil.WriteFile("stories_" + time.Now().Format(time.RFC822)+ ".json", storiesJson, 0666)
	}
}