package models

import "time"

type WikiPage struct {
	tableName   struct{} `sql:"public.wiki_pages"`
	ID          uint64
	WikiID      uint64
	Wiki        *Wiki
	Title       string `doc:"A unique title in Wiki"`
	Content     string
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
