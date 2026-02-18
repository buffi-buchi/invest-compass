SELECT id, sec_id, ticker, short_name, type, extra, create_time
FROM securities
WHERE sec_id = $1
