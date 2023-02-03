--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1
-- Dumped by pg_dump version 14.1

-- Started on 2023-02-03 14:30:59

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

--
-- TOC entry 3 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 3313 (class 0 OID 0)
-- Dependencies: 3
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 210 (class 1259 OID 98913)
-- Name: costumer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.costumer (
    id integer NOT NULL,
    costumer_name text,
    costumer_cont_no text,
    costumer_address text,
    total_buy text,
    creator_id text,
    date timestamp(0) without time zone
);


ALTER TABLE public.costumer OWNER TO postgres;

--
-- TOC entry 209 (class 1259 OID 98912)
-- Name: customer_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.customer_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.customer_id_seq OWNER TO postgres;

--
-- TOC entry 3314 (class 0 OID 0)
-- Dependencies: 209
-- Name: customer_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.customer_id_seq OWNED BY public.costumer.id;


--
-- TOC entry 3164 (class 2604 OID 98916)
-- Name: costumer id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.costumer ALTER COLUMN id SET DEFAULT nextval('public.customer_id_seq'::regclass);


--
-- TOC entry 3307 (class 0 OID 98913)
-- Dependencies: 210
-- Data for Name: costumer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.costumer (id, costumer_name, costumer_cont_no, costumer_address, total_buy, creator_id, date) FROM stdin;
1	fathur	09	depok	10	1	2022-05-18 13:50:19
2	rangga	08	depok	9	2	2023-01-22 22:33:17
3	radit	07	jaksel	8	3	2022-05-18 13:50:19
\.


--
-- TOC entry 3315 (class 0 OID 0)
-- Dependencies: 209
-- Name: customer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.customer_id_seq', 5, true);


--
-- TOC entry 3166 (class 2606 OID 98920)
-- Name: costumer customer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.costumer
    ADD CONSTRAINT customer_pkey PRIMARY KEY (id);


-- Completed on 2023-02-03 14:30:59

--
-- PostgreSQL database dump complete
--

