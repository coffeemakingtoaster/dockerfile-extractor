package db

var initStatement = `
CREATE TABLE IF NOT EXISTS retrieved_repositories (hash TEXT UNIQUE, repo TEXT, scraped INTEGER DEFAULT 0);
`
