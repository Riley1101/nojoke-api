package product

const CreateProductTableQuery = `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price INT NOT NULL,
		description TEXT NOT NULL,
		discount FLOAT,
		rating FLOAT,
		stock INT NOT NULL,
		brand VARCHAR(255) NOT NULL,
		category_id INT,
		thumbnail VARCHAR(255),
		image VARCHAR(255),
		collection_id INTEGER REFERENCES collections(id) ON DELETE CASCADE
	);
`

const CountProductsQuery = `
	SELECT COUNT(*) FROM products;
`

const GetProductsByCollectionQuery = `
	SELECT p.id,p.name,p.price,p.description,p.discount,
	p.rating,p.stock,p.brand,p.category_id,
	p.thumbnail,p.image,p.collection_id
	FROM products p
	JOIN collections c ON p.collection_id = c.id
	WHERE c.id = $1;
`

const GetProductsWithoutCollectionQuery = `
	SELECT 
	p.id,p.name,p.price,p.description,p.discount,
	p.rating,p.stock,p.brand,p.category_id,
	p.thumbnail,p.image,p.collection_id
	FROM products p
	WHERE p.collection_id IS NULL;
`
