package main

import (
	"log"
	"time"
	. "github.com/aerospike/aerospike-client-go"
	"gopkg.in/alexcesaro/statsd.v2"
	"strings"
)

type UserProfileStorage struct {
	AerospikeClient *Client
	ReadPolicy      *BasePolicy
	statsDClient    *statsd.Client
}

type UserTargetingResponse struct {
	Error              error
	IsUserTargeted     bool
	UserTargetedStatus UserTargetedStatus
}

type UserTargetedStatus struct {
	TargetedStaticCampaignsCount     int            `json:"static_campaigns"`
	TargetedStaticCampaignsTimestamp int            `json:"static_campaigns_ts"`
	TargetedShops                    map[string]int `json:"visited_shops"`
	ViewedProducts                   int            `json:"viewed_products"`
}

func NewUserProfileStorage(aerospikeHosts string, statsDClient *statsd.Client) *UserProfileStorage {
	hostNames := strings.Split(aerospikeHosts, ",")

	var hosts []*Host
	for _, hostname := range hostNames {
		hosts = append(hosts,  NewHost(hostname, 3000))
	}

	client, err := NewClientWithPolicyAndHost(NewClientPolicy(), hosts...)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	policy := NewPolicy()
	policy.MaxRetries = 0
	policy.Timeout = time.Duration(time.Millisecond * 2000)

	return &UserProfileStorage{
		AerospikeClient: client,
		ReadPolicy:      policy,
		statsDClient:    statsDClient,
	}
}

func (ups *UserProfileStorage) GetUserTargeting(userId string, channel chan UserTargetingResponse) {
	ups.statsDClient.Increment("requests.ups")
	key, err := NewKey("ups", "targ_users", userId)

	timing := ups.statsDClient.NewTiming()
	rec, err := ups.AerospikeClient.Get(ups.ReadPolicy, key)
	timing.Send("time.requests.ups")

	if err != nil {
		log.Printf("Failed to get ups profile %s\n", err)

		channel <- UserTargetingResponse{
			Error: err,
		}
		return
	}

	if rec == nil {
		channel <- UserTargetingResponse{
			IsUserTargeted: false,
		}
	} else {
		status := UserTargetedStatus{0, 0, nil, 0}
		targetedStaticCampaignsCount, ok := rec.Bins["trg_sc"]
		if ok {
			status.TargetedStaticCampaignsCount = targetedStaticCampaignsCount.(int)
		}

		targetedStaticCampaignsTimestamps, ok := rec.Bins["trg_sc_ts"]
		if ok {
			status.TargetedStaticCampaignsTimestamp = targetedStaticCampaignsTimestamps.(int)
		}

		targetedShopsCount, ok := rec.Bins["trg_shops"]
		if ok {
			status.TargetedShops = make(map[string]int)

			for shopId, productsCount := range targetedShopsCount.(map[interface{}]interface{}) {
				viewedProductsCount := productsCount.(int)
				status.TargetedShops[shopId.(string)] = viewedProductsCount
				status.ViewedProducts += viewedProductsCount
			}
		}

		channel <- UserTargetingResponse{
			IsUserTargeted:     true,
			UserTargetedStatus: status,
		}
	}
}

func (ups *UserProfileStorage) Close() {
	ups.AerospikeClient.Close()
}