SELECT id, code, name, create_time
FROM indexes
ORDER BY id
LIMIT $2 OFFSET $3;
