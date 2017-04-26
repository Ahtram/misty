package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"time"

	"github.com/julienschmidt/httprouter"
)

// GitHubHook represents a GitHub hook we are listening.
type GitHubHook struct {
	GitHubHookEndPoint string
	GitHubHookPort     string
	MistyRef           *Misty
}

type gitHubPushHook struct {
	Commits    []gitHubCommit   `json:"commits"`
	Repository gitHubRepository `json:"repository"`
}

type gitHubRepository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type gitHubCommit struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	TimeStamp string       `json:"timestamp"`
	URL       string       `json:"url"`
	Author    gitHubAuthor `json:"author"`
}

type gitHubAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// StartGitHubHook start the gitHub router.
func (gitHubHook *GitHubHook) StartGitHubHook() {
	router := httprouter.New()
	router.POST(gitHubHook.GitHubHookEndPoint, gitHubHook.receiveGiHubDelivery)
	log.Fatal(http.ListenAndServe(":"+gitHubHook.GitHubHookPort, router))
}

// receiveDelivery gets all build event from Unity Cloud.
func (gitHubHook *GitHubHook) receiveGiHubDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the body content.
	requestByteArray, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(Red("[receiveGiHubDelivery Error] ") + err.Error())
		return
	}

	// Get the json content
	gitHubPushHook := gitHubPushHook{}
	err = json.Unmarshal(requestByteArray, &gitHubPushHook)
	if err != nil {
		fmt.Println(Red("[receiveGiHubDelivery Error] ") + err.Error())
		return
	}

	//Broadcast the download info to channel.
	informMessage := ":package: [GitHub] [" + gitHubPushHook.Repository.Name + "] " + gitHubHook.MistyRef.Line("newRevision", 0) + "\n"
	informMessage += "```Markdown\n"
	for _, commit := range gitHubPushHook.Commits {
		t, _ := time.Parse(time.RFC3339, commit.TimeStamp)
		informMessage += "#[" + t.Format("2006-01-02 15:04:05") + "] [" + commit.Author.Name + "]\n"
		informMessage += "    " + commit.Message + "\n"
	}
	informMessage += "```"
	gitHubHook.MistyRef.broadcastMessage(informMessage)
}
