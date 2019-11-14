CREATE DATABASE "go-api"
    WITH 
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

CREATE TABLE public.data_node
(
    id integer NOT NULL DEFAULT nextval('data_node_id_seq'::regclass),
    name character(128) COLLATE pg_catalog."default",
    data jsonb,
    CONSTRAINT data_node_pkey PRIMARY KEY (id)
)

INSERT INTO public.data_node (name, data) VALUES ('test', '{ "test": "testing" }')
