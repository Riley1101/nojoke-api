package collections

// one to many to products
const CreateCollectionTableQuery = `
CREATE TABLE IF NOT EXISTS collections(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES admin(id) ON DELETE CASCADE
);
`

const CountCollectionsQuery = `
	SELECT COUNT(*) FROM collections;
`
const GetCollectionsQuery = `
	SELECT id, create_at,  user_id FROM collections LIMIT $1 OFFSET $2;
`
