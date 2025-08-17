BEGIN TRANSACTION;
DROP TABLE IF EXISTS "known_guilds";
CREATE TABLE "known_guilds" (
	"internal_id"	INTEGER NOT NULL,
	"guild_id"	INTEGER NOT NULL,
	"guild_name"	TEXT NOT NULL,
	CONSTRAINT "known_guilds_pk" PRIMARY KEY("internal_id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "known_users";
CREATE TABLE "known_users" (
	"user_id"	INTEGER NOT NULL,
	"last_known_username"	TEXT NOT NULL,
	"internal_id"	INTEGER NOT NULL,
	CONSTRAINT "known_users_pk" PRIMARY KEY("internal_id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "reports";
CREATE TABLE "reports" (
	"internal_id"	INTEGER NOT NULL,
	"message_id"	INTEGER NOT NULL,
	"report_link"	TEXT,
	"user_id"	INTEGER NOT NULL,
	CONSTRAINT "reports_pk" PRIMARY KEY("internal_id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "roles";
CREATE TABLE "roles" (
	"internal_id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"guild_id_fk"	INTEGER,
	"role_id"	INTEGER NOT NULL,
	CONSTRAINT "roles_pk" PRIMARY KEY("internal_id" AUTOINCREMENT),
	CONSTRAINT "guild_id_fk" FOREIGN KEY("guild_id_fk") REFERENCES "known_guilds"
);
DROP TABLE IF EXISTS "version";
CREATE TABLE "version" (
	"major"	INTEGER NOT NULL,
	"minor"	INTEGER NOT NULL,
	PRIMARY KEY("major","minor")
);
INSERT INTO "version" VALUES (0,1);
COMMIT;
