-- name: AddMfaSecret :exec
WITH deactivated_secrets AS (
  UPDATE two_factor_secrets
  SET is_active = FALSE
  WHERE user_id = $1 AND is_active = TRUE
  RETURNING user_id
)
INSERT INTO two_factor_secrets (user_id, secret_key, is_active)
SELECT $1, $2, TRUE;

-- name: GetMfaSecret :one
SELECT secret_key
FROM two_factor_secrets
WHERE user_id = $1 AND is_active = TRUE;

-- name: AddRecoveryCodes :exec
WITH marked_used AS (
  UPDATE two_factor_backup_codes
  SET used = TRUE
  WHERE user_id = $1 AND used = FALSE
  RETURNING user_id
)
INSERT INTO two_factor_backup_codes (user_id, code, used)
SELECT $1, unnest($2::text[]), FALSE;

-- name: GetRecoveryCodes :many
SELECT id, code
FROM two_factor_backup_codes
WHERE user_id = $1 AND used = FALSE;
-- from the two_factor_backup_codes table get all codes with their corresponding id where the user_id = $1 and used = false

-- name: UpdatedRecoveryCodeToUsed :one
UPDATE two_factor_backup_codes
SET used = TRUE
WHERE code = $1 AND user_id = $2 AND used = FALSE
RETURNING user_id;

-- name: UpdatedByIdRecoveryCodeToUsed :one
UPDATE two_factor_backup_codes
SET used = TRUE
WHERE id = $1 AND user_id = $2 AND used = FALSE
RETURNING user_id;

-- name: UpdatedSecretToInActive :exec
UPDATE two_factor_secrets
SET is_active = FALSE
WHERE user_id = $1;