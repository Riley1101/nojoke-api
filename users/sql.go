package user

const CreateUserTableQuery = `	
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		phone VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		age INTEGER,
		image VARCHAR(255),
		password VARCHAR(255) NOT NULL
	);
`
