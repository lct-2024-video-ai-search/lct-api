// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import ()

type Video struct {
	ID          int32  `json:"id"`
	Link        string `json:"link"`
	Description string `json:"description"`
}
