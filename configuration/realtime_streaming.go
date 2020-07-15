package configuration

import (
	"encoding/json"
	"net/http"
)

const (
	SlotStatePending = "Pending"
	SlotStateReady   = "Ready"
	SlotStateFailed  = "Failed"
)

type RealtimeStreamingSlot struct {
	Id                  string                     `json:"id,omitempty"`
	State               string                     `json:"state,omitempty"`
	Name                string                     `json:"name,omitempty"`
	Region              string                     `json:"region,omitempty"`
	Profiles            []RealtimeStreamingProfile `json:"profiles,omitempty"`
	Password            string                     `json:"password,omitempty"`
	IPGeoMatch          string                     `json:"ipGeoMatch,omitempty"`
	MediaVaultEnabled   bool                       `json:"mediaVaultEnabled,omitempty"`
	MediaVaultSecretKey string                     `json:"mediaVaultSecretKey,omitempty"`
}

type RealtimeStreamingProfile struct {
	VideoBitrate int `json:"videoBitrate"`
	AudioBitrate int `json:"audioBitrate"`
}

func (c *ConfigurationClient) GetRealtimeStreamingSlot(slotId string, shortname string) (*RealtimeStreamingSlot, *http.Response, error) {
	<-c.rateLimiter
	realtimeStreamingSlot := &RealtimeStreamingSlot{}

	body, response, err := c.Auth.HTTPGet(c.BaseUrl + "/webrtc/shortname/" + shortname + "/slots/" + slotId)

	if err != nil {
		return nil, response, err
	}

	json.Unmarshal(body, realtimeStreamingSlot)

	return realtimeStreamingSlot, response, nil
}

func (c *ConfigurationClient) CreateRealtimeStreamingSlot(shortname string, slot *RealtimeStreamingSlot) (*RealtimeStreamingSlot, *http.Response, error) {
	<-c.rateLimiter

	jsonRequest, _ := json.Marshal(slot)

	body, response, err := c.Auth.HTTPPost(c.BaseUrl+"/webrtc/shortname/"+shortname+"/slots", string(jsonRequest))

	if err != nil {
		return nil, response, err
	}

	responseSlot := &RealtimeStreamingSlot{}
	json.Unmarshal(body, responseSlot)

	return responseSlot, response, nil
}

func (c *ConfigurationClient) DeleteRealtimeStreamingSlot(slotId string, shortname string) (*http.Response, error) {
	<-c.rateLimiter
	_, response, err := c.Auth.HTTPDelete(c.BaseUrl + "/webrtc/shortname/" + shortname + "/slots/" + slotId)

	if err != nil {
		return response, err
	}

	return response, nil
}
