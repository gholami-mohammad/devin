package models

import "time"

type TaskReminderReceiver struct {
	tableName         struct{} `sql:"public.task_reminder_receivers"`
	ID                uint64
	ReminderID        uint64
	Reminder          *TaskReminder
	UserID            uint64
	User              *User
	NotificationTypes string `doc:"A jsonb feild. Allowed types are : sms, email, in-app, telegram"`
	CreatedByID       uint64
	CreatedBy         *User
	CreatedAt         time.Time
}
