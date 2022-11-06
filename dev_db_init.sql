create table IF NOT EXISTS users
(
    id            serial
        constraint users_pk primary key,
    common_id     int     not null,
    email         varchar not null,
    password_hash text,
    roles         text
);

create unique index IF NOT EXISTS users_email_uindex
    on users (email);

create unique index IF NOT EXISTS users_id_uindex
    on users (id);

create unique index IF NOT EXISTS users_common_id_uindex
    on users (common_id);

INSERT INTO users (email, password_hash, common_id, roles)
VALUES ('test@test.ru', MD5('test'), 555, 'test,best')
ON CONFLICT (email) DO NOTHING;

create table IF NOT EXISTS refresh_sessions
(
    id            serial
        constraint refresh_sessions_pk
            primary key,
    subject       int          not null,
    refresh_token text         not null,
    fingerprint   varchar(200) not null,
    expires_in    bigint,
    created_at    timestamp default now()
);

create unique index IF NOT EXISTS refresh_sessions_subject_fingerprint_uindex
    on refresh_sessions (subject, fingerprint);

create unique index IF NOT EXISTS refresh_sessions_refresh_token_uindex
    on refresh_sessions (refresh_token);
