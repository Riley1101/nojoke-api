package auth

const CreateAdminTableQuery = `
	CREATE TABLE IF NOT EXISTS admin (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`
