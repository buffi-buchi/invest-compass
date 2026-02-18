SELECT id, ticker, name, create_time
FROM indexes
ORDER BY id
LIMIT $1 OFFSET $2;