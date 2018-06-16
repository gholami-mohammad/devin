package models

// ProjectSearch is model to performe search on Project model
type ProjectSearch struct {
	ID          *uint64
	Name        *string
	Title       *string
	Description *string
	StatusID    *uint

	// Just search for project which this user is its owner
	OwnerUserID *uint64

	// If UserID != nil, then search for projects with OwnerUserID and all project members
	// to find projects which this user is its owner or its member
	UserID *uint64
}
