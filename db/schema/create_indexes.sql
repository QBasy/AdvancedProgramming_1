-- create_indexes.sql

-- Create an index on the username column in the users table
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- Create an index on the email column in the users table
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- Create an index on the title column in the photos table
CREATE INDEX IF NOT EXISTS idx_photos_title ON photos (title);