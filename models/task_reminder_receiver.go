package models

type TaskReminderReceiver struct {
	ID                uint64
	ReminderID        uint64
	Reminder          *TaskReminder
	UserID            uint64
	User              *User
	NotificationTypes string `doc:"A jsonb feild. Allowed types are : sms, email, in-app, telegram"`
}
