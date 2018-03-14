package main

import (
	"log"
	. "github.com/aerospike/aerospike-client-go"
	"time"
)

type UserProfileStorage struct {
	AerospikeClient *Client
	ReadPolicy      *BasePolicy
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

func NewUserProfileStorage(aerospikeHost string) *UserProfileStorage {

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
	}
}

func (ups *UserProfileStorage) GetUserTargeting(userId string, channel chan UserTargetingResult) {
	key, err := NewKey("ups", "targ_users", userId)
	rec, err := ups.AerospikeClient.Get(ups.ReadPolicy, key)

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
			UserTargeting: userTargeting,
		}
	}
}
