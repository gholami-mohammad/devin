package models

import "time"

type IssueLabel struct {
	ID         uint64
	Label      string
	Color      string
	ProjectID  uint64
	Project    *Project
	IsBugLable bool `doc:"زمانیکه لیبل های یک پروژه برای بخش مسایل ایجاد میشود اگر این مقدار ۱ باشد یعنی این مسئله مطرح شده یک باگ است و به صورت خودکار به بخش باگ ها نیز ارجاع میشود"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
