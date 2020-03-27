CREATE TABLE public.content_item
(
    id VARCHAR (36) PRIMARY KEY,
    name VARCHAR (50),
    attrs jsonb NOT NULL,
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
);
