package migrations

import "devin/database"

// Migrate the database to a new version
func (Migration) MigrateMilestoneCommentsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec(`CREATE TABLE IF NOT EXISTS public.milestone_comments (
    id bigserial NOT NULL,
    milestone_id bigint NOT NULL,
    reply_to_id bigint,
    comment text,
    attachment_path varchar(512),
    created_by_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,

    CONSTRAINT milestone_comments_pkey PRIMARY KEY (id),
    CONSTRAINT milestone_comments_milestone_id_milestones_id FOREIGN KEY (milestone_id)
        REFERENCES public.milestones (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT milestone_comments_reply_to_id_milestones_id FOREIGN KEY (reply_to_id)
        REFERENCES public.milestones (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT milestone_comments_created_by_id_users_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
    );`)

	return
}

// Rollback the database to previous version
func (Migration) RollbackMilestoneCommentsTable() (e error) {
	db := database.NewPGInstance()
	defer db.Close()
	_, e = db.Exec("DROP TABLE IF EXISTS public.milestone_comments;")

	return
}
