package db

var initStatement = `
CREATE TABLE IF NOT EXISTS retrieved_repositories (hash TEXT UNIQUE, url TEXT);
`
