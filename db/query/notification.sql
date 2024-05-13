-- name: GetFCMTokenInSession :many
SELECT sessions.fcm_token FROM "sessions"
WHERE "user_id" = ANY($1::uuid[])
AND invalidated_at IS NULL; -- TODO: Change to use the created_at, so as to not skip if the created_at is > 30days ago.
-- AND CURRENT_TIMESTAMP - sessions.created_at >= INTERVAL '30 days';
-- AND  sessions.created_at < NOW() - INTERVAL '30 days';


-- name: CreateNotification :many
INSERT INTO notifications (user_id, type) -- payload,
SELECT user_id, $2
FROM UNNEST($1::uuid[]) AS user_id
RETURNING *;

-- name: CreateNotificationFollow :exec
INSERT INTO follow_notifications (notification_id, following_user_id, author_id)
VALUES ($1, $2, $3);

-- name: CreateNotificationPostComment :exec
INSERT INTO post_comment_notifications (notification_id, post_id, comment_user_id, author_id)
VALUES ($1, $2, $3, $4);

-- name: CreateNotificationPostLike :exec
INSERT INTO post_like_notifications (notification_id, post_id, like_user_id, author_id)
VALUES ($1, $2, $3, $4);

-- name: CreateNotificationAnnouncement :exec
INSERT INTO announcement_notifications (user_id, notification_id, news_id, author_id)
SELECT user_ids.user_id, notification_ids.notification_id, $3, $4
FROM UNNEST($1::uuid[]) WITH ORDINALITY AS user_ids(user_id, ord)
JOIN UNNEST($2::int[]) WITH ORDINALITY AS notification_ids(notification_id, ord)
ON user_ids.ord = notification_ids.ord;


-- name: CreateNotificationAnnouncementOld :exec
INSERT INTO announcement_notifications (user_id, notification_id, news_id, author_id)
SELECT user_id, notification_id, $3, $4
FROM UNNEST($1::uuid[]) AS user_id, UNNEST($2::int[]) AS notification_id;


-- name: CreateNotificationPostMention :exec
INSERT INTO mentions_notifications (mentioning_user_id, notification_id, post_id, author_id)
SELECT user_ids.mentioning_user_id, notification_ids.notification_id, $3, $4
FROM UNNEST($1::uuid[]) WITH ORDINALITY AS user_ids(mentioning_user_id, ord)
JOIN UNNEST($2::int[]) WITH ORDINALITY AS notification_ids(notification_id, ord)
ON user_ids.ord = notification_ids.ord;

-- name: CreateNotificationPrayerInvite :exec
INSERT INTO prayer_invite_notifications (inviting_user_id, notification_id, room_id, author_id)
SELECT user_ids.inviting_user_id, notification_ids.notification_id, $3, $4
FROM UNNEST($1::uuid[]) WITH ORDINALITY AS user_ids(inviting_user_id, ord)
JOIN UNNEST($2::int[]) WITH ORDINALITY AS notification_ids(notification_id, ord)
ON user_ids.ord = notification_ids.ord;


-- name: GetUserNotificationsOld :many
WITH NotificationDetails AS (
    SELECT n.id AS notification_id,
           n.type,
           n.read,
           n.time,
           COALESCE(m.post_id, l.post_id, c.post_id, a.news_id) AS post_announcement_id,
           COALESCE(m.author_id, l.author_id, c.author_id, a.author_id, f.author_id) AS author_id,
           m.mentioning_user_id AS mentioned_user_id,
           l.like_user_id AS liked_user_id,
           c.comment_user_id AS commented_user_id,
           f.following_user_id As following_user_id
    FROM notifications n
    LEFT JOIN mentions_notifications m ON n.id = m.notification_id AND n.type = 'postMention'
    LEFT JOIN post_like_notifications l ON n.id = l.notification_id AND n.type = 'postLike'
    LEFT JOIN post_comment_notifications c ON n.id = c.notification_id AND n.type = 'postComment'
    LEFT JOIN announcement_notifications a ON n.id = a.notification_id AND n.type = 'churchAnnouncement'
    LEFT JOIN follow_notifications f ON n.id = f.notification_id AND n.type = 'follow'
    WHERE n.user_id = $1
      AND (n.type = $2 OR $2 IS NULL OR $2 = '')
)
SELECT nd.notification_id,
       nd.type,
       nd.read,
       nd.time,
       ap.id,
       ap.user_type,
       ap.username AS author_username,
       ep.image_url AS author_profile_image,
       CASE WHEN nd.type = 'postMention' THEN p.content
            WHEN nd.type = 'postLike' THEN p.content
            WHEN nd.type = 'postComment' THEN p.content
            WHEN nd.type = 'churchAnnouncement' THEN an.title
            ELSE NULL
       END AS post_announcement_content,
       nd.post_announcement_id AS post_announcement_id
FROM NotificationDetails nd
LEFT JOIN authentications ap ON nd.author_id = ap.id
LEFT JOIN entity_profiles ep ON ap.id = ep.user_id
LEFT JOIN posts p ON nd.post_announcement_id = p.id AND nd.type IN ('postMention', 'postLike', 'postComment')
LEFT JOIN announcements an ON nd.post_announcement_id = an.id AND nd.type = 'churchAnnouncement';


