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

DROP DATABASE number_guess;
--
-- Name: number_guess; Type: DATABASE; Schema: -; Owner: freecodecamp
--

CREATE DATABASE number_guess WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C.UTF-8' LC_CTYPE = 'C.UTF-8';


ALTER DATABASE number_guess OWNER TO freecodecamp;

\connect number_guess

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
    guess integer DEFAULT 0 NOT NULL,
    player_id integer NOT NULL
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
-- Name: players; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.players (
    player_id integer NOT NULL,
    username character varying(22) NOT NULL
);


ALTER TABLE public.players OWNER TO freecodecamp;

--
-- Name: players_player_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.players_player_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.players_player_id_seq OWNER TO freecodecamp;

--
-- Name: players_player_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.players_player_id_seq OWNED BY public.players.player_id;


--
-- Name: games game_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.games ALTER COLUMN game_id SET DEFAULT nextval('public.games_game_id_seq'::regclass);


--
-- Name: players player_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.players ALTER COLUMN player_id SET DEFAULT nextval('public.players_player_id_seq'::regclass);


--
-- Data for Name: games; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.games VALUES (1, 2, 1);
INSERT INTO public.games VALUES (2, 1, 1);
INSERT INTO public.games VALUES (3, 1, 1);
INSERT INTO public.games VALUES (4, 344, 8);
INSERT INTO public.games VALUES (5, 863, 8);
INSERT INTO public.games VALUES (6, 414, 9);
INSERT INTO public.games VALUES (7, 497, 9);
INSERT INTO public.games VALUES (8, 409, 8);
INSERT INTO public.games VALUES (9, 578, 8);
INSERT INTO public.games VALUES (10, 394, 8);
INSERT INTO public.games VALUES (11, 2, 1);
INSERT INTO public.games VALUES (12, 3, 1);
INSERT INTO public.games VALUES (13, 6, 1);
INSERT INTO public.games VALUES (14, 3, 1);
INSERT INTO public.games VALUES (15, 1, 1);
INSERT INTO public.games VALUES (16, 1, 1);
INSERT INTO public.games VALUES (17, 1, 1);
INSERT INTO public.games VALUES (18, 396, 10);
INSERT INTO public.games VALUES (19, 181, 10);
INSERT INTO public.games VALUES (20, 102, 11);
INSERT INTO public.games VALUES (21, 64, 11);
INSERT INTO public.games VALUES (22, 708, 10);
INSERT INTO public.games VALUES (23, 270, 10);
INSERT INTO public.games VALUES (24, 602, 10);
INSERT INTO public.games VALUES (25, 364, 12);
INSERT INTO public.games VALUES (26, 893, 12);
INSERT INTO public.games VALUES (27, 78, 13);
INSERT INTO public.games VALUES (28, 65, 13);
INSERT INTO public.games VALUES (29, 78, 12);
INSERT INTO public.games VALUES (30, 167, 12);
INSERT INTO public.games VALUES (31, 155, 12);
INSERT INTO public.games VALUES (32, 289, 16);
INSERT INTO public.games VALUES (33, 350, 16);
INSERT INTO public.games VALUES (34, 346, 17);
INSERT INTO public.games VALUES (35, 188, 17);
INSERT INTO public.games VALUES (36, 211, 16);
INSERT INTO public.games VALUES (37, 45, 16);
INSERT INTO public.games VALUES (38, 154, 16);
INSERT INTO public.games VALUES (39, 952, 18);
INSERT INTO public.games VALUES (40, 145, 18);
INSERT INTO public.games VALUES (41, 302, 19);
INSERT INTO public.games VALUES (42, 710, 19);
INSERT INTO public.games VALUES (43, 913, 18);
INSERT INTO public.games VALUES (44, 206, 18);
INSERT INTO public.games VALUES (45, 306, 18);
INSERT INTO public.games VALUES (46, 361, 20);
INSERT INTO public.games VALUES (47, 487, 20);
INSERT INTO public.games VALUES (48, 712, 21);
INSERT INTO public.games VALUES (49, 599, 21);
INSERT INTO public.games VALUES (50, 372, 20);
INSERT INTO public.games VALUES (51, 130, 20);
INSERT INTO public.games VALUES (52, 109, 20);
INSERT INTO public.games VALUES (53, 692, 22);
INSERT INTO public.games VALUES (54, 268, 22);
INSERT INTO public.games VALUES (55, 954, 23);
INSERT INTO public.games VALUES (56, 469, 23);
INSERT INTO public.games VALUES (57, 111, 22);
INSERT INTO public.games VALUES (58, 889, 22);
INSERT INTO public.games VALUES (59, 410, 22);
INSERT INTO public.games VALUES (60, 26, 24);
INSERT INTO public.games VALUES (61, 899, 24);
INSERT INTO public.games VALUES (62, 937, 25);
INSERT INTO public.games VALUES (63, 987, 25);
INSERT INTO public.games VALUES (64, 217, 24);
INSERT INTO public.games VALUES (65, 465, 24);
INSERT INTO public.games VALUES (66, 649, 24);
INSERT INTO public.games VALUES (67, 847, 26);
INSERT INTO public.games VALUES (68, 475, 26);
INSERT INTO public.games VALUES (69, 532, 27);
INSERT INTO public.games VALUES (70, 81, 27);
INSERT INTO public.games VALUES (71, 538, 26);
INSERT INTO public.games VALUES (72, 313, 26);
INSERT INTO public.games VALUES (73, 81, 26);


