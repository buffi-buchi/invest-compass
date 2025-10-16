SELECT id, email, password, create_time
FROM users
WHERE id = $1
LIMIT 1;
