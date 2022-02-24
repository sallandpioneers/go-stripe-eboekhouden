-- name: CreateCustomer :exec
INSERT INTO customers
(id, stripe_id, boekhouden_id, boekhouden_code)
VALUES
(?,?,?,?);

-- name: GetCustomer :one
SELECT id, stripe_id, boekhouden_id, boekhouden_code
FROM customers
WHERE id = ?;

-- name: GetCustomerBasedOnStripeID :one
SELECT id, stripe_id, boekhouden_id, boekhouden_code
FROM customers
WHERE stripe_id = ?;

-- name: GetCustomerBasedOnBoekhoudenID :one
SELECT id, stripe_id, boekhouden_id, boekhouden_code
FROM customers
WHERE boekhouden_id = ?;

