CREATE TABLE public.content_item
(
    id serial PRIMARY KEY,
    name VARCHAR (50),
    data jsonb,
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
)

INSERT INTO public.content_item (name, data, created_on, updated_on) VALUES ('test', '{ "test": "testing" }', '2000-01-01 00:00:00 +00:00', '2000-01-01 00:00:00 +00:00');
INSERT INTO public.content_item (name, data, created_on, updated_on) VALUES ('test2', '{ "test": "testing" }', '2000-01-01 00:00:00 +00:00', '2000-01-01 00:00:00 +00:00');
