INSERT INTO users (id, email, password, create_time)
VALUES ($1, $2, $3, $4)
RETURNING id, create_time;
