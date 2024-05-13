CREATE TABLE "posts" (
  "id" uuid PRIMARY KEY NOT NULL,
  "user_id" uuid NOT NULL,
  "content" varchar NOT NULL,
  "total_images" int DEFAULT 0,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz,
  "suspended_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "posts_images" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "post_id" uuid NOT NULL,
  "image_url" varchar,
  "caption" varchar
);

CREATE TABLE "follow" (
  "follower_user_id" uuid not null,
  "following_user_id" uuid not null,
  "created_at" timestamptz DEFAULT (now()),
  PRIMARY KEY ("follower_user_id", "following_user_id")
);


CREATE TABLE "announcements" (
  "id" uuid PRIMARY KEY NOT NULL,
  "user_id" uuid NOT NULL,
  "title" varchar NOT NULL,
  "content" varchar NOT NULL,
  "total_images" int DEFAULT 0,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz,
  "suspended_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE INDEX ON "posts" ("user_id", "id", "created_at");

CREATE INDEX ON "announcements" ("user_id", "id", "created_at");

CREATE INDEX ON "posts_images" ("post_id", "id", "image_url");

CREATE INDEX ON "follow" ("follower_user_id", "following_user_id");

ALTER TABLE "follow" ADD FOREIGN KEY ("follower_user_id") REFERENCES "authentications" ("id");

ALTER TABLE "follow" ADD FOREIGN KEY ("following_user_id") REFERENCES "authentications" ("id");

ALTER TABLE "announcements" ADD FOREIGN KEY ("user_id") REFERENCES "authentications" ("id");

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "authentications" ("id");

ALTER TABLE "posts_images" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");