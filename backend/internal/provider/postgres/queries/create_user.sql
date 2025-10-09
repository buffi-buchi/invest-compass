INSERT INTO users (email, password, create_time)
VALUES ($1, $2, $3)
RETURNING id, create_time;
