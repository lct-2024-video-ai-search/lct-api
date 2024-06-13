package api

import (
	"context"
	"github.com/gin-gonic/gin"
	db "lct-backend/db/sqlc"
	"net/http"
)

const DefaultPageSize = 20
const DefaultPageNumber = 1

const InsertVideoRequest = "INSERT INTO VideoIndex (link, audio_description, video_description, idx, user_description) VALUES (?, ?, ?, ?) RETURNING *;"

type indexVideoRequest struct {
	Link        string `json:"link" binding:"required"`
	Description string `json:"description"`
}

type searchVideoRequest struct {
	Query    string `json:"query" binding:"required"`
	Page     int32  `json:"page" binding:"min=1"`
	PageSize int32  `json:"pageSize" binding:"min=5,max=100"`
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

	s.db.ExecContext(context.Background(), InsertVideoRequest)
	ctx.Status(http.StatusCreated)
}

func (s *Server) searchVideo(ctx *gin.Context) {
	var req searchVideoRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Page == 0 {
		req.Page = DefaultPageNumber
	}

	if req.PageSize == 0 {
		req.PageSize = DefaultPageSize
	}

	arg := db.ListVideosParams{
		Offset: req.PageSize * (req.Page - 1),
		Limit:  req.PageSize,
	}

	videos, err := s.store.ListVideos(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, videos)
}
