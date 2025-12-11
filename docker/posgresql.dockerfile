FROM postgres:15.15-trixie

ENV POSTGRES_DB=myapp \
    POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=postgres

# COPY from the build context root: point to the actual location of init.sql
COPY docker-entrypoint-initdb.d/init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432