SELECT id, portfolio_id, security_id, amount, price, trade_date, type, note
FROM transactions
WHERE portfolio_id = $1
ORDER BY trade_date DESC
LIMIT $2 OFFSET $3
