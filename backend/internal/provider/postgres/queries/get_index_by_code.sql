SELECT id, index_code, name, create_time
FROM indexes
WHERE index_code = $1
