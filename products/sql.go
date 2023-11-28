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
		image VARCHAR(255)
	);
`
