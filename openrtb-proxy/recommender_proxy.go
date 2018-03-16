package main

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"time"
	"log"
	"gopkg.in/alexcesaro/statsd.v2"
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
}

func NewRecommenderProxy(recommenderUrl string, recommenderTimeoutMs int, statsDClient *statsd.Client) *RecommenderProxy {
	return &RecommenderProxy{recommenderUrl, createHTTPClient(recommenderTimeoutMs), statsDClient}
}

func (recommender *RecommenderProxy) Recommend(userId string, adRequest *AdRequest, userTargetedStatus *UserTargetedStatus, response chan RecommenderResponse) {
	recommender.statsDClient.Increment("requests.recommender")
	adRequest.ProxyData = ProxyData{
		UserId:userId,
		UserTargetedStatus: *userTargetedStatus,
	}

	reqBody, err := adRequest.MarshalJSON()
	if err != nil {
		log.Printf("Failed to serialize AdRequest %s\n", err)
		response <- RecommenderResponse{err, nil}
	}

	timing := recommender.statsDClient.NewTiming()
	res, err := recommender.httpClient.Post(recommender.RecommenderUrl, "application/json; charset=utf-8", bytes.NewReader(reqBody))

	if err != nil {
		recommender.statsDClient.Increment("error.recommender")

		log.Printf("Failed to post data to Recommender %s\n", err)
		response <- RecommenderResponse{err, nil}
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	timing.Send("time.requests.recommender")

	if err != nil {
		log.Printf("Failed to read data from Recommender %s\n", err)
		response <- RecommenderResponse{err, nil}
	} else {
		response <- RecommenderResponse{nil, body}
	}
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
