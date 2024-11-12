-- table: users
CREATE TABLE users (
				id TEXT PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				password TEXT NOT NULL,
				email TEXT NOT NULL UNIQUE,
				private_key TEXT NOT NULL,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CHECK (length(username) >= 3)
);

-- table: auth
CREATE TABLE auth (
				id TEXT PRIMARY KEY,
				user_id TEXT NOT NULL,
				token TEXT NOT NULL,
				user_agent TEXT,
				expires_at DATETIME NOT NULL,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
				CHECK (expires_at > created_at)
);

-- Index: auth
CREATE INDEX idx_user_id ON auth (user_id);

-- table: tasks
CREATE TABLE tasks (
				id TEXT PRIMARY KEY,
				title TEXT NOT NULL,
				description TEXT NOT NULL,
				status TEXT NOT NULL CHECK (status IN ('pending', 'in_progress', 'completed')),
			    priority_id TEXT,
				due_date DATETIME,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				user_id TEXT NOT NULL,
				FOREIGN KEY (user_id) REFERENCES users (id),
                FOREIGN KEY (priority_id) REFERENCES task_priorities (id)
);

-- table: task_priorities
CREATE TABLE task_priorities (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL UNIQUE,
                level TEXT NOT NULL CHECK (level IN ('low', 'medium', 'high')),
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- table: categories
CREATE TABLE categories (
				id TEXT PRIMARY KEY,
				name TEXT UNIQUE NOT NULL,
				color TEXT NOT NULL CHECK (color REGEXP '^#[0-9A-F]{6}$')
);

-- table: task_categories
CREATE TABLE task_categories (
				task_id TEXT NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
				category_id TEXT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
				PRIMARY KEY (task_id, category_id)
);
