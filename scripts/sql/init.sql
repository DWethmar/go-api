CREATE TABLE public.content_item
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

CREATE TABLE public.content_item_translation
(
    content_item_id VARCHAR (36) REFERENCES content_item(id) ON DELETE CASCADE,
    locale VARCHAR (50),
    attrs jsonb NOT NULL,
    PRIMARY KEY (content_item_id, locale)
);
