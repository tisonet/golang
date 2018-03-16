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
	adRequestWriter  *AdRequestWriter
	statsDClient     *statsd.Client
}

func NewIbbHandler(ups *UserProfileStorage, recommenderProxy *RecommenderProxy, adrequestWriter *AdRequestWriter, statsDClient *statsd.Client) *IbbHandler {
	return &IbbHandler{
		ups, recommenderProxy, adrequestWriter, statsDClient,
	}
}

func (handler *IbbHandler) handle(ctx *fasthttp.RequestCtx) bool {
	handler.statsDClient.Increment("requests.ibb")
	defer handler.statsDClient.NewTiming().Send("time.requests.ibb")

	adRequest := &AdRequest{}


	if err := adRequest.UnmarshalJSON(ctx.PostBody()); err != nil {
		log.Printf("Failed to parse AdRequest %s\n", err)
		return handler.finalizeEmptyResponse(adRequest, RCMD_FAILED, ctx)
	}

	if !adRequest.isValid() {
		return handler.finalizeEmptyResponse(adRequest, INVALID_AD_REQUEST, ctx)
	}

	userId, isValid := getUserId(adRequest)
	if !isValid {
		return handler.finalizeEmptyResponse(adRequest, INVALID_BBID, ctx)
	}

	userTargetingResponse := handler.callUserProfileStorage(userId)
	if userTargetingResponse.Error != nil {
		log.Printf("Failed to get user targeting %s\n", userTargetingResponse.Error)
		return handler.finalizeEmptyResponse(adRequest, RCMD_FAILED, ctx)
	}

	if !userTargetingResponse.IsUserTargeted {
		return handler.finalizeEmptyResponse(adRequest, NOT_TARGETED_USER, ctx)
	}

	recommenderResponse := handler.callRecommender(userId, adRequest, &userTargetingResponse.UserTargetedStatus)
	if recommenderResponse.Error != nil {
		return handler.finalizeEmptyResponse(adRequest, RCMD_FAILED, ctx)
	}

	ctx.Write(recommenderResponse.Response)
	return true
}


func (handler *IbbHandler) finalizeEmptyResponse(adRequest *AdRequest, status int, ctx *fasthttp.RequestCtx) bool {
	go handler.adRequestWriter.write(adRequest, status)
	ctx.SetStatusCode(fasthttp.StatusNoContent)

	return false
}

func (handler *IbbHandler) callUserProfileStorage(userId string) UserTargetingResponse {
	userTargetingResultChannel := make(chan UserTargetingResponse)
	go handler.ups.GetUserTargeting(userId, userTargetingResultChannel)
	response := <- userTargetingResultChannel
	close(userTargetingResultChannel)
	return response
}

func (handler *IbbHandler) callRecommender(userId string, adRequest *AdRequest, status *UserTargetedStatus) RecommenderResponse {
	recommenderResponseChannel := make(chan RecommenderResponse)
	go handler.recommenderProxy.Recommend(userId, adRequest, status, recommenderResponseChannel)
	response := <-recommenderResponseChannel
	close(recommenderResponseChannel)
	return response
}

func getUserId(adRequest *AdRequest) (string, bool) {
	parts := ibbidToUserIdRegex.FindStringSubmatch(adRequest.Ibbid)
	if len(parts) == 2 {
		return parts[1], true
	}
	return "", false
}