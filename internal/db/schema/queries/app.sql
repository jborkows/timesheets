-- name: CreateApp :exec
INSERT INTO app (a, b, c) VALUES (:a, :b, :c);

-- name: GetApp :one
SELECT * FROM app WHERE id = :id;

-- name: ListApps :many
SELECT * FROM app;

