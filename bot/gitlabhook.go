package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GitLabHook represents a GitLab hook we are listening.
type GitLabHook struct {
	GitLabHookEndPoint string
	GitLabHookPort     string
	MistyRef           *Misty
}

type gitLabPushHook struct {
	Project          gitLabProject  `json:"project"`
	Commits          []gitLabCommit `json:"commits"`
	TotalCommitCount int            `json:"total_commits_count"`
}

type gitLabProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type gitLabCommit struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	TimeStamp string       `json:"timestamp"`
	URL       string       `json:"url"`
	Author    gitLabAuthor `json:"author"`
}

type gitLabAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// StartGitLabHook start the gitlab router.
func (gitLabHook *GitLabHook) StartGitLabHook() {
	router := httprouter.New()
	router.POST(gitLabHook.GitLabHookEndPoint, gitLabHook.receiveGitLabDelivery)
	log.Fatal(http.ListenAndServe(":"+gitLabHook.GitLabHookPort, router))
}

// receiveDelivery gets all build event from Unity Cloud.
func (gitLabHook *GitLabHook) receiveGitLabDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the body content.
	requestByteArray, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(Red("[receiveGitLabDelivery Error] ") + err.Error())
		return
	}

	// Get the json content
	gitLabPushHook := gitLabPushHook{}
	err = json.Unmarshal(requestByteArray, &gitLabPushHook)
	if err != nil {
		fmt.Println(Red("[receiveGitLabDelivery Error] ") + err.Error())
		return
	}

	//Broadcast the download info to channel.
	informMessage := ":bookmark: [" + gitLabPushHook.Project.Name + "] " + gitLabHook.MistyRef.Line("gitLabNewRevision", 0) + "\n"
	informMessage += "```Markdown\n"
	for _, commit := range gitLabPushHook.Commits {
		informMessage += "#[" + commit.TimeStamp + "]\n"
		informMessage += "    " + commit.Message
	}
	informMessage += "```"
	gitLabHook.MistyRef.broadcastMessage(informMessage)
}
