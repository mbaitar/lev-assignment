CREATE TABLE IF NOT EXISTS integration_tokens (
  token_id INTEGER PRIMARY KEY AUTOINCREMENT,
  access_token TEXT NOT NULL,
  refresh_token TEXT NOT NULL,
  connected_account_id TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  account_id TEXT,
  FOREIGN KEY (account_id) REFERENCES users (id)
);