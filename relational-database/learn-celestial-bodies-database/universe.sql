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

DROP DATABASE universe;
--
-- Name: universe; Type: DATABASE; Schema: -; Owner: freecodecamp
--

CREATE DATABASE universe WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C.UTF-8' LC_CTYPE = 'C.UTF-8';


ALTER DATABASE universe OWNER TO freecodecamp;

\connect universe

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
-- Name: galaxy; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.galaxy (
    galaxy_id integer NOT NULL,
    name character varying(30),
    description text,
    is_spherical boolean NOT NULL,
    age_in_millions_of_years integer NOT NULL,
    stars_count numeric
);


ALTER TABLE public.galaxy OWNER TO freecodecamp;

--
-- Name: galaxy_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.galaxy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.galaxy_id_seq OWNER TO freecodecamp;

--
-- Name: galaxy_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.galaxy_id_seq OWNED BY public.galaxy.galaxy_id;


--
-- Name: galaxy_types; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.galaxy_types (
    galaxy_types_id integer NOT NULL,
    name character varying(30) NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.galaxy_types OWNER TO freecodecamp;

--
-- Name: galaxy_types_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.galaxy_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.galaxy_types_id_seq OWNER TO freecodecamp;

--
-- Name: galaxy_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.galaxy_types_id_seq OWNED BY public.galaxy_types.galaxy_types_id;


--
-- Name: moon; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.moon (
    moon_id integer NOT NULL,
    name character varying(30) NOT NULL,
    description text NOT NULL,
    distance_from_earth integer,
    planet_id integer
);


ALTER TABLE public.moon OWNER TO freecodecamp;

--
-- Name: moon_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.moon_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.moon_id_seq OWNER TO freecodecamp;

--
-- Name: moon_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.moon_id_seq OWNED BY public.moon.moon_id;


--
-- Name: planet; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.planet (
    planet_id integer NOT NULL,
    name character varying(30) NOT NULL,
    description text NOT NULL,
    distance_from_earth integer,
    has_life boolean,
    star_id integer
);


ALTER TABLE public.planet OWNER TO freecodecamp;

--
-- Name: planet_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.planet_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.planet_id_seq OWNER TO freecodecamp;

--
-- Name: planet_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.planet_id_seq OWNED BY public.planet.planet_id;


--
-- Name: star; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.star (
    star_id integer NOT NULL,
    name character varying(30) NOT NULL,
    description text NOT NULL,
    distance_from_earth integer,
    galaxy_id integer
);


ALTER TABLE public.star OWNER TO freecodecamp;

--
-- Name: star_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.star_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.star_id_seq OWNER TO freecodecamp;

--
-- Name: star_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.star_id_seq OWNED BY public.star.star_id;


--
-- Name: galaxy galaxy_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy ALTER COLUMN galaxy_id SET DEFAULT nextval('public.galaxy_id_seq'::regclass);


--
-- Name: galaxy_types galaxy_types_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy_types ALTER COLUMN galaxy_types_id SET DEFAULT nextval('public.galaxy_types_id_seq'::regclass);


--
-- Name: moon moon_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon ALTER COLUMN moon_id SET DEFAULT nextval('public.moon_id_seq'::regclass);


--
-- Name: planet planet_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet ALTER COLUMN planet_id SET DEFAULT nextval('public.planet_id_seq'::regclass);


--
-- Name: star star_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star ALTER COLUMN star_id SET DEFAULT nextval('public.star_id_seq'::regclass);


--
-- Data for Name: galaxy; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.galaxy VALUES (1, 'Andromeda', 'The Andromeda Galaxy, also known as Messier 31, M31, or NGC 224 and originally the Andromeda Nebula, is a barred spiral galaxy with the diameter of about 46.56 kiloparsecs approximately 2.5 million light-years from Earth and the nearest large galaxy to the Milky Way.', false, 10, 55343);
INSERT INTO public.galaxy VALUES (3, 'Messier 81', 'Messier 81 is a grand design spiral galaxy about 12 million light-years away in the constellation Ursa Major.', true, 12, 9234);
INSERT INTO public.galaxy VALUES (4, 'Centaurus', 'Centaurus A is a galaxy in the constellation of Centaurus. It was discovered in 1826 by Scottish astronomer James Dunlop from his home in Parramatta, in New South Wales, Australia.', false, 15, 432412);
INSERT INTO public.galaxy VALUES (5, 'Sculptor', 'The Sculptor Galaxy is an intermediate spiral galaxy in the constellation Sculptor. The Sculptor Galaxy is a starburst galaxy, which means that it is currently undergoing a period of intense star formation.', false, 9, 1923);
INSERT INTO public.galaxy VALUES (6, 'Milky Way', 'The Milky Way is the galaxy that includes our Solar System, with the name describing the galaxy appearance from Earth', false, 5, 973);
INSERT INTO public.galaxy VALUES (7, 'Triangulum', 'The Triangulum Galaxy is a spiral galaxy 2.73 million light-years from Earth in the constellation Triangulum. It is catalogued as Messier 33 or NGC 598.', true, 6, 400);


--
-- Data for Name: galaxy_types; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.galaxy_types VALUES (1, 'elliptical', 'An elliptical galaxy is a type of galaxy with an approximately ellipsoidal shape and a smooth, nearly featureless image.');
INSERT INTO public.galaxy_types VALUES (2, 'shell', 'A shell galaxy is a type of elliptical galaxy where the stars in its halo are arranged in concentric shells.');
INSERT INTO public.galaxy_types VALUES (3, 'spiral', 'Spiral galaxies form a class of galaxy originally described by Edwin Hubble in his 1936 work.');


--
-- Data for Name: moon; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.moon VALUES (2, 'Moon', 'The Moon is Earth only natural satellite. It is the fifth largest satellite in the Solar System and the largest and most massive relative to its parent planet, with a diameter about one-quarter that of Earth.', 23432, 2);
INSERT INTO public.moon VALUES (3, 'Titan', 'Titan is the largest moon of Saturn and the second-largest natural satellite in the Solar System. It is the only moon known to have a dense atmosphere, and is the only known object in space other than Earth on which clear evidence of stable bodies of surface liquid has been found.', 324252, 5);
INSERT INTO public.moon VALUES (4, 'Europa', 'Europa, or Jupiter II, is the smallest of the four Galilean moons orbiting Jupiter, and the sixth-closest to the planet of all the 80 known moons of Jupiter. It is also the sixth-largest moon in the Solar System.', 342342, 3);
INSERT INTO public.moon VALUES (5, 'Ganymede', 'Ganymede, a satellite of Jupiter, is the largest and most massive of the Solar System moons. The ninth-largest object of the Solar System, it is the largest without a substantial atmosphere.', 2134, 3);
INSERT INTO public.moon VALUES (6, 'Callisto', 'Callisto, or Jupiter IV, is the second-largest moon of Jupiter, after Ganymede. It is the third-largest moon in the Solar System after Ganymede and Saturn largest moon Titan, and the largest object in the Solar System that may not be properly differentiated. Callisto was discovered in 1610 by Galileo Galilei.', 92341, 3);
INSERT INTO public.moon VALUES (7, 'Io', 'Io, or Jupiter I, is the innermost and third-largest of the four Galilean moons of the planet Jupiter.', 4321, 3);
INSERT INTO public.moon VALUES (8, 'Amalthea', 'Amalthea is a moon of Jupiter. It has the third closest orbit around Jupiter among known moons and was the fifth moon of Jupiter to be discovered, so it is also known as Jupiter V. It is also the fifth largest moon of Jupiter, after the four Galilean Moons.', 234151, 5);
INSERT INTO public.moon VALUES (9, 'Himalia', 'Himalia, or Jupiter VI, is the largest irregular satellite of Jupiter, with a diameter of at least 140 km. It is the sixth largest Jovian satellite, after the four Galilean moons and Amalthea.', 93241, 3);
INSERT INTO public.moon VALUES (10, 'Mimas', 'Mimas, also designated Saturn I, is a moon of Saturn discovered in 1789 by William Herschel. It is named after Mimas, a son of Gaia in Greek mythology.', 231451, 7);
INSERT INTO public.moon VALUES (11, 'Enceladus', 'Enceladus is the sixth-largest moon of Saturn. It is about 500 kilometers in diameter, about a tenth of that of Saturn largest moon, Titan. Enceladus is mostly covered by fresh, clean ice, making it one of the most reflective bodies of the Solar System.', 9234, 3);
INSERT INTO public.moon VALUES (12, 'Hyperion', 'Hyperion, also known as Saturn VII, is a moon of Saturn discovered by William Cranch Bond, his son George Phillips Bond and William Lassell in 1848. It is distinguished by its irregular shape, its chaotic rotation, and its unexplained sponge-like appearance. It was the first non-round moon to be discovered.', 4321, 3);
INSERT INTO public.moon VALUES (13, 'Iapetus', 'Iapetus is a moon of Saturn. It is the 24th of Saturn 83 known moons. With an estimated diameter of 1,469 km, it is the third-largest moon of Saturn and the eleventh-largest in the Solar System. Named after the Titan Iapetus, the moon was discovered in 1671 by Giovanni Domenico Cassini.', 934, 3);
INSERT INTO public.moon VALUES (14, 'Adrastea', 'Adrastea, also known as Jupiter XV, is the second by distance, and the smallest of the four inner moons of Jupiter.', 3412, 10);
INSERT INTO public.moon VALUES (15, 'Dione', 'Dione is a moon of Saturn. It was discovered by Italian astronomer Giovanni Domenico Cassini in 1684. It is named after the Titaness Dione of Greek mythology.', 4321785, 11);
INSERT INTO public.moon VALUES (16, 'Thebe', 'Thebe, also known as Jupiter XIV, is the fourth of Jupiter moons by distance from the planet. It was discovered by Stephen P. Synnott in images from the Voyager 1 space probe taken on March 5, 1979, while making its flyby of Jupiter.', 84312, 12);
INSERT INTO public.moon VALUES (17, 'Epimetheus', 'Epimetheus is an inner satellite of Saturn. It is also known as Saturn XI. It is named after the mythological Epimetheus, brother of Prometheus.', 4321, 3);
INSERT INTO public.moon VALUES (18, 'Tethys', 'Tethys, or Saturn III, is a mid-sized moon of Saturn about 1,060 km across. It was discovered by G. D. Cassini in 1684 and is named after the titan Tethys of Greek mythology.', 7563, 4);
INSERT INTO public.moon VALUES (19, 'Carpo', 'Carpo, also Jupiter XLVI, is a natural satellite of Jupiter. It was discovered by a team of astronomers from the University of Hawaii led by Scott S. Sheppard in 2003, and was provisionally designated as S/2003 J 20 until it received its name in early 2005.', 3412412, 4);
INSERT INTO public.moon VALUES (20, 'Ananke', 'Ananke is a retrograde irregular moon of Jupiter. It was discovered by Seth Barnes Nicholson at Mount Wilson Observatory in 1951 and is named after the Greek mythological Ananke, the personification of necessity, and the mother of the Moirai by Zeus. The adjectival form of the name is Anankean.', 41320, 4);
INSERT INTO public.moon VALUES (21, 'Iocaste', 'Iocaste, also known as Jupiter XXIV, is a retrograde irregular satellite of Jupiter. It was discovered by a team of astronomers from the University of Hawaii including: David C. Jewitt, Yanga R. Fernandez, and Eugene Magnier led by Scott S. Sheppard in 2000, and given the temporary designation S/2000 J 3.', 5830, 13);


--
-- Data for Name: planet; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.planet VALUES (2, 'Earth', 'Earth is the third planet from the Sun and the only place in the universe known to harbor life. While large volumes of water can be found throughout the Solar System, only Earth sustains liquid surface water. About 71% of Earth surface is made up of the ocean, dwarfing Earth polar ice, lakes, and rivers.', NULL, true, 7);
INSERT INTO public.planet VALUES (3, 'Jupiter', 'Jupiter is the fifth planet from the Sun and the largest in the Solar System. It is a gas giant with a mass more than two and a half times that of all the other planets in the Solar System combined, while being slightly less than one-thousandth the mass of the Sun.', 234, false, 8);
INSERT INTO public.planet VALUES (4, 'Saturn', 'Saturn is the sixth planet from the Sun and the second-largest in the Solar System, after Jupiter. It is a gas giant with an average radius of about nine and a half times that of Earth. It has only one-eighth the average density of Earth, but is over 95 times more massive.', 12, false, 8);
INSERT INTO public.planet VALUES (5, 'Mercury', 'Mercury is the smallest planet in the Solar System and the closest to the Sun. Its orbit around the Sun takes 87.97 Earth days, the shortest of all the Sun planets.', 28, false, 8);
INSERT INTO public.planet VALUES (6, 'Mars', 'Mars is the fourth planet from the Sun and the second-smallest planet in the Solar System, only being larger than Mercury. In the English language, Mars is named for the Roman god of war.', 23, false, 8);
INSERT INTO public.planet VALUES (7, 'Venus', 'Venus is the second planet from the Sun. It is sometimes called Earth "sister" or "twin" planet as it is almost as large and has a similar composition. As an interior planet to Earth, Venus appears in Earth sky never far from the Sun, either as morning star or evening star.', 34, false, 9);
INSERT INTO public.planet VALUES (8, 'Uranus', 'Uranus is the seventh planet from the Sun. It is named after Greek sky deity Uranus, who in Greek mythology is the father of Cronus, a grandfather of Zeus and great-grandfather of Ares. Uranus has the third-largest planetary radius and fourth-largest planetary mass in the Solar System.', 53, false, 10);
INSERT INTO public.planet VALUES (9, 'Neptune', 'Neptune is the eighth planet from the Sun and the farthest known planet in the Solar System. It is the fourth-largest planet in the Solar System by diameter, the third-most-massive planet, and the densest giant planet. It is 17 times the mass of Earth, and slightly more massive than its near-twin Uranus.', 19, false, 12);
INSERT INTO public.planet VALUES (10, 'Arion', 'Arion was a genius of poetry and music in ancient Greece. According to legend, his life was saved at sea by dolphins after attracting their attention by the playing of his kithara.', 234, false, 12);
INSERT INTO public.planet VALUES (11, 'Orbitar', 'Orbitar is a contrived word paying homage to the space launch and orbital operations of NASA.', 53, false, 11);
INSERT INTO public.planet VALUES (12, 'Taphao Thong', 'Taphao Thong is one of two sisters associated with the Thai folk tale of Chalawan.', 93, false, 11);
INSERT INTO public.planet VALUES (13, 'Dimidium', 'Dimidium is Latin for half, referring to the planet mass of at least half the mass of Jupiter.', 111, false, 12);


--
-- Data for Name: star; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.star VALUES (7, 'Sirius', 'Sirius is the brightest star in the night sky. Its name is derived from the Greek word Σείριος, or Seirios, meaning lit. glowing or scorching. The star is designated α Canis Majoris, Latinized to Alpha Canis Majoris, and abbreviated Alpha CMa or α CMa', 2324, 1);
INSERT INTO public.star VALUES (8, 'Betelgeuse', 'Betelgeuse is a red supergiant of spectral type M1-2 and one of the largest stars visible to the naked eye. It is usually the tenth-brightest star in the night sky and, after Rigel, the second-brightest in the constellation of Orion.', 342, 7);
INSERT INTO public.star VALUES (9, 'Vega', 'Vega is the brightest star in the northern constellation of Lyra. It has the Bayer designation α Lyrae, which is Latinised to Alpha Lyrae and abbreviated Alpha Lyr or α Lyr. This star is relatively close at only 25 light-years from the Sun, and one of the most luminous stars in the Sun neighborhood.', 9432, 3);
INSERT INTO public.star VALUES (10, 'Alpha Persei', 'Alpha Persei, formally named Mirfak, is the brightest star in the northern constellation of Perseus, outshining the constellation best-known star, Algol. Alpha Persei has an apparent visual magnitude of 1.8, and is a circumpolar star when viewed from mid-northern latitudes.', 3421, 4);
INSERT INTO public.star VALUES (11, 'Delta Canis Majoris', 'Delta Canis Majoris, officially named Wezen, is a star in the constellation of Canis Major. It is a yellow-white F-type supergiant with an apparent magnitude of +1.83. Since 1943, the spectrum of this star has served as one of the stable anchor points by which other stars are classified.', 84, 5);
INSERT INTO public.star VALUES (12, 'Fomalhaut', 'Fomalhaut is the brightest star in the southern constellation of Piscis Austrinus, the "Southern Fish", and one of the brightest stars in the night sky. It has the Bayer designation Alpha Piscis Austrini, which is Latinized from α Piscis Austrini, and is abbreviated Alpha PsA or α PsA.', 943, 6);


--
-- Name: galaxy_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.galaxy_id_seq', 7, true);


--
-- Name: galaxy_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.galaxy_types_id_seq', 3, true);


--
-- Name: moon_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.moon_id_seq', 21, true);


--
-- Name: planet_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.planet_id_seq', 13, true);


--
-- Name: star_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.star_id_seq', 12, true);


--
-- Name: galaxy galaxy_name_unique; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy
    ADD CONSTRAINT galaxy_name_unique UNIQUE (name);


--
-- Name: galaxy galaxy_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy
    ADD CONSTRAINT galaxy_pkey PRIMARY KEY (galaxy_id);


--
-- Name: galaxy_types galaxy_types_name; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy_types
    ADD CONSTRAINT galaxy_types_name UNIQUE (name);


--
-- Name: galaxy_types galaxy_types_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy_types
    ADD CONSTRAINT galaxy_types_pkey PRIMARY KEY (galaxy_types_id);


--
-- Name: moon moon_name_unique; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon
    ADD CONSTRAINT moon_name_unique UNIQUE (name);


--
-- Name: moon moon_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon
    ADD CONSTRAINT moon_pkey PRIMARY KEY (moon_id);


--
-- Name: planet planet_name_unique; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet
    ADD CONSTRAINT planet_name_unique UNIQUE (name);


--
-- Name: planet planet_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet
    ADD CONSTRAINT planet_pkey PRIMARY KEY (planet_id);


--
-- Name: star star_name_unique; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star
    ADD CONSTRAINT star_name_unique UNIQUE (name);


--
-- Name: star star_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star
    ADD CONSTRAINT star_pkey PRIMARY KEY (star_id);


--
-- Name: moon moon_planet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon
    ADD CONSTRAINT moon_planet_id_fkey FOREIGN KEY (planet_id) REFERENCES public.planet(planet_id);


--
-- Name: planet planet_star_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet
    ADD CONSTRAINT planet_star_id_fkey FOREIGN KEY (star_id) REFERENCES public.star(star_id);


--
-- Name: star star_galaxy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star
    ADD CONSTRAINT star_galaxy_id_fkey FOREIGN KEY (galaxy_id) REFERENCES public.galaxy(galaxy_id);


--
-- PostgreSQL database dump complete
--
