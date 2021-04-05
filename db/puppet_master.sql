CREATE TABLE IF NOT EXISTS permissions (
  id SERIAL NOT NULL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS roles (
  id SERIAL NOT NULL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL NOT NULL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  email VARCHAR(80) NOT NULL UNIQUE,
  password VARCHAR(80) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permission_role (
  permission_id SMALLINT NOT NULL REFERENCES permissions (id) ON UPDATE CASCADE ON DELETE CASCADE,
  role_id SMALLINT NOT NULL REFERENCES roles (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS role_user (
  role_id SMALLINT NOT NULL REFERENCES roles (id) ON UPDATE CASCADE ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);


INSERT INTO permissions ("id", "name", "description", "created_at", "updated_at") VALUES
(1,	'view user',	'Can view a user',	'2021-04-05 13:32:49.483076+00',	'2021-04-05 13:32:49.483076+00'),
(2,	'create user',	'Can create user',	'2021-04-05 13:32:57.862105+00',	'2021-04-05 13:32:57.862105+00'),
(3,	'edit user',	'Can edit user',	'2021-04-05 13:33:07.451287+00',	'2021-04-05 13:33:07.451287+00'),
(4,	'delete user',	'Can delete user',	'2021-04-05 13:33:15.617075+00',	'2021-04-05 13:33:15.617075+00'),
(5,	'view role',	'Can view role',	'2021-04-05 13:33:25.620844+00',	'2021-04-05 13:33:25.620844+00'),
(6,	'create role',	'Can create role',	'2021-04-05 13:33:35.849704+00',	'2021-04-05 13:33:35.849704+00'),
(7,	'edit role',	'Can edit role',	'2021-04-05 13:33:43.189558+00',	'2021-04-05 13:33:43.189558+00'),
(8,	'delete role',	'Can delete role',	'2021-04-05 13:33:54.616185+00',	'2021-04-05 13:33:54.616185+00'),
(9,	'view permission',	'Can view permission',	'2021-04-05 13:34:33.37066+00',	'2021-04-05 13:34:33.37066+00'),
(10,	'create permission',	'Can create permission',	'2021-04-05 13:35:00.82438+00',	'2021-04-05 13:35:00.82438+00'),
(11,	'edit permission',	'Can edit permission',	'2021-04-05 13:35:31.411928+00',	'2021-04-05 13:35:31.411928+00'),
(12,	'delete permission',	'Can delete permission',	'2021-04-05 13:35:41.087124+00',	'2021-04-05 13:35:41.087124+00'),
(13,	'give permission to role',	'Can give permission to role',	'2021-04-05 13:36:14.140329+00',	'2021-04-05 13:36:14.140329+00'),
(14,	'remove permission to role',	'Can remove permission to role',	'2021-04-05 13:36:28.691047+00',	'2021-04-05 13:36:28.691047+00'),
(15,	'sync permission to role',	'Can sync permission to role',	'2021-04-05 13:36:29.224283+00',	'2021-04-05 13:36:29.224283+00');

INSERT INTO roles ("id", "name", "description", "created_at", "updated_at") VALUES
(1,	'Admin',	'Admin of the system',	'2021-04-05 13:37:48.531415+00',	'2021-04-05 13:37:48.531415+00');

INSERT INTO users ("id", "name", "email", "password", "created_at", "updated_at") VALUES
(1,	'The Admin',	'admin@admin.com',	'$2a$06$DDKssq9NZAFGSGaLx8mjB.6Cl0NdnkQNSla49s8I6u1g8g7nmNK42',	'2021-04-05 13:38:57.594285+00',	'2021-04-05 13:38:57.594285+00');

INSERT INTO "role_user" ("role_id", "user_id") VALUES (1,	1);

-- Workaround to fix primary key out of sync
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users)+1);
SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles)+1);
SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions)+1);