--
-- Data for Name: players; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.players VALUES (1, 'Oleks');
INSERT INTO public.players VALUES (2, 'user_1672940869405');
INSERT INTO public.players VALUES (3, 'user_1672940869404');
INSERT INTO public.players VALUES (4, 'user_1672940926003');
INSERT INTO public.players VALUES (5, 'user_1672940926002');
INSERT INTO public.players VALUES (6, 'Vasya');
INSERT INTO public.players VALUES (7, 'Ia');
INSERT INTO public.players VALUES (8, 'user_1672942439865');
INSERT INTO public.players VALUES (9, 'user_1672942439864');
INSERT INTO public.players VALUES (10, 'user_1672942947030');
INSERT INTO public.players VALUES (11, 'user_1672942947029');
INSERT INTO public.players VALUES (12, 'user_1672942981238');
INSERT INTO public.players VALUES (13, 'user_1672942981237');
INSERT INTO public.players VALUES (14, 'Yurii');
INSERT INTO public.players VALUES (15, 'KLer');
INSERT INTO public.players VALUES (16, 'user_1672943077487');
INSERT INTO public.players VALUES (17, 'user_1672943077486');
INSERT INTO public.players VALUES (18, 'user_1672943125093');
INSERT INTO public.players VALUES (19, 'user_1672943125092');
INSERT INTO public.players VALUES (20, 'user_1672943574326');
INSERT INTO public.players VALUES (21, 'user_1672943574325');
INSERT INTO public.players VALUES (22, 'user_1672943621476');
INSERT INTO public.players VALUES (23, 'user_1672943621475');
INSERT INTO public.players VALUES (24, 'user_1672943744654');
INSERT INTO public.players VALUES (25, 'user_1672943744653');
INSERT INTO public.players VALUES (26, 'user_1672943888216');
INSERT INTO public.players VALUES (27, 'user_1672943888215');


--
-- Name: games_game_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.games_game_id_seq', 73, true);


--
-- Name: players_player_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.players_player_id_seq', 27, true);


--
-- Name: games games_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_pkey PRIMARY KEY (game_id);


--
-- Name: players players_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_pkey PRIMARY KEY (player_id);


--
-- Name: games games_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(player_id);


--
-- PostgreSQL database dump complete
--

