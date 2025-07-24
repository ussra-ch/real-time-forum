CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        nickname TEXT NOT NULL,
        age INTEGER,
        gender TEXT,
        first_name TEXT,
        last_name TEXT,
        email TEXT UNIQUE,
        password TEXT
    );

CREATE TABLE
    IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        content TEXT,
        title TEXT,
        interest TEXT,
        photo TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    sessions (
        id TEXT PRIMARY KEY,
        user_id INTEGER,
        expires_at DATETIME
    );