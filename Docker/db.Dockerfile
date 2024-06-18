# syntax=docker/dockerfile:1
FROM postgres:latest
COPY init.sql /docker-entrypoint-initdb.d/
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB library