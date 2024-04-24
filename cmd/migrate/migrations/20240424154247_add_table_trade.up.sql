CREATE TABLE IF NOT EXISTS trades (
    trade_id INTEGER PRIMARY KEY AUTOINCREMENT,
    buyer TEXT,
    arr_traded REAL,
    discount_rate REAL,
    trade_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    roi REAL,
    roi_percentage REAL,
    net_profit REAL
);
