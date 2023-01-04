--
-- PostgreSQL database dump
--

-- Dumped from database version 12.9 (Ubuntu 12.9-2.pgdg20.04+1)
-- Dumped by pg_dump version 12.9 (Ubuntu 12.9-2.pgdg20.04+1)

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

DROP DATABASE worldcup;
--
-- Name: worldcup; Type: DATABASE; Schema: -; Owner: freecodecamp
--

CREATE DATABASE worldcup WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C.UTF-8' LC_CTYPE = 'C.UTF-8';


ALTER DATABASE worldcup OWNER TO freecodecamp;

\connect worldcup

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

--
-- Name: games; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.games (
    game_id integer NOT NULL,
    year integer NOT NULL,
    round character varying(30) NOT NULL,
    winner_goals integer NOT NULL,
    opponent_goals integer NOT NULL,
    winner_id integer NOT NULL,
    opponent_id integer NOT NULL
);


ALTER TABLE public.games OWNER TO freecodecamp;

--
-- Name: games_game_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.games_game_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.games_game_id_seq OWNER TO freecodecamp;

--
-- Name: games_game_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.games_game_id_seq OWNED BY public.games.game_id;


--
-- Name: teams; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.teams (
    team_id integer NOT NULL,
    name character varying(30) NOT NULL
);


ALTER TABLE public.teams OWNER TO freecodecamp;

--
-- Name: teams_team_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.teams_team_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.teams_team_id_seq OWNER TO freecodecamp;

--
-- Name: teams_team_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.teams_team_id_seq OWNED BY public.teams.team_id;


--
-- Name: games game_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.games ALTER COLUMN game_id SET DEFAULT nextval('public.games_game_id_seq'::regclass);


--
-- Name: teams team_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.teams ALTER COLUMN team_id SET DEFAULT nextval('public.teams_team_id_seq'::regclass);


--
-- Data for Name: games; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.games VALUES (5, 2018, 'Final', 4, 2, 22, 23);
INSERT INTO public.games VALUES (6, 2018, 'Third Place', 2, 0, 24, 25);
INSERT INTO public.games VALUES (7, 2018, 'Semi-Final', 2, 1, 23, 25);
INSERT INTO public.games VALUES (8, 2018, 'Semi-Final', 1, 0, 22, 24);
INSERT INTO public.games VALUES (9, 2018, 'Quarter-Final', 3, 2, 23, 26);
INSERT INTO public.games VALUES (10, 2018, 'Quarter-Final', 2, 0, 25, 27);
INSERT INTO public.games VALUES (11, 2018, 'Quarter-Final', 2, 1, 24, 28);
INSERT INTO public.games VALUES (12, 2018, 'Quarter-Final', 2, 0, 22, 29);
INSERT INTO public.games VALUES (13, 2018, 'Eighth-Final', 2, 1, 25, 30);
INSERT INTO public.games VALUES (14, 2018, 'Eighth-Final', 1, 0, 27, 31);
INSERT INTO public.games VALUES (15, 2018, 'Eighth-Final', 3, 2, 24, 32);
INSERT INTO public.games VALUES (16, 2018, 'Eighth-Final', 2, 0, 28, 33);
INSERT INTO public.games VALUES (17, 2018, 'Eighth-Final', 2, 1, 23, 34);
INSERT INTO public.games VALUES (18, 2018, 'Eighth-Final', 2, 1, 26, 35);
INSERT INTO public.games VALUES (19, 2018, 'Eighth-Final', 2, 1, 29, 36);
INSERT INTO public.games VALUES (20, 2018, 'Eighth-Final', 4, 3, 22, 37);
INSERT INTO public.games VALUES (21, 2014, 'Final', 1, 0, 38, 37);
INSERT INTO public.games VALUES (22, 2014, 'Third Place', 3, 0, 39, 28);
INSERT INTO public.games VALUES (23, 2014, 'Semi-Final', 1, 0, 37, 39);
INSERT INTO public.games VALUES (24, 2014, 'Semi-Final', 7, 1, 38, 28);
INSERT INTO public.games VALUES (25, 2014, 'Quarter-Final', 1, 0, 39, 40);
INSERT INTO public.games VALUES (26, 2014, 'Quarter-Final', 1, 0, 37, 24);
INSERT INTO public.games VALUES (27, 2014, 'Quarter-Final', 2, 1, 28, 30);
INSERT INTO public.games VALUES (28, 2014, 'Quarter-Final', 1, 0, 38, 22);
INSERT INTO public.games VALUES (29, 2014, 'Eighth-Final', 2, 1, 28, 41);
INSERT INTO public.games VALUES (30, 2014, 'Eighth-Final', 2, 0, 30, 29);
INSERT INTO public.games VALUES (31, 2014, 'Eighth-Final', 2, 0, 22, 42);
INSERT INTO public.games VALUES (32, 2014, 'Eighth-Final', 2, 1, 38, 43);
INSERT INTO public.games VALUES (33, 2014, 'Eighth-Final', 2, 1, 39, 33);
INSERT INTO public.games VALUES (34, 2014, 'Eighth-Final', 2, 1, 40, 44);
INSERT INTO public.games VALUES (35, 2014, 'Eighth-Final', 1, 0, 37, 31);
INSERT INTO public.games VALUES (36, 2014, 'Eighth-Final', 2, 1, 24, 45);


--
-- Data for Name: teams; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.teams VALUES (22, 'France');
INSERT INTO public.teams VALUES (23, 'Croatia');
INSERT INTO public.teams VALUES (24, 'Belgium');
INSERT INTO public.teams VALUES (25, 'England');
INSERT INTO public.teams VALUES (26, 'Russia');
INSERT INTO public.teams VALUES (27, 'Sweden');
INSERT INTO public.teams VALUES (28, 'Brazil');
INSERT INTO public.teams VALUES (29, 'Uruguay');
INSERT INTO public.teams VALUES (30, 'Colombia');
INSERT INTO public.teams VALUES (31, 'Switzerland');
INSERT INTO public.teams VALUES (32, 'Japan');
INSERT INTO public.teams VALUES (33, 'Mexico');
INSERT INTO public.teams VALUES (34, 'Denmark');
INSERT INTO public.teams VALUES (35, 'Spain');
INSERT INTO public.teams VALUES (36, 'Portugal');
INSERT INTO public.teams VALUES (37, 'Argentina');
INSERT INTO public.teams VALUES (38, 'Germany');
INSERT INTO public.teams VALUES (39, 'Netherlands');
INSERT INTO public.teams VALUES (40, 'Costa Rica');
INSERT INTO public.teams VALUES (41, 'Chile');
INSERT INTO public.teams VALUES (42, 'Nigeria');
INSERT INTO public.teams VALUES (43, 'Algeria');
INSERT INTO public.teams VALUES (44, 'Greece');
INSERT INTO public.teams VALUES (45, 'United States');


--
-- Name: games_game_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.games_game_id_seq', 36, true);


--
-- Name: teams_team_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.teams_team_id_seq', 45, true);


--
-- Name: games games_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_pkey PRIMARY KEY (game_id);


--
-- Name: teams teams_name_key; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.teams
    ADD CONSTRAINT teams_name_key UNIQUE (name);


--
-- Name: teams teams_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.teams
    ADD CONSTRAINT teams_pkey PRIMARY KEY (team_id);


--
-- Name: games opponent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT opponent_id_fkey FOREIGN KEY (opponent_id) REFERENCES public.teams(team_id);


--
-- Name: games winner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT winner_id_fkey FOREIGN KEY (winner_id) REFERENCES public.teams(team_id);


--
-- PostgreSQL database dump complete
--

