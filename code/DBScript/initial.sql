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
    fullname VARCHAR(400) NOT NULl DEFAULT '',
    about TEXT NOT NULl default '',
    email VARCHAR(150) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forum(
    slug VARCHAR(100) NOT NULL primary key,
    title VARCHAR(100) NOT NULL,
    host VARCHAR(100) NOT NULL,
    posts bigint,
    threads int,
    foreign key (host) references actor (nickname)
        on DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS thread(
    id bigint primary key,
    title VARCHAR(300) NOT NULL,
    author VARCHAR(100) NOT NULL,
    forum VARCHAR(100),
    message TEXT NOT NULL,
    votes int,
    slug VARCHAR(150),
    created timestamp DEFAULT now(),
    foreign key (author) references actor (nickname)
        on DELETE CASCADE,
    foreign key (forum) references forum (slug)
        on DELETE CASCADE
);
