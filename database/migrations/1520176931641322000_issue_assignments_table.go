package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateIssueAssignmentsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.issue_assignments (
    id bigserial NOT NULL,
    issue_id bigint NOT NULL,
    user_id bigint NOT NULL,
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT issue_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT issue_assignments_issue_id_issues_id FOREIGN KEY (issue_id)
        REFERENCES public.issues (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT issue_assignments_user_id_users_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT issue_assignments_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackIssueAssignmentsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.issue_assignments;")

	return
}
