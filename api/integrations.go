package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const GetDescriptionURI = "/api/get_descriptions"
const PostDescriptionURI = "/create_video_index"
const SearchIndexURI = "/search"

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
	httpResp, err := s.sendJSON(req, s.videoProcessingURL+GetDescriptionURI, "POST")
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

type IndexInfo struct {
	VideoURL          string `json:"video_url"`
	VideoDescription  string `json:"video_desc"`
	VideoMovementDesc string `json:"video_movement_desc"`
	SpeechDescription string `json:"speech_desc"`
	Index             int64  `json:"index"`
}

type postIndexRequest struct {
	VideoDescription  string `json:"VideoDescription"`
	VideoMovementDesc string `json:"VideoMovementDesc"`
	SpeechDescription string `json:"SpeechDescription"`
	Index             uint64 `json:"Index"`
}

type postIndexResponse struct {
	Status string `json:"status"`
}

func (s *Server) postIndex(postIndexReq postIndexRequest) (postIndexResponse, error) {
	httpResp, err := s.sendJSON(postIndexReq, s.videoIndexingURL+PostDescriptionURI, "POST")
	if err != nil {
		return postIndexResponse{}, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != 200 {
		return postIndexResponse{}, fmt.Errorf("indexing service return bad status code: %d", httpResp.StatusCode)
	}
	var resp postIndexResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return postIndexResponse{}, err
	}

	return resp, nil
}

type searchIndexResponse struct {
	Indexes []int `json:"ids"`
}

func (s *Server) searchIndex(query string) (searchIndexResponse, error) {
	query = url.QueryEscape(query)
	httpResp, err := s.client.Get(fmt.Sprintf("%s%s?query=%s", s.videoIndexingURL, SearchIndexURI, query))
	if err != nil {
		return searchIndexResponse{}, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != 200 {
		return searchIndexResponse{}, fmt.Errorf("search index service returned bad status code: %d", httpResp.StatusCode)
	}

	var resp searchIndexResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return searchIndexResponse{}, err
	}
	return resp, nil
}

func (s *Server) sendJSON(v any, url, method string) (*http.Response, error) {
	asJson, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(method, url, bytes.NewReader(asJson))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	defer httpReq.Body.Close()

	return s.client.Do(httpReq)
}
