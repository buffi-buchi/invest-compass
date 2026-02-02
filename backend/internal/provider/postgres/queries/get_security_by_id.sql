SELECT id, sec_id, ticker, short_name, type, extra, create_time
FROM securities
WHERE id = $1
