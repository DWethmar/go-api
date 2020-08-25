![Logo](/packages/my-app/src/krat_logo.png)

# go-api

Make sure $GOPATH/bin is in your $PATH so the shell can discover it.

- run: make run
- build: make build
- watch: make watch

# UI

- yarn install
- yarn start

//https://github.com/smhanov/auth


 structure based on: https://github.com/grafana/grafana/tree/master/pkg

 # TODO

 - fix empty post error
 - query builder
 - content: move name to subfield (stnd?)
 - contenttype: unique constraint on key + content_model_id.

 Versioning for content:

 table: content_model_version

 id (content id)
 version (int)
 content_model_id
 fields:
{
    <local>: {
        "<content_field_id>": "<value>"
    }
}
created_at
updated_at

delete field when content_field is deleted.

Versioning for content_models:
 thinking....