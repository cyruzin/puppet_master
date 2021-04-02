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
