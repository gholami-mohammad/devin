package models

import (
	"time"
)

// Project is the ORM model to interact with projects table in database
type Project struct {
	tableName               struct{}        `sql:"projects"`
	ID                      uint64          `doc:"Auto increment ID"`
	Name                    string          `doc:"Unique name of the project. Unique rule is : a-z, A-Z, 0-9, dash(-), underscore(_) "`
	Title                   string          `doc:"Nullable, Like name without any naming rule."`
	OwnerID                 uint64          `doc:"کد یکتای مالک و سازنده ی پروژه"`
	Owner                   *User           ``
	Description             string          `doc:"Nullable, Full description of project. Possible to link to a notebook."`
	ScheduledStartDate      *time.Time      `doc:"تاریخ پیش بینی شده برای شروع پروژه"`
	StartDate               *time.Time      `doc:"تاریخ شروع پروژه"`
	ScheduledCompletionDate *time.Time      `doc:"تاریخ پیش بینی شده برای کامل و تمام شدن پروژه"`
	CompletionDate          *time.Time      `doc:"تاریخ واقعی اتمام پروژه که توسط مدیر کل پروژه این تاریخ ثبت میشود"`
	Users                   []*ProjectUser  `doc:"List of users who can access this project. This list must be from the company peoples."`
	Tags                    []*TaggedObject `doc:"A HasMany relation, where ModuleID = models.MODULE_PROJECT"`
	DefaultTaskViewID       uint            `doc:"For now, 2 task view is availabel: 1=List view ; 2=Board view"`
	StatusID                uint            `doc:"active, archived, pending, etc"`
	Status                  *ProjectStatus
	OwnerUserID             uint64 `doc:"مالک وسازنده ی این پروژه"`
	OwnerUser               *User
	OwnerCompanyID          uint64 `doc:"این پروژه در کدام سازمان ساخته شده است"`
	OwnerCompany            *User
	ProjectManagerID        uint64 `doc:"مدیر پروژه و مسئول این پروژه کیست"`
	ProjectManager          *User
	GitRepositories         []*Repository
	CreatedByID             uint64
	CreatedBy               *User
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               *time.Time

	EnableWikiModule         bool
	AllowPublicWiki          bool `doc:"If wiki is enable, Is it public?"`
	EnableTasksModule        bool
	EnableMilestonesModule   bool
	EnableFilesModule        bool
	EnableMessagesModule     bool
	EnableTimeLogsModule     bool
	EnableNotebooksModule    bool
	EnableRisksModule        bool
	EnableLinksModule        bool
	EnableBillingModule      bool
	EnableGitModule          bool
	EnableIssueTrackerModule bool
	AllowPublicIssues        bool
	EnableBugTrackerModule   bool
	AllowPublicBugs          bool
	EnableProjectComments    bool
}
