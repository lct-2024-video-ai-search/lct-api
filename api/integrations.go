package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const GetDescriptionURI = "/api/get_descriptions"

type getDescriptionsRequest struct {
	VideoURL         string `json:"video_url"`
	VideoDescription string `json:"video_desc"`
}

type getDescriptionsResponse struct {
	VideoURL          string `json:"video_url"`
	VideoDescription  string `json:"video_desc"`
	VideoMovementDesc string `json:"video_movement_desc"`
	SpeechDescription string `json:"speech_desc"`
}

func (s *Server) getDescriptions(req getDescriptionsRequest) (getDescriptionsResponse, error) {
	asJson, err := json.Marshal(req)
	if err != nil {
		return getDescriptionsResponse{}, err
	}

	httpReq, err := http.NewRequest("POST", GetDescriptionURI, bytes.NewReader(asJson))
	if err != nil {
		return getDescriptionsResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	defer httpReq.Body.Close()

	httpResp, err := s.client.Do(httpReq)
	if err != nil {
		return getDescriptionsResponse{}, err
	}
	defer httpResp.Body.Close()

	var resp getDescriptionsResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return getDescriptionsResponse{}, err
	}

	return resp, nil
}
