-- name: GetMcdonaldsProduct :one
SELECT * FROM mcdonalds
WHERE id = $1 LIMIT 1;

-- name: GetMcdonaldsProducts :many
SELECT * FROM mcdonalds;

-- name: GetSubwayProducts :many
SELECT * FROM subway;

-- name: GetChickFilAProducts :many
SELECT * FROM chickfila;

-- name: GetBurgerKingProducts :many
SELECT * FROM burgerking;

-- name: GetTacoBellProducts :many
SELECT * FROM tacobell;

-- name: GetPopeyesProducts :many
SELECT * FROM popeyes;

-- name: CreateOnboardingResponse :one
INSERT INTO onboarding_responses (
    user_id, goal, weight, height_feet, height_inches, age_range, days_per_week, activity_level, cravings, dislikes 
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;
