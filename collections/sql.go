package collections

const CreateCollectionTableQuery = `
	CREATE TABLE IF NOT EXISTS collection (
		id SERIAL PRIMARY KEY,
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER REFERENCES admin(id) ON DELETE CASCADE
	);
`

const CreateCollectionTypeTableQuery = `
	CREATE TABLE IF NOT EXISTS collection_type (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`

const CountCollectionsQuery = `
	SELECT COUNT(*) FROM collection;
`
const GetCollectionsQuery = `
	SELECT id, create_at, items, user_id FROM collection;
`
