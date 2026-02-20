SELECT ticker, short_name, create_time
FROM indexes
WHERE ticker = $1