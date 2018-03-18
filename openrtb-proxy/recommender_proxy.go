package main

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"time"
	"log"
	"gopkg.in/alexcesaro/statsd.v2"
	"github.com/valyala/fasthttp"
)

const (
	MaxIdleConnections = 20
)

type RecommenderProxy struct {
	RecommenderUrl string
	httpClient     *http.Client
	statsDClient   *statsd.Client
}

type RecommenderResponse struct {
	Error    error
	Response []byte
	StatusCode int
}

func NewRecommenderProxy(recommenderUrl string, recommenderTimeoutMs int, statsDClient *statsd.Client) *RecommenderProxy {
	return &RecommenderProxy{recommenderUrl, createHTTPClient(recommenderTimeoutMs), statsDClient}
}

func (recommender *RecommenderProxy) Recommend(userId string, adRequest *AdRequest, userTargetedStatus *UserTargetedStatus, responseChannel chan RecommenderResponse) {
	recommender.statsDClient.Increment("requests.recommender")
	adRequest.ProxyData = ProxyData{
		UserId:userId,
		UserTargetedStatus: *userTargetedStatus,
	}

	reqBody, err := adRequest.MarshalJSON()
	if err != nil {
		log.Printf("Failed to serialize AdRequest %s\n", err)
		responseChannel <- createErrorResponse(err)
	}

	defer recommender.statsDClient.NewTiming().Send("time.requests.recommender")

	res, err := recommender.httpClient.Post(recommender.RecommenderUrl, "application/json; charset=utf-8", bytes.NewReader(reqBody))
	if err != nil {
		recommender.statsDClient.Increment("error.recommender")

		log.Printf("Failed to post data to Recommender %s\n", err)
		responseChannel <- createErrorResponse(err)
		return
	}

	if res.StatusCode == fasthttp.StatusNoContent {
		responseChannel <- RecommenderResponse{nil, nil, fasthttp.StatusNoContent}
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read data from Recommender %s\n", err)
		responseChannel <- createErrorResponse(err)
	} else {
		responseChannel <- RecommenderResponse{nil, body, fasthttp.StatusOK}
	}

	defer res.Body.Close()
}

func createHTTPClient(recommenderTimeoutMs int) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(recommenderTimeoutMs) * time.Millisecond,
	}

	return client
}

func createErrorResponse(err error) RecommenderResponse {
	return  RecommenderResponse{err, nil, fasthttp.StatusNoContent}
}
