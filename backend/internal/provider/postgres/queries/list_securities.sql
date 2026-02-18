SELECT id, sec_id, ticker, short_name, type, extra, create_time
FROM securities
WHERE ($3::uuid[] IS NULL OR id = ANY($3))
ORDER BY create_time DESC
LIMIT $1 OFFSET $2
