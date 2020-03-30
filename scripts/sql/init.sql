CREATE TABLE public.content_entry
(
    id VARCHAR (36) PRIMARY KEY,
    name VARCHAR (50),
    created_on timestamptz NOT NULL,
    updated_on timestamptz NOT NULL
    /* 
    Added a version number would be nice.
    version interger NOT NULL

    And with each update increate it. 
    example:  SET version = version + 1
    **/
);

CREATE TABLE public.content_entry_translation
(
    content_entry_id VARCHAR (36) REFERENCES content_entry(id) ON DELETE CASCADE,
    locale VARCHAR (50),
    fields jsonb NOT NULL,
    PRIMARY KEY (content_entry_id, locale)
);
