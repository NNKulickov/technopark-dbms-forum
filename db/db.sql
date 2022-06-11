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

CREATE UNLOGGED TABLE IF NOT EXISTS actor (
    nickname VARCHAR(100) primary key,
    fullname VARCHAR(400) NOT NULl,
    about TEXT NOT NULl default '',
    email VARCHAR(150)
);

CREATE UNIQUE INDEX if not exists test_email on actor (lower(email));
CREATE UNIQUE INDEX if not exists test_nickname on actor (lower(nickname));

CREATE UNLOGGED TABLE IF NOT EXISTS forum(
    slug VARCHAR(100) primary key,
    title VARCHAR(100) NOT NULL,
    host VARCHAR(100) NOT NULL,
    foreign key (host) references actor (nickname)
        on DELETE CASCADE
);

CREATE UNIQUE INDEX if not exists test_forum on forum (lower(slug));

CREATE SEQUENCE IF NOT EXISTS thread_id_seq;

CREATE UNLOGGED TABLE IF NOT EXISTS thread(
    id bigint primary key default nextval('thread_id_seq'),
    title VARCHAR(300) NOT NULL,
    author VARCHAR(100) NOT NULL,
    forum VARCHAR(100),
    message TEXT NOT NULL,
    slug VARCHAR(150),
    created timestamp with time zone DEFAULT now(),
    votes int default 0,
    foreign key (author) references actor (nickname)
        on DELETE CASCADE,
    foreign key (forum) references forum (slug)
        on DELETE CASCADE
);

CREATE UNIQUE INDEX if not exists test_thread on thread (lower(slug));

CREATE SEQUENCE IF NOT EXISTS post_id_seq;

CREATE UNLOGGED TABLE IF NOT EXISTS post(
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

-- CREATE UNIQUE INDEX if not exists upost_parent_author on post (lower(author),message,parent);

CREATE SEQUENCE IF NOT EXISTS vote_id_seq;

CREATE UNLOGGED TABLE IF NOT EXISTS vote(
    id bigint primary key default nextval('vote_id_seq'),
    threadid bigint references thread (id)
    on DELETE CASCADE not null,
    nickname VARCHAR(100) references actor (nickname)
    on DELETE CASCADE not null,
    voice smallint not null,
    constraint unique_voice unique(threadid,nickname)
);

CREATE OR REPLACE FUNCTION insertPathTree() RETURNS trigger as
$insertPathTree$
Declare
    parent_path         BIGINT[];
begin
    if ( new.parent = 0 ) then
        new.pathtree := array_append(new.pathtree,new.id);
    else
        select pathtree from post where id = new.parent into parent_path;
    new.pathtree := new.pathtree || parent_path || new.id;
    end if;
    return new;
end;
$insertPathTree$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insertThreadsVotes() RETURNS trigger as
$insertThreadsVotes$
begin
    update thread set votes = votes + new.voice where id = new.threadid;
    return new;
end;
$insertThreadsVotes$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION updateThreadsVotes() RETURNS trigger as
$updateThreadsVotes$
begin
    update thread set votes = votes + new.voice - old.voice where id = new.threadid;
    return new;
end;
$updateThreadsVotes$ LANGUAGE plpgsql;

drop trigger if exists insertPathTreeTrigger on post;
drop trigger if exists insertThreadsVotesTrigger on vote;
drop trigger if exists updateThreadsVotesTrigger on vote;

CREATE TRIGGER  insertPathTreeTrigger
    BEFORE INSERT
    on post
    for each row
EXECUTE Function insertPathTree();

CREATE TRIGGER  insertThreadsVotesTrigger
    AFTER INSERT
    on vote
    for each row
EXECUTE Function insertThreadsVotes();

CREATE TRIGGER  updateThreadsVotesTrigger
    AFTER UPDATE
    on vote
    for each row
EXECUTE Function updateThreadsVotes();

CREATE INDEX IF NOT EXISTS forum_slug_hash ON forum using hash (slug);
CREATE INDEX IF NOT EXISTS thread_parenttree_post on post (threadid,pathtree);
CREATE INDEX IF NOT EXISTS first_parent_post on post ((pathtree[1]),pathtree);
