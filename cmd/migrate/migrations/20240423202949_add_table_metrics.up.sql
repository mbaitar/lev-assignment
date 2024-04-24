CREATE TABLE IF NOT EXISTS metrics (
   id INTEGER PRIMARY KEY,
   user_id TEXT,
   mrr REAL,
   churn REAL,
   net_growth INTEGER,
   trading_limit REAL,
   last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);