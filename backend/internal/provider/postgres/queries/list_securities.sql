SELECT ticker, short_name, create_time
FROM securities
WHERE ($3::text[] IS NULL OR ticker = ANY($3))
ORDER BY ticker
    LIMIT $1 OFFSET $2;
