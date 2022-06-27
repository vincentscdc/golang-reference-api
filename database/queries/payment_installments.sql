-- name: CreatePaymentInstallments :one
INSERT INTO payment_installments (id, payment_plan_id, currency, amount, due_at, status) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, payment_plan_id, currency, amount, due_at, status, created_at, updated_at;

-- name: ListPaymentInstallmentsByPlanID :many
SELECT id, payment_plan_id, currency, amount, due_at, status, created_at, updated_at FROM payment_installments
WHERE payment_plan_id = $1
ORDER BY due_at;
