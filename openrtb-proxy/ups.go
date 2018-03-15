package main

import (
	"log"
	"time"
	. "github.com/aerospike/aerospike-client-go"
	"gopkg.in/alexcesaro/statsd.v2"
)

type UserProfileStorage struct {
	AerospikeClient *Client
	ReadPolicy      *BasePolicy
	statsDClient    *statsd.Client
}

type UserTargetingResult struct {
	IsUserTargeted bool
	UserTargeting  UserTargeting
}

type UserTargeting struct {
	TargetedStaticCampaignsCount     int
	TargetedStaticCampaignsTimestamp int
	TargetedShopsCount               map[interface{}]interface{}
}

func NewUserProfileStorage(aerospikeHost string, statsDClient *statsd.Client) *UserProfileStorage {

	client, err := NewClient(aerospikeHost, 3000)

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

func (ups *UserProfileStorage) GetUserTargeting(userId string, channel chan UserTargetingResult) {
	ups.statsDClient.Increment("requests.ups")
	key, err := NewKey("ups", "targ_users", userId)

	timing := ups.statsDClient.NewTiming()
	rec, err := ups.AerospikeClient.Get(ups.ReadPolicy, key)
	timing.Send("time.requests.ups")

	if err != nil {
		log.Printf("Failed to get ups profile %s\n", err)

		channel <- UserTargetingResult{
			IsUserTargeted: false,
		}
		return
	}

	if rec == nil {
		channel <- UserTargetingResult{
			IsUserTargeted: false,
		}
	} else {
		userTargeting := UserTargeting{0, 0, nil}
		targetedStaticCampaignsCount, ok := rec.Bins["trg_sc"]
		if ok {
			userTargeting.TargetedStaticCampaignsCount = targetedStaticCampaignsCount.(int)
		}

		targetedStaticCampaignsTimestamps, ok := rec.Bins["trg_sc_ts"]
		if ok {
			userTargeting.TargetedStaticCampaignsTimestamp = targetedStaticCampaignsTimestamps.(int)
		}

		targetedShopsCount, ok := rec.Bins["trg_shops"]
		if ok {
			userTargeting.TargetedShopsCount = targetedShopsCount.(map[interface{}]interface{})
		}

		channel <- UserTargetingResult{
			IsUserTargeted: true,
			UserTargeting:  userTargeting,
		}
	}
}
