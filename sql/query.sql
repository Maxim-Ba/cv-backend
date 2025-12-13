-- name: GetTechnology :one
SELECT * FROM technology
WHERE id = $1;

-- name: ListTechnologies :many
SELECT * FROM technology
ORDER BY title;

-- name: CreateTechnology :one
INSERT INTO technology (id, title, description, logo_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTechnology :one
UPDATE technology
SET title = $2, description = $3, logo_url = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTechnology :exec
DELETE FROM technology
WHERE id = $1;

-- name: GetTag :one
SELECT * FROM tag
WHERE id = $1;

-- name: ListTags :many
SELECT * FROM tag
ORDER BY name;

-- name: CreateTag :one
INSERT INTO tag (id, name, hex_color)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEducation :one
SELECT * FROM education
WHERE id = $1;

-- name: ListEducations :many
SELECT * FROM education
ORDER BY year DESC;

-- name: CreateEducation :one
INSERT INTO education (id, name, year, course, organization)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWorkHistory :one
SELECT * FROM work_history
WHERE id = $1;

-- name: ListWorkHistories :many
SELECT * FROM work_history
ORDER BY period_start DESC;

-- name: CreateWorkHistory :one
INSERT INTO work_history (id, name, about, logo_url, period_start, period_end, what_i_did, projects)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetTechnologiesByTag :many
SELECT t.* FROM technology t
JOIN technologies_tag tt ON t.id = tt.technology_id
WHERE tt.tag_id = $1
ORDER BY t.title;

-- name: GetTechnologiesByWorkHistory :many
SELECT t.* FROM technology t
JOIN work_history_technology wht ON t.id = wht.technology_id
WHERE wht.work_history_id = $1
ORDER BY t.title;

-- name: AddTechnologyToTag :exec
INSERT INTO technologies_tag (tag_id, technology_id)
VALUES ($1, $2);

-- name: RemoveTechnologyFromTag :exec
DELETE FROM technologies_tag
WHERE tag_id = $1 AND technology_id = $2;

-- name: AddTechnologyToWorkHistory :exec
INSERT INTO work_history_technology (work_history_id, technology_id)
VALUES ($1, $2);

-- name: RemoveTechnologyFromWorkHistory :exec
DELETE FROM work_history_technology
WHERE work_history_id = $1 AND technology_id = $2;

-- name: SearchTechnologies :many
SELECT * FROM technology
WHERE title ILIKE '%' || $1 || '%'
   OR description ILIKE '%' || $1 || '%'
ORDER BY title;