-- name: GetUserNotifications :many
WITH NotificationDetails AS (
    SELECT n.id AS notification_id,
           n.type,
           n.read,
           n.time,
          COALESCE(m.post_id::VARCHAR, l.post_id::VARCHAR, c.post_id::VARCHAR, a.news_id::VARCHAR, pi.room_id, '') AS post_announcement_id,
           COALESCE(m.author_id, l.author_id, c.author_id, a.author_id, f.author_id, pi.author_id) AS author_id,
           m.mentioning_user_id AS mentioned_user_id,
           l.like_user_id AS liked_user_id,
           c.comment_user_id AS commented_user_id,
           f.following_user_id As following_user_id,
           pi.inviting_user_id AS invited_user_id,
           pr.start_time,
           pr.end_time,
           pr.name AS meeting_name,
           pr.room_id,
           pp.status AS invitation_status
    FROM notifications n
    LEFT JOIN mentions_notifications m ON n.id = m.notification_id AND n.type = 'postMention'
    LEFT JOIN post_like_notifications l ON n.id = l.notification_id AND n.type = 'postLike'
    LEFT JOIN post_comment_notifications c ON n.id = c.notification_id AND n.type = 'postComment'
    LEFT JOIN announcement_notifications a ON n.id = a.notification_id AND n.type = 'churchAnnouncement'
    LEFT JOIN follow_notifications f ON n.id = f.notification_id AND n.type = 'follow'
    LEFT JOIN prayer_invite_notifications pi ON n.id = pi.notification_id AND n.type = 'prayerInvite'
    LEFT JOIN prayer_room pr ON pi.room_id = pr.room_id
    LEFT JOIN prayer_participants pp ON pr.room_id = pp.room_id AND pp.user_id = n.user_id
    WHERE n.user_id = $1
    AND (n.type = $2 OR $2 IS NULL OR $2 = '')
)
SELECT nd.notification_id,
       nd.type,
       nd.read,
       nd.time,
       ap.id,
       ap.user_type,
       ap.username AS author_username,
       ep.image_url AS author_profile_image,
       CASE WHEN nd.type = 'postMention' THEN p.content
            WHEN nd.type = 'postLike' THEN p.content
            WHEN nd.type = 'postComment' THEN p.content
            WHEN nd.type = 'churchAnnouncement' THEN an.title
            WHEN nd.type = 'prayerInvite' THEN nd.meeting_name
            ELSE NULL
       END AS post_announcement_content,
       nd.post_announcement_id AS post_announcement_id,
       nd.start_time,
       nd.end_time,
       nd.room_id,
       nd.invitation_status
FROM NotificationDetails nd
LEFT JOIN authentications ap ON nd.author_id = ap.id
LEFT JOIN entity_profiles ep ON ap.id = ep.user_id
LEFT JOIN posts p ON nd.post_announcement_id = p.id::VARCHAR AND nd.type IN ('postMention', 'postLike', 'postComment')
LEFT JOIN announcements an ON nd.post_announcement_id = an.id::VARCHAR AND nd.type = 'churchAnnouncement';


-- name: MarkNotificationAsRead :exec
UPDATE notifications SET read = true WHERE id = $1;


-- TODO: Add delete for prayerInvite
-- name: DeleteNotifications :exec
WITH ExpandedIds AS (
    SELECT UNNEST($1::int[]) AS notification_id
)
, DeletedNotifications AS (
    DELETE FROM post_comment_notifications 
    WHERE notification_id IN (SELECT id FROM notifications WHERE type = 'postComment' AND id IN (SELECT notification_id FROM ExpandedIds))
    RETURNING notification_id
)
, DeletedFollows AS (
    DELETE FROM follow_notifications 
    WHERE notification_id IN (SELECT id FROM notifications WHERE type = 'follow' AND id IN (SELECT notification_id FROM ExpandedIds))
    RETURNING notification_id
)
, DeletedLikes AS (
    DELETE FROM post_like_notifications 
    WHERE notification_id IN (SELECT id FROM notifications WHERE type = 'postLike' AND id IN (SELECT notification_id FROM ExpandedIds))
    RETURNING notification_id
)
, DeletedMentions AS (
    DELETE FROM mentions_notifications 
    WHERE notification_id IN (SELECT id FROM notifications WHERE type = 'postMention' AND id IN (SELECT notification_id FROM ExpandedIds))
    RETURNING notification_id
)
, DeletedAnnouncements AS (
    DELETE FROM announcement_notifications 
    WHERE notification_id IN (SELECT id FROM notifications WHERE type = 'churchAnnouncement' AND id IN (SELECT notification_id FROM ExpandedIds))
    RETURNING notification_id
)
DELETE FROM notifications
WHERE id IN (
    SELECT notification_id FROM DeletedNotifications
    UNION ALL
    SELECT notification_id FROM DeletedFollows
    UNION ALL
    SELECT notification_id FROM DeletedLikes
    UNION ALL
    SELECT notification_id FROM DeletedMentions
    UNION ALL
    SELECT notification_id FROM DeletedAnnouncements
);