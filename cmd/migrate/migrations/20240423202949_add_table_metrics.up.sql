CREATE TABLE IF NOT EXISTS metrics (
   id INTEGER PRIMARY KEY,
   user_id TEXT,
   mrr REAL,
   churned_amount INTEGER,
   churned_mrr REAL,
   churned_percentage REAL,
   net_growth INTEGER,
   trading_limit REAL,
   last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);