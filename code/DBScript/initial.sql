--
-- PostgreSQL database dump
--

-- Dumped from database version 14.2 (Debian 14.2-1.pgdg110+1)
-- Dumped by pg_dump version 14.2 (Debian 14.2-1.pgdg110+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;
SET search_path = public, pg_catalog;

CREATE TABLE IF NOT EXISTS actor (
    nickname VARCHAR(100) primary key,
    fullname VARCHAR(400) NOT NULl,
    about TEXT NOT NULl default '',
    email VARCHAR(150)
);

CREATE UNIQUE INDEX if not exists test_email on actor (lower(email));
CREATE UNIQUE INDEX if not exists test_nickname on actor (lower(nickname));

CREATE TABLE IF NOT EXISTS forum(
    slug VARCHAR(100) primary key,
    title VARCHAR(100) NOT NULL,
    host VARCHAR(100) NOT NULL,
    foreign key (host) references actor (nickname)
        on DELETE CASCADE
);

CREATE UNIQUE INDEX if not exists test_forum on forum (lower(slug));

CREATE SEQUENCE IF NOT EXISTS thread_id_seq;

CREATE TABLE IF NOT EXISTS thread(
    id bigint primary key default nextval('thread_id_seq'),
    title VARCHAR(300) NOT NULL,
    author VARCHAR(100) NOT NULL,
    forum VARCHAR(100),
    message TEXT NOT NULL,
    slug VARCHAR(150),
    created timestamp with time zone DEFAULT now(),
    foreign key (author) references actor (nickname)
        on DELETE CASCADE,
    foreign key (forum) references forum (slug)
        on DELETE CASCADE
);

CREATE UNIQUE INDEX if not exists test_thread on thread (lower(slug));

CREATE SEQUENCE IF NOT EXISTS post_id_seq;

CREATE TABLE IF NOT EXISTS post(
    id bigint primary key default nextval('post_id_seq'),
    parent bigint,
    author VARCHAR(100) references actor (nickname)
    on DELETE CASCADE,
    message TEXT NOT NULL,
    isEdited boolean,
    forum VARCHAR(100) references forum(slug)
    on DELETE CASCADE not null,
    threadid bigint references thread(id)
    on DELETE CASCADE not null,
    created timestamp DEFAULT now(),
    pathtree bigint[]  default array []::bigint[]
);

CREATE SEQUENCE IF NOT EXISTS vote_id_seq;

CREATE TABLE IF NOT EXISTS vote(
    id bigint primary key default nextval('vote_id_seq'),
    threadid bigint references thread (id)
    on DELETE CASCADE not null,
    nickname VARCHAR(100) references actor (nickname)
    on DELETE CASCADE not null,
    voice smallint not null,
    constraint unique_voice unique(threadid,nickname)
);


CREATE OR REPLACE FUNCTION updatePathTree() RETURNS trigger as
$updatePathTree$
Declare
    parent_path         BIGINT[];
begin
    if ( new.parent = 0 ) then
        new.pathtree := array_append(new.pathtree,new.id);
    else
        select pathtree from post where id = new.parent into parent_path;
    new.pathtree := new.pathtree || parent_path || new.id;
    end if;
    Return new;
end;
$updatePathTree$ LANGUAGE plpgsql;

drop trigger if exists updatePathTreeTrigger on post;

CREATE TRIGGER  updatePathTreeTrigger
    BEFORE INSERT
    on post
    for each row
EXECUTE Function updatePathTree();