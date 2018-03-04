package migrations

import "gogit/database"

// Migrate the database to a new version
func (Migration) MigrateIssueCommentsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.issue_comments (
    id bigserial NOT NULL,
    issue_id bigint NOT NULL,
    reply_to_id bigint NOT NULL,
    comment text NOT NULL,
    attachment_path varchar(255),
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT issue_comments_pkey PRIMARY KEY (id),
    CONSTRAINT issue_comments_issue_id_issues_id FOREIGN KEY (issue_id)
        REFERENCES public.issues (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT issue_comments_reply_to_id_issue_comments_id FOREIGN KEY (reply_to_id)
        REFERENCES public.issue_comments (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT issue_comments_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON DELETE CASCADE
        ON UPDATE CASCADE
    )`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackIssueCommentsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.issue_comments;")

	return
}
