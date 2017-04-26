package bot

import (
	"log"
	"net/http"

	"fmt"
	"strings"

	"encoding/json"
	"io/ioutil"

	"strconv"

	"github.com/julienschmidt/httprouter"
)

// The hook port we are listening. (Moved to config)
// var uCloudHookListenPort = ":48769"
var uCloudAccessToken = "f82a9507bb519996aadbe2caebe4cc51"

var uCloudAPIURL = "https://build-api.cloud.unity3d.com"

var uCloudProjectBuildSuccessEvent = "ProjectBuildSuccess"

var uCloudShareLinkURL = "https://developer.cloud.unity3d.com/share/"

// UCloudHook represents an Unity cloud hook we are listening.
type UCloudHook struct {
	UCloudHookEndPoint string
	UCloudHookPort     string
	UCloudAccessToken  string
	MistyRef           *Misty
}

type uCloudProjectBuildSuccess struct {
	ProjectName       string                         `json:"projectName"`
	BuildTargetName   string                         `json:"buildTargetName"`
	ProjectGUID       string                         `json:"projectGuid"`
	OrgForeignKey     string                         `json:"orgForeignKey"`
	BuildNumber       int                            `json:"buildNumber"`
	BuildStatus       string                         `json:"buildStatus"`
	LastBuiltRevision string                         `json:"lastBuiltRevision"`
	StartedBy         string                         `json:"startedBy"`
	Platform          string                         `json:"platform"`
	ScmType           string                         `json:"scmType"`
	Links             uCloudProjectBuildSuccessLinks `json:"links"`
}

type uCloudProjectBuildSuccessLinks struct {
	APISelf          uCloudLink `json:"api_self"`
	DashboardURL     uCloudLink `json:"dashboard_url"`
	DashboardProject uCloudLink `json:"dashboard_project"`
	DashboardSummary uCloudLink `json:"dashboard_summary"`
	DashboardLog     uCloudLink `json:"dashboard_log"`
}

type uCloudLink struct {
	Method string `json:"method"`
	Href   string `json:"href"`
}

type uCloudLinkToShare struct {
	ShareID string `json:"shareid"`
}

// StartUCloudHook start the ucloud router.
func (uCloudHook *UCloudHook) StartUCloudHook() {
	router := httprouter.New()
	router.POST(uCloudHook.UCloudHookEndPoint, uCloudHook.receiveUCloudDelivery)
	log.Fatal(http.ListenAndServe(":"+uCloudHook.UCloudHookPort, router))
}

// receiveDelivery gets all build event from Unity Cloud.
func (uCloudHook *UCloudHook) receiveUCloudDelivery(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check if this is a build success notification. (Header)
	event := r.Header.Get("X-Unitycloudbuild-Event")
	if strings.Compare(event, uCloudProjectBuildSuccessEvent) == 0 {
		// This is a success event!
		// Get the body content.
		requestByteArray, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(Red("[receiveUCloudDelivery Error] ") + err.Error())
			return
		}

		// Get the json content.
		uCloudProjectBuildSuccess := uCloudProjectBuildSuccess{}
		err = json.Unmarshal(requestByteArray, &uCloudProjectBuildSuccess)
		if err != nil {
			fmt.Println(Red("[receiveUCloudDelivery Error] ") + err.Error())
			return
		}

		fmt.Println(Green("[Project Build Success!] [" + uCloudProjectBuildSuccess.ProjectName + "] [" + uCloudProjectBuildSuccess.BuildTargetName + "] [" + strconv.Itoa(uCloudProjectBuildSuccess.BuildNumber) + "]"))

		requestShareLinkURL := uCloudAPIURL + uCloudProjectBuildSuccess.Links.APISelf.Href + "/share"

		// Request a share link.
		// fmt.Println("Requesting a share link: " + requestShareLinkURL)

		shareLinkRequest, err := http.NewRequest("POST", requestShareLinkURL, nil)
		if err != nil {
			fmt.Println(Red("[POST Request Error] ") + err.Error())
			return
		}
		shareLinkRequest.Header.Set("Authorization", "Basic "+uCloudHook.UCloudAccessToken)
		shareLinkRequest.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		shareLinkResp, err := client.Do(shareLinkRequest)
		if err != nil {
			fmt.Println(Red("[POST Request Error] ") + err.Error())
			return
		}

		// //Just a log
		// requestDump, err := httputil.DumpResponse(shareLinkResp, true)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(Green("[Share Link Response]") + string(requestDump))

		respShareLinkByteArray, err := ioutil.ReadAll(shareLinkResp.Body)
		if err != nil {
			fmt.Println(Red("[respShareLinkByteArray Error] ") + err.Error())
			return
		}

		uCloudLinkToShare := uCloudLinkToShare{}
		err = json.Unmarshal(respShareLinkByteArray, &uCloudLinkToShare)
		if err != nil {
			fmt.Println(Red("[Unmarshal respShareLinkByteArray Error] ") + err.Error())
			return
		}
		fmt.Println(Green("[Got Share Link]") + uCloudShareLinkURL + uCloudLinkToShare.ShareID)

		//Broadcast the download info to channel.
		informMessage := ":inbox_tray: " + uCloudHook.MistyRef.Line("uCloudNewBuild", 0)
		informMessage += " [" + uCloudProjectBuildSuccess.ProjectName + "] [" + uCloudProjectBuildSuccess.BuildTargetName + "] [" + strconv.Itoa(uCloudProjectBuildSuccess.BuildNumber) + "]\n"
		informMessage += uCloudShareLinkURL + uCloudLinkToShare.ShareID
		uCloudHook.MistyRef.broadcastMessage(informMessage)

		return
	}
	fmt.Println("[Ignoring uCloud event]: " + event)
}
