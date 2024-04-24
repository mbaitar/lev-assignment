CREATE TABLE IF NOT EXISTS metrics (
   id INTEGER PRIMARY KEY,
   user_id TEXT,
   mrr INTEGER,
   churn INTEGER,
   net_growth INTEGER,
   trading_limit INTEGER,
   last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);