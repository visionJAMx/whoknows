CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pages (
    title TEXT PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    language TEXT NOT NULL DEFAULT 'en'
        CHECK (language IN ('en', 'da')),
    last_updated TIMESTAMP,
    content TEXT NOT NULL
);