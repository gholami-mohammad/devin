package models

import (
	"time"
)

// Project is the ORM model to interact with projects table in database
type Project struct {
	ID uint64

	// Unique name of the project. Unique rule is : a-z, A-Z, 0-9, dash(-), underscore(_)
	Name string

	// Like name without any naming rule.
	Title *string

	// 1 = private : Only accessabe to owner, members and roots
	// 2 = public in organization : Only accessable to owner, members, root and organization members
	// 3 = public for all : Globaly accessable for all
	VisibilityTypeID uint

	// Nullable, Full description of project. Possible to link to a notebook.
	Description *string

	// Predicted datetime to start the project
	ScheduledStartDate *time.Time

	// Real datetime that the project starts
	StartDate *time.Time

	// Predicted datetime of project completion
	ScheduledCompletionDate *time.Time

	// Real completion datetime
	CompletionDate *time.Time

	// List of users who can access this project.
	// This list must be from the organization peoples.
	Users []ProjectUser

	// A HasMany relation, where ModuleID = models.MODULE_PROJECT
	Tags []TaggedObject

	// For now, 2 task views are availabel: 1=List view ; 2=Board view
	DefaultTaskViewID uint

	// active, archived, pending, etc
	StatusID uint
	Status   *ProjectStatus

	// Who's the owner of the project, It may be not equal to creator of project
	OwnerUserID uint64
	OwnerUser   *User

	// This project created under this organization
	// This field can be NULL
	// If no owner organization selected for the projet, it will be in the 'Personnal projects' group
	OwnerOrganizationID *uint64
	OwnerOrganization   *User

	// Who is the project manager? By default project owner is the project manager.
	ProjectManagerID uint64
	ProjectManager   *User

	// A project can has many git repository to store its source codes.
	GitRepositories []Repository

	// The Creator of this record
	CreatedByID uint64
	CreatedBy   *User

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	EnableWikiModule         bool
	AllowPublicWiki          bool
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
