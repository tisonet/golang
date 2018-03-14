package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"github.com/valyala/fasthttp"
	"regexp"
	"github.com/galdor/go-cmdline"
	"os"
	"strconv"
)

var ups *UserProfileStorage
var recommenderProxy *RecommenderProxy
var adrequestWriter *AdRequestWriter
var IbbidToUserIdRegex = regexp.MustCompile(`^BBID-[\d]+-(?P<userId>[\d]+)$`)

const (
	UNKNOWN = 0
	BID = 1
	INVALID_AD_REQUEST = 2
	INVALID_BBID = 3
	NO_ACTIVE_CAMPAIGN = 4
	NOT_TARGETED_USER = 5
	NO_CAMPAIGN_FOR_USER = 6
	WB_LISTS_LIMITED = 7
	FREQ_LIMITED = 8
	BUDGET_LIMITED = 9
	NO_FORMAT_FOR_POSITIONS = 10
	RCMD_FAILED = 11
	RCMD_TIMEOUT = 12

)

func main() {

	cmdline := cmdline.New()

	cmdline.AddOption("u", "recommender-url", "https://recommender-static.gaussalgo.com/","Recommender server address")
	cmdline.AddOption("o","recommender-timeout-ms", "200", "Recommender time for responding in milliseconds")
	cmdline.AddOption("t", "ad-requests-topic", "ad_requests","Kafka topic for storing ad requests")

	cmdline.Parse(os.Args)


	ups = NewUserProfileStorage()
	recommenderTimeoutMs, _ := strconv.Atoi(cmdline.OptionValue("recommender-timeout-ms"))
	recommenderProxy = NewRecommenderProxy(cmdline.OptionValue("recommender-url"), recommenderTimeoutMs)
	adrequestWriter = NewAdRequestWriter("gaussalgo22.colpirio.intra:9092,gaussalgo23.colpirio.intra:9092,gaussalgo44.colpirio.intra:9092", cmdline.OptionValue("ad-requests-topic"))
	log.Fatal(fasthttp.ListenAndServe(":8081", HttpRouter))
}

func HttpRouter(ctx *fasthttp.RequestCtx)  {
	switch string(ctx.Path()) {
	case "/":
		IbbHandler(ctx)
	default:
		ctx.Error("not found", fasthttp.StatusNotFound)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func IbbHandler(ctx *fasthttp.RequestCtx) {
	adRequest:=&AdRequest{}

	ctx.SetContentType( "application/json; charset=UTF-8")


	if err := adRequest.UnmarshalJSON(ctx.PostBody()); err != nil {
		log.Printf("Failed to parse AdRequest %s\n", err)

		go adrequestWriter.write(adRequest, RCMD_FAILED)
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
		return
	}

	if !adRequest.isValid() {
		go adrequestWriter.write(adRequest, INVALID_AD_REQUEST)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}


	var userId, isValid = GetUserId(adRequest)
	if !isValid {
		go adrequestWriter.write(adRequest, INVALID_BBID)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	userTargetingResultChannel := make(chan UserTargetingResult)
	go ups.GetUserTargeting(userId, userTargetingResultChannel)
	userTargetingResult := <-userTargetingResultChannel

	if userTargetingResult.IsUserTargeted {
		recommenderResponseChannel := make(chan RecommenderResponse)
		go recommenderProxy.Recommend(adRequest, recommenderResponseChannel)
		recommenderResponse := <-recommenderResponseChannel

		if recommenderResponse.Success {
			ctx.Write(recommenderResponse.Response)
		}else {
			go adrequestWriter.write(adRequest, RCMD_FAILED)
			ctx.SetStatusCode(fasthttp.StatusNoContent)
		}
	} else {
		go adrequestWriter.write(adRequest, NOT_TARGETED_USER)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}
}

func GetUserId(adRequest *AdRequest) (string, bool) {
	parts := IbbidToUserIdRegex.FindStringSubmatch(adRequest.Ibbid)
	if len(parts) == 2 {
		return parts[1], true
	}
	return "", false
}
