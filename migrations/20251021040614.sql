-- Create "conversations" table
CREATE TABLE "conversations" (
  "id" bigserial NOT NULL,
  "public_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "title" character varying(160) NOT NULL DEFAULT 'Untitled chat',
  "user_id" bigint NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_conversations_deleted_at" to table: "conversations"
CREATE INDEX "idx_conversations_deleted_at" ON "conversations" ("deleted_at");
-- Create index "idx_conversations_public_id" to table: "conversations"
CREATE UNIQUE INDEX "idx_conversations_public_id" ON "conversations" ("public_id");
-- Create index "idx_conversations_user_id" to table: "conversations"
CREATE INDEX "idx_conversations_user_id" ON "conversations" ("user_id");
-- Create "messages" table
CREATE TABLE "messages" (
  "id" bigserial NOT NULL,
  "public_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "conversation_id" bigint NOT NULL,
  "role" character varying(16) NOT NULL,
  "content" text NOT NULL,
  "content_vector" tsvector NULL,
  "ref_id" bigint NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_conversations_messages" FOREIGN KEY ("conversation_id") REFERENCES "conversations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_messages_conversation_id" to table: "messages"
CREATE INDEX "idx_messages_conversation_id" ON "messages" ("conversation_id");
-- Create index "idx_messages_public_id" to table: "messages"
CREATE UNIQUE INDEX "idx_messages_public_id" ON "messages" ("public_id");
-- Create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "public_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "email" character varying(100) NOT NULL,
  "password" character varying(255) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_public_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_public_id" ON "users" ("public_id");
-- Create "user_matrices" table
CREATE TABLE "user_matrices" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "is_create" boolean NULL,
  "is_read" boolean NULL,
  "is_update" boolean NULL,
  "is_delete" boolean NULL,
  "is_upload" boolean NULL,
  "is_download" boolean NULL,
  "is_archive" boolean NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_user_matrices_user_id" UNIQUE ("user_id"),
  CONSTRAINT "fk_users_user_matrix" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
