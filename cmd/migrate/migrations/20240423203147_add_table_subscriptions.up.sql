CREATE TABLE IF NOT EXISTS subscriptions (
     id TEXT PRIMARY KEY,
     customer TEXT,
    user_id TEXT,
     status TEXT,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     amount REAL,
     currency TEXT,
    end_date TIMESTAMP,
     cancel_at_period_end INTEGER
);