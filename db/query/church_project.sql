-- name: CreateChurchProject :exec
INSERT INTO church_projects (church_id, project_name, project_description, target_amount, donated_amount, start_date, end_date)
VALUES ($1, $2, $3, $4, 0.0, $5, $6);

-- name: UpdateChurchProject :exec
UPDATE church_projects
SET
    project_name = COALESCE(sqlc.narg('project_name'),project_name),
    project_description = COALESCE(sqlc.narg('project_description'),project_description),
    target_amount = COALESCE(sqlc.narg('target_amount'),target_amount),
    start_date = COALESCE(sqlc.narg('start_date'),start_date),
    end_date = COALESCE(sqlc.narg('end_date'),end_date),
    visibility = COALESCE(sqlc.narg('visibility'),visibility),
    updated_at = NOW()
WHERE
    id = sqlc.arg('id');

-- name: MarkChurchProjectAsCompleted :exec
UPDATE church_projects
SET completed = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: GetChurchProjects :many
SELECT cp.*, COALESCE(SUM(pd.donation_amount), 0) AS total_donated_amount
FROM church_projects cp
LEFT JOIN project_donations pd ON cp.id = pd.project_id
WHERE cp.church_id = $1
GROUP BY cp.id;

-- name: CreateProjectsDonation :exec
INSERT INTO project_donations (project_id, user_id, donation_amount, donated_at)
VALUES ($1, $2, $3, now());

-- name: GetChurchProjectDetailsOld :many
SELECT cp.*,
       COALESCE(SUM(pd.donation_amount), 0) AS total_donated_amount,
       (
           SELECT JSON_AGG(
                      JSON_BUILD_OBJECT(
                          'user_id', pd.user_id,
                          'total_donation', SUM(pd.donation_amount)
                      )
                  )
           FROM project_donations pd
           WHERE pd.project_id = cp.id
           GROUP BY pd.user_id
       ) AS contributors
FROM church_projects cp
LEFT JOIN project_donations pd ON cp.id = pd.project_id
WHERE cp.id = $1
GROUP BY cp.id;

-- ! Top contributors
-- name: GetChurchProjectTopContributors :many
SELECT user_id, SUM(donation_amount) AS total_donation
FROM project_donations
WHERE project_id = $1
GROUP BY user_id
ORDER BY total_donation DESC
LIMIT 10;

-- ! Donation trends
-- name: GetChurchProjectTrend :many
SELECT DATE_TRUNC('month', donated_at) AS month,
       SUM(donation_amount) AS total_donation
FROM project_donations
WHERE project_id = $1
GROUP BY month
ORDER BY month;

-- name: GetChurchProjectDetailsWithContributors :many
WITH project_details AS (
    SELECT cp.*,
           COALESCE(SUM(pd.donation_amount), 0) AS total_donated_amount
    FROM church_projects cp
    LEFT JOIN project_donations pd ON cp.id = pd.project_id
    WHERE cp.id = $1
    GROUP BY cp.id
), contributors AS (
    SELECT user_id,
           SUM(donation_amount) AS total_donation
    FROM project_donations
    WHERE project_id = $1
    GROUP BY user_id
)
SELECT pd.*,
       (
           SELECT JSON_AGG(
                      JSON_BUILD_OBJECT(
                          'user_id', c.user_id,
                          'total_donation', c.total_donation
                      )
                  )
           FROM contributors c
       ) AS contributors
FROM project_details pd;


-- name: GetChurchProjectDetails :many
WITH project_details AS (
    SELECT cp.*,
           COALESCE(SUM(pd.donation_amount), 0) AS total_donated_amount
    FROM church_projects cp
    LEFT JOIN project_donations pd ON cp.id = pd.project_id
    WHERE cp.id = $1
    GROUP BY cp.id
)
-- , contributors AS (
--     SELECT user_id,
--            SUM(donation_amount) AS total_donation
--     FROM project_donations
--     WHERE project_id = $1
--     GROUP BY user_id
-- )
SELECT pd.*
-- ,
--        (
--            SELECT JSON_AGG(
--                       JSON_BUILD_OBJECT(
--                           'user_id', c.user_id,
--                           'total_donation', c.total_donation
--                       )
--                   )
--            FROM contributors c
--        ) AS contributors
FROM project_details pd;


-- TODO: Paginate
-- name: GetChurchProjectContributors :many
With contributors AS (
    SELECT user_id,
           SUM(donation_amount) AS total_donation
    FROM project_donations
    WHERE project_id = $1
    GROUP BY user_id
)
SELECT 
       (
           SELECT JSON_AGG(
                      JSON_BUILD_OBJECT(
                          'user_id', c.user_id,
                          'total_donation', c.total_donation
                      )
                  )
           FROM contributors c
       ) AS contributors LIMIT $2 OFFSET $3;
-- FROM project_details pd;