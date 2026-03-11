SELECT ticker, short_name, create_time
FROM securities
ORDER BY ticker
    LIMIT $1 OFFSET $2;
