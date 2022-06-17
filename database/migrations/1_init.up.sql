CREATE TABLE IF NOT EXISTS tests (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS commits (
    id INTEGER PRIMARY KEY,
    commit_sha STRING NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS uploads (
    id INTEGER PRIMARY KEY,
    commit_id INTEGER NOT NULL REFERENCES commits (id),
    time INTEGER NOT NULL,
    url TEXT
);

CREATE TABLE IF NOT EXISTS results (
    test_id INTEGER NOT NULL REFERENCES tests(id),
    upload_id INTEGER NOT NULL REFERENCES uploads(id),
    commit_id INTEGER NOT NULL REFERENCES commits(id),
    success INTEGER NOT NULL,
    output TEXT
);
