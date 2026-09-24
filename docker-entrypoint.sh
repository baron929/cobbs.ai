#!/bin/bash
set -euo pipefail

: "${MYSQL_ROOT_PASSWORD:?MYSQL_ROOT_PASSWORD is required}"

service mariadb start

mysqladmin -u root password ${MYSQL_ROOT_PASSWORD}

export dataSourceName="root:${MYSQL_ROOT_PASSWORD}@tcp(127.0.0.1:3306)/"

exec /server --createDatabase=true