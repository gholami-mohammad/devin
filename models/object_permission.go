package models

import "time"

// ObjectPermission :
// تمامی دسترسی های مرتبط با ویرایش، دیدن و حذف یک رکورد برای هر فرد در این جدول مدیریت میشود
// دسترسی بر اساس کد ماژول و کد آیجکت قابل تعیین است
type ObjectPermission struct {
	ID          uint64
	UserID      uint64
	User        *User
	ModuleID    uint `doc:"From this list: models.MODULE_PROJECT, models.MODULE_TASK, models.MODULE_MILESTONE, models.MODULE_REPOSITORY, models.MODULE_SPENT_TIME, models.MODULE_ISSUE_TRACKER, models.MODULE_BUG_TRACKER"`
	ObjectID    uint64
	CanRead     bool `doc:"Permission to view in list of tasks and its details"`
	CanUpdate   bool
	CanDelete   bool
	CreatedByID uint64
	CreatedBy   *User
	CreatedAt   time.Time
}
