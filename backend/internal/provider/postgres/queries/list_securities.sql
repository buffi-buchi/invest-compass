SELECT id, sec_id, ticker, short_name, type, extra, create_time
FROM securities
ORDER BY create_time DESC
LIMIT $1 OFFSET $2
