SELECT id, ticker, name, create_time
FROM indexes
ORDER BY ticker
LIMIT $1 OFFSET $2;