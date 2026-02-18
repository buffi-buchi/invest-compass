SELECT id, ticker, name, create_time
FROM indexes
WHERE ticker = $1
