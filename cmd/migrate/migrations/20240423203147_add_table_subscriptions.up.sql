CREATE TABLE IF NOT EXISTS subscriptions (
     id TEXT PRIMARY KEY,
     customer TEXT,
     status TEXT,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     amount INTEGER,
     currency TEXT
);