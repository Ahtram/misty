package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/julienschmidt/httprouter"
)

// StartGitLabHook start the gitlab router.
func (misty *Misty) StartGitLabHook(endPointID string, port string) {
	router := httprouter.New()
	router.POST("/gitlab/"+endPointID, receiveGitLabDelivery)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// receiveDelivery gets all build event from Unity Cloud.
func receiveGitLabDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[receiveGitLabDelivery got]: " + string(requestDump))
}
