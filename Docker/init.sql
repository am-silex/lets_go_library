create table public.books
(
    id       serial,
    authorid integer,
    year     integer,
    isbn     varchar
);

alter table public.books
    owner to postgres;

create table public.authors
(
    id            serial,
    first_name    varchar,
    last_name     varchar,
    bio           varchar,
    date_of_birth integer
);

alter table public.authors
    owner to postgres;

