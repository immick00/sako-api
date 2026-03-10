-- name: GetMcdonaldsProduct :one
SELECT * FROM mcdonalds
WHERE id = $1 LIMIT 1;

-- name: GetMcdonaldsProducts :many
SELECT * FROM mcdonalds;

-- name: GetSubwayProducts :many
SELECT * FROM subway;