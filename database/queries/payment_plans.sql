-- name: CreatePaymentPlan :one
INSERT INTO payment_plans (id, user_id, currency, amount, status) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, user_id, currency, amount, status, created_at, updated_at;

-- name: ListPaymentPlansByUserID :many
SELECT id, user_id, currency, amount, status, created_at, updated_at FROM payment_plans
WHERE user_id = $1
ORDER BY created_at DESC;
