-- name: CreatePostImage :exec
INSERT INTO posts_images (post_id, image_url, caption) VALUES ($1, $2, $3);