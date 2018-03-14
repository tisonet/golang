package main

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"time"
	"log"
)

const (
	MaxIdleConnections = 20
)

type RecommenderProxy struct {
	RecommenderUrl string
	httpClient *http.Client
}

type RecommenderResponse struct {
	Success  bool
	Response []byte
}

func NewRecommenderProxy(recommenderUrl string, recommenderTimeoutMs int) *RecommenderProxy {
	return &RecommenderProxy{recommenderUrl, createHTTPClient(recommenderTimeoutMs)}
}

func (recommender *RecommenderProxy) Recommend(adRequest *AdRequest, response chan RecommenderResponse) {
	adRequest.Proxy = true

	reqBody, err := adRequest.MarshalJSON()
	if err != nil {
		log.Printf("Failed to serialize AdRequest %s\n", err)
		response <- RecommenderResponse{false, nil}
	}

	res, err := recommender.httpClient.Post(recommender.RecommenderUrl, "application/json; charset=utf-8", bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("Failed to post data to Recommender %s\n", err)
		response <- RecommenderResponse{false, nil}
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Printf("Failed to read data from Recommender %s\n", err)
		response <- RecommenderResponse{false, nil}
	} else {
		response <- RecommenderResponse{true, body}
	}
}

func createHTTPClient(recommenderTimeoutMs int) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(recommenderTimeoutMs)* time.Millisecond,
	}

	return client
}
