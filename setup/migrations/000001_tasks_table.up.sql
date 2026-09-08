CREATE TYPE task_status AS ENUM ('TODO', 'IN_PROGRESS', 'DONE');
CREATE TYPE task_priority AS ENUM ('LOW', 'MEDIUM', 'HIGH');

CREATE table tasks (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		status task_status NOT NULL,
		priority task_priority NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE UNIQUE INDEX idx_unique_tasks_title_user_id ON tasks(title, user_id);
