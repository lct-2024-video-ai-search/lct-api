package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const DefaultPageSize = 4
const DefaultPageNumber = 1

const InsertVideoRequest = "INSERT INTO VideoIndex (link, audio_description, video_description, idx, user_description) VALUES (?, ?, ?, ?, ?)"
const GetMaxIdxRequest = "SELECT max(idx) FROM VideoIndex;"
const GetAllInIndexesRequest = "SELECT DISTINCT link, user_description, idx FROM VideoIndex WHERE idx IN (?)"
const GetAllPagedRequest = "SELECT DISTINCT link, user_description FROM VideoIndex LIMIT ? OFFSET ?"

type indexVideoRequest struct {
	Link        string `json:"link" binding:"required"`
	Description string `json:"description"`
}

type indexVideoResponse struct {
	indexVideoRequest
}

func (s *Server) indexVideo(ctx *gin.Context) {
	var req indexVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	desc, err := s.getDescriptions(getDescriptionsRequest{
		VideoURL:         req.Link,
		VideoDescription: req.Description,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var maxIdx int64
	_ = s.db.QueryRow(GetMaxIdxRequest).Scan(&maxIdx) // zero on error

	newIdx := maxIdx + 1
	_, err = s.db.ExecContext(context.Background(), InsertVideoRequest,
		desc.VideoURL, desc.SpeechDescription, desc.VideoMovementDesc, newIdx, desc.VideoDescription)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// TODO: send to video processing

	_, err = s.postIndex(postIndexRequest{
		VideoDescription:  desc.VideoDescription,
		VideoMovementDesc: desc.VideoMovementDesc,
		SpeechDescription: desc.SpeechDescription,
		Index:             maxIdx + 1,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, indexVideoResponse{req})
}

type videoEntry struct {
	Link        string `json:"link"`
	Description string `json:"description"`
}

type videoEntryResponse []videoEntry

func (s *Server) videosPaged(ctx *gin.Context) {
	page, size := DefaultPageNumber, DefaultPageSize
	pageQuery, sizeQuery := ctx.Query("page"), ctx.Query("size")
	if maybePage, err := strconv.Atoi(pageQuery); err == nil && pageQuery != "" && maybePage > 0 {
		page = maybePage
	}
	if maybeSize, err := strconv.Atoi(sizeQuery); err == nil && sizeQuery != "" && maybeSize > 0 {
		size = maybeSize
	}

	rows, err := s.db.Query(GetAllPagedRequest, size, (page-1)*size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var entries videoEntryResponse
	var entry videoEntry
	for rows.Next() {
		err = rows.Scan(&entry.Link, &entry.Description)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		entries = append(entries, entry)
	}

	ctx.JSON(http.StatusOK, entries)
}

func (s *Server) searchVideo(ctx *gin.Context) {
	searchQuery := ctx.Query("text")
	if searchQuery == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("empty search query")))
		return
	}
	resp, err := s.searchIndex(searchQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println(resp.Indexes[:20])
	videos, err := s.fetchVideosSortedByIndex(resp.Indexes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, videos)
}

func (s *Server) fetchVideosSortedByIndex(indexes []int) ([]videoEntry, error) {
	priorities := make(map[int]int)
	for i, v := range indexes {
		priorities[v] = i
	}

	rows, err := s.db.Query(GetAllInIndexesRequest, indexes)
	if err != nil {
		return nil, err
	}

	entries := make(videoEntryResponse, len(priorities))
	var entry videoEntry
	for rows.Next() {
		var idx int
		err = rows.Scan(&entry.Link, &entry.Description, &idx)
		if err != nil {
			return nil, err
		}
		entries[priorities[idx]] = entry
	}

	return entries, nil
}
