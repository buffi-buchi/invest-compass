SELECT ticker, short_name, create_time
FROM securities
WHERE ticker = ANY($3)
ORDER BY ticker
    LIMIT $1 OFFSET $2;
