package models

type Permission struct {
	ID uint

	/**
	 * Company level
	 */
	CanViewProject   bool
	CanCreateProject bool
	CanUpdateProject bool
	CanDeleteProject bool
}
