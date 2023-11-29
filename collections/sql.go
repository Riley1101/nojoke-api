package collections

// one to many to products
const CreateProductCollectionTableQuery = `
CREATE TABLE IF NOT EXISTS product_collections (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES admin(id) ON DELETE CASCADE
);
`

const CountProductCollectionsQuery = `
	SELECT COUNT(*) FROM product_collections;
`
const GetProductCollectionsQuery = `
	SELECT id, create_at, products, user_id FROM product_collections;
`
