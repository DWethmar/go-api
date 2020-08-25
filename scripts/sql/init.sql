
CREATE TABLE public.content_model
(
    id VARCHAR (36) PRIMARY KEY,
    name VARCHAR (50),
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
);

CREATE TABLE public.content_model_field
(
    id VARCHAR (36) PRIMARY KEY,
    content_model_id VARCHAR (36) REFERENCES content_model(id) ON DELETE CASCADE,
    key VARCHAR (50),  
    name VARCHAR (50),
    type VARCHAR (50),
    length smallint,
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL,
    UNIQUE (content_model_id, key)
);

CREATE INDEX content_model_field_type ON content_model_field(type);

CREATE TABLE public.content
(
    id VARCHAR (36) PRIMARY KEY,
    content_model_id VARCHAR (36) REFERENCES content_model(id) ON DELETE CASCADE,
    name VARCHAR (50),
    content_model_id 
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
);

CREATE TABLE public.content_document
(
    content_id VARCHAR (36) REFERENCES content(id) ON DELETE CASCADE,
    locale VARCHAR (50),
    fields jsonb NOT NULL,
    PRIMARY KEY (content_id, locale)
);
