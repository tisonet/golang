package main

//easyjson:json
type AdRequest struct {
	PageviewId string               `json:"pageview_id"`
	Ibbid      string               `json:"ibbid"`
	Ip         string               `json:"ip"`
	Url        string               `json:"url"`
	UserAgent  string               `json:"user_agent"`
	Source     string               `json:"source"`
	Positions  [] AdRequestPosition `json:"positions"`
	ProxyData  ProxyData            `json:"proxy"`
}

//easyjson:json
type AdRequestPosition struct {
	ImpressionId string `json:"impression_id"`
	ImpId        string `json:"impid"`
	PositionId   string `json:"positionId"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type ProxyData struct {
	UserId string `json:"user_id"`
	UserTargetedStatus UserTargetedStatus `json:"targeted_status"`
}

func (adRequest *AdRequest) isValid() bool {
	return (adRequest.Url != "") && (adRequest.Ibbid != "") && (adRequest.PageviewId != "") && (adRequest.hasValidPositions())
}

func (adRequest *AdRequest) hasValidPositions() bool {
	for _, position := range adRequest.Positions {
		if !position.isValid() {
			return false
		}
	}
	return true
}

func (adRequestPosition *AdRequestPosition) isValid() bool {
	return (adRequestPosition.ImpressionId != "") && (adRequestPosition.ImpId != "") && (adRequestPosition.Width > 0) && (adRequestPosition.Height > 0)
}
