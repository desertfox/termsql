
FROM mysql:latest

# Copy certificates from the first stage
COPY ./docker/certs /certs

RUN chown mysql: /certs/*

ENV MYSQL_ROOT_PASSWORD=root

COPY ./docker/my.cnf /etc/mysql/my.cnf

RUN chown :mysql /etc/mysql/my.cnf

# When container starts, this script will be executed.
# It creates a new user with the specified username and password.
COPY ./docker/scripts/ /docker-entrypoint-initdb.d/