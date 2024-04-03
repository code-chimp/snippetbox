#!/usr/bin/env pwsh

docker run --name snippetbox_dev `
    -p 3306:3306 `
    -e 'MARIADB_ROOT_PASSWORD=P@ssw0rd' `
    -e 'MARIADB_USER=snippets-admin' `
    -e 'MARIADB_PASSWORD=S00p3r*S3kr1t' `
    -e 'MARIADB_DATABASE=snippetbox' `
    -d mariadb:11
