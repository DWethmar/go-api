CREATE TABLE public.content
(
    id VARCHAR (36) PRIMARY KEY,
    name VARCHAR (50),
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
);

CREATE TABLE public.content_fields_translation
(
    contenty_id VARCHAR (36) REFERENCES content(id) ON DELETE CASCADE,
    locale VARCHAR (50),
    fields jsonb NOT NULL,
    PRIMARY KEY (contenty_id, locale)
);

CREATE TABLE public.content_model
(
    id VARCHAR (36) PRIMARY KEY,
    name VARCHAR (50),
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
);

CREATE TABLE public.content_model_field
(
    content_model_id VARCHAR (36) REFERENCES content_model(id) ON DELETE CASCADE,
    key VARCHAR (50),  
    name VARCHAR (50),
    type VARCHAR (50),
    length smallint,
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
    PRIMARY KEY (content_model_id, field_key)
);

CREATE INDEX content_model_field_type ON content_model_field (type);