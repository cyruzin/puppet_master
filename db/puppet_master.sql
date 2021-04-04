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
(5,	'view user',	'Can view a user',	'2021-04-02 06:04:49.709798+00',	'2021-04-02 06:04:49.709798+00'), 
(6,	'create user',	'Can create a user',	'2021-04-02 06:04:58.674248+00',	'2021-04-02 06:04:58.674248+00'), 
(7,	'edit user',	'Can edit a user',	'2021-04-02 06:05:05.695544+00',	'2021-04-02 06:05:05.695544+00'), 
(8,	'delete user',	'Can delete a user',	'2021-04-02 06:05:19.941473+00',	'2021-04-02 06:05:19.941473+00'), 
(1,	'create role',	'Can create a role',	'2021-04-02 05:17:53.643469+00',	'2021-04-02 05:17:53.643469+00'), 
(2,	'edit role',	'Can edit a role',	'2021-04-02 05:18:01.63039+00',	'2021-04-02 05:18:01.63039+00'), 
(3,	'view role',	'Can view a role',	'2021-04-02 05:18:10.89572+00',	'2021-04-02 05:18:10.89572+00'), 
(4,	'delete role',	'Can delete a role',	'2021-04-02 05:18:20.596286+00',	'2021-04-02 05:18:20.596286+00'), 
(9,	'view permission',	'Can view permissions',	'2021-04-03 17:16:21.882724+00',	'2021-04-03 17:16:21.882724+00'), 
(10,	'create permission',	'Can create permission',	'2021-04-03 17:16:38.100071+00',	'2021-04-03 17:16:38.100071+00'), 
(11,	'edit permission',	'Can edit permission',	'2021-04-03 17:16:45.132572+00',	'2021-04-03 17:16:45.132572+00'), 
(12,	'delete permission',	'Can delete permission',	'2021-04-03 17:16:53.457964+00',	'2021-04-03 17:16:53.457964+00');

INSERT INTO roles ("id", "name", "description", "created_at", "updated_at") VALUES 
(1,	'Admin',	'Admin of the system',	'2021-04-02 05:19:12.73807+00',	'2021-04-02 05:19:12.73807+00');

INSERT INTO users ("name", "email", "password", "created_at", "updated_at")  VALUES 
('The Admin',	'admin@admin.com',	'$2a$06$DDKssq9NZAFGSGaLx8mjB.6Cl0NdnkQNSla49s8I6u1g8g7nmNK42',	'2021-04-02 05:25:28.999125+00',	'2021-04-02 05:25:28.999125+00');

INSERT INTO role_user ("role_id", "user_id") VALUES (1,	1);