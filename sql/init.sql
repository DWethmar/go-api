CREATE TABLE public.data_node
(
    id serial PRIMARY KEY,
    name VARCHAR (50),
    data jsonb
);

INSERT INTO public.data_node (name, data) VALUES ('test', '{ "test": "testing" }');
INSERT INTO public.data_node (name, data) VALUES ('test2', '{ "test": "testing" }');
