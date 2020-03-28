CREATE TABLE public.content_item
(
    id VARCHAR (36) PRIMARY KEY,
    name VARCHAR (50),
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
);

CREATE TABLE public.content_item_translation
(
    content_item_id VARCHAR (36) REFERENCES content_item(id),
    locale VARCHAR (50),
    attrs jsonb NOT NULL,
    PRIMARY KEY (content_item_id, locale)
);
