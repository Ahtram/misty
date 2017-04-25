package bot

import (
	"fmt"
	"log"
	"net/http"

	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

// The hook port we are listening.
var uCloudHookListenPort = ":48769"

// StartUCloudHook start the ucloud router.
func StartUCloudHook() {
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/ucloud/projecta", receiveDelivery)
	log.Fatal(http.ListenAndServe(uCloudHookListenPort, router))
}

// index is a blank landing page.
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "There's no spoon.\n")
}

// receiveDelivery gets all build event from Unity Cloud.
func receiveDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bodyStr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("receiveDelivery read body error! ")
	}

	fmt.Printf("receiveDelivery got: %s", bodyStr)
}
