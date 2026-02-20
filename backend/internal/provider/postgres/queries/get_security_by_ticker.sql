SELECT ticker, short_name, create_time
FROM securities
WHERE ticker = $1;
