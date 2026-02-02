SELECT ps.id, ps.portfolio_id, ps.security_id, ps.amount, ps.create_time
FROM portfolio_securities ps
WHERE ps.portfolio_id = $1
ORDER BY ps.create_time DESC
