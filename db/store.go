package db

import (
	"context"
	"database/sql"
	"time"
)

const InsertVideoRequest = "INSERT INTO VideoIndex (link, audio_description, video_description, idx, user_description) VALUES (?, ?, ?, ?, ?)"
const GetMaxIdxRequest = "SELECT max(idx) FROM VideoIndex;"
const GetAllInIndexesRequest = "SELECT DISTINCT link, user_description, idx FROM VideoIndex WHERE idx IN (?)"
const GetAllPagedRequest = "SELECT DISTINCT link, user_description FROM VideoIndex LIMIT ? OFFSET ?"

type Video struct {
	Index            uint64
	Link             string
	AudioDescription string
	VideoDescription string
	UserDescription  string
	CreatedAt        time.Time
}

type Videos []Video

type InsertVideoParameters struct {
	Link             string
	AudioDescription string
	VideoDescription string
	UserDescription  string
}

type VideoStore interface {
	InsertVideo(ctx context.Context, params InsertVideoParameters) (uint64, error)
	GetAllVideoLinksAndUserDescriptionsWithIndexes(ctx context.Context, indexes []int) (Videos, error)
	GetAllVideoLinksAndUserDescriptionsPaged(ctx context.Context, page, size int) (Videos, error)
}

type SQLVideoStore struct {
	db *sql.DB
}

func NewSQLVideoStore(db *sql.DB) SQLVideoStore {
	return SQLVideoStore{
		db: db,
	}
}

func (s *SQLVideoStore) InsertVideo(ctx context.Context, params InsertVideoParameters) (uint64, error) {
	index := s.generateUniqueIndex()
	_, err := s.db.ExecContext(ctx, InsertVideoRequest,
		params.Link, params.AudioDescription, params.VideoDescription, index, params.UserDescription)
	if err != nil {
		return 0, err
	}
	return index, nil
}

func (s *SQLVideoStore) GetAllVideoLinksAndUserDescriptionsWithIndexes(ctx context.Context, indexes []int) (Videos, error) {
	priorities := make(map[int]int)
	for i, v := range indexes {
		priorities[v] = i
	}

	rows, err := s.db.QueryContext(ctx, GetAllInIndexesRequest, indexes)
	if err != nil {
		return nil, err
	}

	entries := make(Videos, len(priorities))
	var entry Video
	for rows.Next() {
		var idx int
		err = rows.Scan(&entry.Link, &entry.UserDescription, &idx)
		if err != nil {
			return nil, err
		}
		entries[priorities[idx]] = entry
	}

	return entries, nil
}

func (s *SQLVideoStore) GetAllVideoLinksAndUserDescriptionsPaged(ctx context.Context, page, size int) (Videos, error) {
	rows, err := s.db.QueryContext(ctx, GetAllPagedRequest, size, (page-1)*size)
	if err != nil {
		return nil, err
	}

	var entries Videos
	var entry Video
	for rows.Next() {
		err = rows.Scan(&entry.Link, &entry.UserDescription)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *SQLVideoStore) generateUniqueIndex() uint64 {
	var maxIdx uint64
	_ = s.db.QueryRow(GetMaxIdxRequest).Scan(&maxIdx) // zero on error

	return maxIdx + 1
}
