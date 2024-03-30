# Snippetbox

This is a simple web application written in Go that allows users to create, view,
edit, and delete snippets of text. It is based on the book ["Let's Go"][letsgo]
by [Alex Edwards](https://www.alexedwards.net/).

## Prep the Database

Create a Docker instance of MariaDB for development purposes:

```shell
docker run --name snippetbox_dev \
    -p 3306:3306 \
    -e 'MARIADB_ROOT_PASSWORD=P@ssw0rd' \
    -e 'MARIADB_USER=snippets-admin' \
    -e 'MARIADB_PASSWORD=S00p3r*S3kr1t' \
    -e 'MARIADB_DATABASE=snippetbox' \
    -d mariadb:11
```

Create table and seed some data - login as `snippets-admin` and run:

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
  ('An old silent pond', 'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n- Matsuo Basho', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)),
  ('Over the wintry forest', 'Over the\nwintry forest, winds howl in rage\nwith no leaves to blow.\n\n- Natsume Soseki', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)),
  ('First autumn morning', 'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n- Murakami Kijo', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY));

```

[letsgo]: https://lets-go.alexedwards.net/ "Let's Go"
