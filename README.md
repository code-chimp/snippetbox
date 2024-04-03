# Snippetbox

This is a simple web application written in Go that allows users to create, view,
edit, and delete snippets of text. It is based on the book ["Let's Go"][letsgo]
by [Alex Edwards](https://www.alexedwards.net/).

## Prerequisites

- [Go 1.22](https://go.dev/) or later
- [Air](https://github.com/cosmtrek/air) for live reloading:
  ```shell
  go install github.com/cosmtrek/air
  ```
- [Docker](https://www.docker.com/) for running MariaDB in a container

## Prep the Database

To create a Docker instance of MariaDB for development purposes, run the appropriate script from the `./data-scripts` directory:

```shell
./data-scripts/setup.sh    
```
or:

```powershell
.\data-scripts\setup.ps1
```

Create table and seed some data - login as `snippets-admin` and run:

Contents available in `./data-scripts/seed_schema.sql`:
```mariadb
CREATE TABLE snippets (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(100) NOT NULL,
  content TEXT,
  created DATETIME NOT NULL,
  expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets (created);

INSERT INTO snippets
  (title, content, created, expires)
VALUES
  ('An old silent pond', 'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n- Matsuo Bashō', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)),
  ('Over the wintry forest', 'Over the\nwintry forest, winds howl in rage\nwith no leaves to blow.\n\n- Natsume Soseki', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)),
  ('First autumn morning', 'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n- Murakami Kijo', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY));

-- For the latter chapter on session management
CREATE TABLE sessions (
  token  CHAR(43)     PRIMARY KEY,
  data   BLOB         NOT NULL,
  expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry ON sessions (expiry);
```

## Generate a Self-Signed Certificate

For development purposes, generate a self-signed certificate:

```shell
cd ./tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```
On Windows if you have the default Go installation path:

```powershell
cd .\tls
go run "C:\Program Files\Go\src\crypto\tls\generate_cert.go" --rsa-bits=2048 --host=localhost
```


[letsgo]: https://lets-go.alexedwards.net/ "Let's Go"
