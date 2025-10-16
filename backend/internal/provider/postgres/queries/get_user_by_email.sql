SELECT id, email, password, create_time
FROM users
WHERE email = $1
LIMIT 1;
