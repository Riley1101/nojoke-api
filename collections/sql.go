package collections

const CreateCollectionTableQuery = `
	CREATE TABLE IF NOT EXISTS collection (
		id SERIAL PRIMARY KEY,
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		items INTEGER REFERENCES item(id) ON DELETE CASCADE,
		user_id INTEGER REFERENCES admin(id) ON DELETE CASCADE
	);
`
