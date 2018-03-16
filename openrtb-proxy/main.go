package main

import (
	"log"
	"github.com/valyala/fasthttp"
	"github.com/galdor/go-cmdline"
	"os"
	"strconv"
	"gopkg.in/alexcesaro/statsd.v2"
)

var ibbHandler *IbbHandler

const (
	UNKNOWN                 = 0
	BID                     = 1
	INVALID_AD_REQUEST      = 2
	INVALID_BBID            = 3
	NO_ACTIVE_CAMPAIGN      = 4
	NOT_TARGETED_USER       = 5
	NO_CAMPAIGN_FOR_USER    = 6
	WB_LISTS_LIMITED        = 7
	FREQ_LIMITED            = 8
	BUDGET_LIMITED          = 9
	NO_FORMAT_FOR_POSITIONS = 10
	RCMD_FAILED             = 11
	RCMD_TIMEOUT            = 12
)

func main() {

	cmdline := cmdline.New()

	cmdline.AddOption("u", "recommender-url", "https://recommender-static.gaussalgo.com/", "Recommender server address")
	cmdline.AddOption("o", "recommender-timeout-ms", "200", "Recommender time for responding in milliseconds")
	cmdline.AddOption("t", "ad-requests-topic", "ad_requests", "Kafka topic for storing ad requests")
	cmdline.AddOption("a", "aerospike-host", "gaussalgo9.colpirio.intra", "Aerospike hostname")
	cmdline.AddOption("k", "kafka-hosts", "gaussalgo22.colpirio.intra:9092,gaussalgo23.colpirio.intra:9092,gaussalgo44.colpirio.intra:9092", "Kafka brokers hostnames")
	cmdline.AddOption("s", "statsd-server", "statsd.marathon.mesos:8125", "Statsd server address")
	cmdline.Parse(os.Args)

	statsDClient, err := statsd.New(statsd.Address(cmdline.OptionValue("statsd-server")), statsd.Prefix("openrtb_proxy"))
	if err != nil {
		log.Fatalf("Cant connect to statsD server %s\n", cmdline.OptionValue("statsd-server"))
		return
	}


	ups := NewUserProfileStorage(cmdline.OptionValue("aerospike-host"), statsDClient.Clone(statsd.SampleRate(0.2)))
	recommenderTimeoutMs, _ := strconv.Atoi(cmdline.OptionValue("recommender-timeout-ms"))
	recommenderProxy := NewRecommenderProxy(cmdline.OptionValue("recommender-url"), recommenderTimeoutMs,  statsDClient.Clone(statsd.SampleRate(0.2)))
	adRequestWriter := NewAdRequestWriter(cmdline.OptionValue("kafka-hosts"), cmdline.OptionValue("ad-requests-topic"))

	ibbHandler = NewIbbHandler(ups, recommenderProxy, adRequestWriter, statsDClient.Clone(statsd.SampleRate(0.2)))

	log.Fatal(fasthttp.ListenAndServe(":8080", HttpRouter))

	statsDClient.Close()
}

func HttpRouter(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		ctx.SetContentType("application/json; charset=UTF-8")

		ibbHandler.handle(ctx)
	default:
		ctx.Error("not found", fasthttp.StatusNotFound)
	}
}


