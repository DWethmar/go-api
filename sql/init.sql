CREATE TABLE public.content_item
(
    id serial PRIMARY KEY,
    name VARCHAR (50),
    data jsonb
);

INSERT INTO public.content_item (name, data) VALUES ('test', '{ "test": "testing" }');
INSERT INTO public.content_item (name, data) VALUES ('test2', '{ "test": "testing" }');
