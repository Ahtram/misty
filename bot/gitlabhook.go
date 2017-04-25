package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/julienschmidt/httprouter"
)

// The hook port we are listening.
var gitLabHookListenPort = ":48770"

// StartGitLabHook start the gitlab router.
func StartGitLabHook() {
	router := httprouter.New()
	router.POST("/gitlab/projecta", receiveGitLabDelivery)
	log.Fatal(http.ListenAndServe(gitLabHookListenPort, router))
}

// receiveDelivery gets all build event from Unity Cloud.
func receiveGitLabDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[receiveGitLabDelivery got]: " + string(requestDump))
}
