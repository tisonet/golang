package main

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"github.com/valyala/fasthttp"
	"log"
	"regexp"
)

var ibbidToUserIdRegex = regexp.MustCompile(`^BBID-[\d]+-(?P<userId>[\d]+)$`)

type IbbHandler struct {
	ups              *UserProfileStorage
	recommenderProxy *RecommenderProxy
	adrequestWriter  *AdRequestWriter
	statsDClient     *statsd.Client
}

func NewIbbHandler(ups *UserProfileStorage, recommenderProxy *RecommenderProxy, adrequestWriter *AdRequestWriter, statsDClient *statsd.Client) *IbbHandler {
	return &IbbHandler{
		ups, recommenderProxy, adrequestWriter, statsDClient,
	}
}

func (handler *IbbHandler) handle(ctx *fasthttp.RequestCtx) {
	handler.statsDClient.Increment("requests.ibb")
	defer handler.statsDClient.NewTiming().Send("time.requests.ibb")

	adRequest := &AdRequest{}

	ctx.SetContentType("application/json; charset=UTF-8")

	if err := adRequest.UnmarshalJSON(ctx.PostBody()); err != nil {
		log.Printf("Failed to parse AdRequest %s\n", err)

		go handler.adrequestWriter.write(adRequest, RCMD_FAILED)
		ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
		return
	}

	if !adRequest.isValid() {
		go handler.adrequestWriter.write(adRequest, INVALID_AD_REQUEST)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	userId, isValid := getUserId(adRequest)
	if !isValid {
		go handler.adrequestWriter.write(adRequest, INVALID_BBID)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	userTargetingResultChannel := make(chan UserTargetingResult)
	go handler.ups.GetUserTargeting(userId, userTargetingResultChannel)
	userTargetingResult := <-userTargetingResultChannel

	if userTargetingResult.IsUserTargeted {
		recommenderResponseChannel := make(chan RecommenderResponse)
		go handler.recommenderProxy.Recommend(adRequest, recommenderResponseChannel)
		recommenderResponse := <-recommenderResponseChannel

		if recommenderResponse.Success {
			ctx.Write(recommenderResponse.Response)
		} else {
			go handler.adrequestWriter.write(adRequest, RCMD_FAILED)
			ctx.SetStatusCode(fasthttp.StatusNoContent)
		}
	} else {
		go handler.adrequestWriter.write(adRequest, NOT_TARGETED_USER)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}
}


func getUserId(adRequest *AdRequest) (string, bool) {
	parts := ibbidToUserIdRegex.FindStringSubmatch(adRequest.Ibbid)
	if len(parts) == 2 {
		return parts[1], true
	}
	return "", false
}