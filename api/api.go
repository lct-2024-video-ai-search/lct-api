package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"lct-backend/db"
	"lct-backend/transform"
	"net/http"
	"strconv"
)

const DefaultSearchQuery = "fire"
const DefaultPageSize = 4
const DefaultPageNumber = 1

type indexVideoRequest struct {
	Link        string `json:"link" binding:"required"`
	Description string `json:"description"`
}

type indexVideoResponse struct {
	indexVideoRequest
}

// indexVideo godoc
// @Summary      Индексировать видео
// @Description  Индексировать видео в базе сервиса
// @Accept       json
// @Produce      json
// @Param        video body api.indexVideoRequest true "ссылка и описание видео"
// @Success      200  {object}  api.indexVideoResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /index [post]
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

	if errors.Is(err, ErrBadRequest) {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	idx, err := s.store.InsertVideo(ctx, db.InsertVideoParameters{
		Link:             desc.VideoURL,
		AudioDescription: desc.SpeechDescription,
		VideoDescription: desc.VideoMovementDesc,
		UserDescription:  desc.VideoDescription,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = s.postIndex(postIndexRequest{
		VideoDescription:  desc.VideoDescription,
		VideoMovementDesc: desc.VideoMovementDesc,
		SpeechDescription: desc.SpeechDescription,
		Index:             idx,
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

	videos, err := s.store.GetAllVideoLinksAndUserDescriptionsPaged(ctx, page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := transform.Map(videos, func(videos db.Video) videoEntry {
		return videoEntry{
			Link:        videos.Link,
			Description: videos.UserDescription,
		}
	})

	ctx.JSON(http.StatusOK, response)
}

// searchVideo godoc
// @Summary      Поиск видео
// @Description  Поиск видео по заданному тексту
// @Produce      json
// @Param        text query string true "поисковый запрос"
// @Success      200  {object}  api.videoEntryResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /search [get]
func (s *Server) searchVideo(ctx *gin.Context) {
	searchQuery := ctx.Query("text")
	if searchQuery == "" {
		searchQuery = DefaultSearchQuery // stub to filter inappropriate mature content
	}
	resp, err := s.searchIndex(searchQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	videos, err := s.store.GetAllVideoLinksAndUserDescriptionsWithIndexes(ctx, resp.Indexes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := transform.Map(videos, func(videos db.Video) videoEntry {
		return videoEntry{
			Link:        videos.Link,
			Description: videos.UserDescription,
		}
	})

	ctx.JSON(http.StatusOK, response)
}

func (s *Server) mainPage(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"Документация": "https://api-zvezdolet.ddns.net/docs",
		"Индексация":   "https://api-zvezdolet.ddns.net/index",
		"Поиск":        "https://api-zvezdolet.ddns.net/search",
	})
}
