INSERT INTO portfolio_securities (id, portfolio_id, security_id, amount, create_time)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (portfolio_id, security_id)
DO UPDATE SET amount = EXCLUDED.amount
