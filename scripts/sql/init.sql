CREATE TABLE public.entry
(
    id VARCHAR (36) PRIMARY KEY,
    name VARCHAR (50),
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
);

CREATE TABLE public.entry_translation
(
    entry_id VARCHAR (36) REFERENCES entry(id) ON DELETE CASCADE,
    locale VARCHAR (50),
    fields jsonb NOT NULL,
    PRIMARY KEY (entry_id, locale)
);


