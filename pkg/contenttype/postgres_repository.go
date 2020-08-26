package contenttype

import (
	"database/sql"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/database"
)

// PostgresRepository repository for operating on content data.
type PostgresRepository struct {
	db *sql.DB
}

var (
	allContentType = `
	SELECT 
		id, 
		name, 
		created_on, 
		updated_on
	FROM public.content_type c
	ORDER BY created_on ASC`

	allContentTypeFields = `
	SELECT 
		id,
		content_model_id, 
		key,
		type,
		length, 
		created_on, 
		updated_on
	FROM public.content_type_field c
	WHERE c.content_model_id = $1
	ORDER BY created_on ASC`

	singleContentType = `
	SELECT 
		id, 
		name, 
		created_on, 
		updated_on
	FROM public.content_type c
	WHERE c.id = $1
	ORDER BY created_on ASC`

	insertContentType = `
	INSERT INTO public.content_type (id, name, created_on, updated_on)
	VALUES ($1, $2, $3, $4)`

	insertContentTypeField = `
	INSERT INTO public.content_type_field (id, content_model_id, key, name, type, length, created_on, updated_on)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	updateContent = `
	UPDATE public.content_type SET (name, updated_on) = ($1, $2)
	WHERE id = $3`

	updateContentTypeField = `
	UPDATE public.content_type_field (name, length, updated_on)
	VALUES ($1, $2, $3)
	WHERE id = $3`

	deleteContentType = `
	DELETE FROM public.content_type WHERE id = $1
	`
)

func (repo *PostgresRepository) getFields(contentModelID common.ID) ([]*Field, error) {
	rows, err := repo.db.Query(
		allContentTypeFields,
		contentModelID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entries := make([]*Field, 0)

	for rows.Next() {
		entry := &Field{}
		err := rows.Scan(
			&entry.EntryModelID,
			&entry.Key,
			&entry.Name,
			&entry.FieldType,
			&entry.Length,
			&entry.CreatedOn,
			&entry.UpdatedOn,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// List get all entries.
func (repo *PostgresRepository) List() ([]*ContentType, error) {
	rows, err := repo.db.Query(allContentType)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entries := make([]*ContentType, 0)

	for rows.Next() {
		entry := &ContentType{}
		err := rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.CreatedOn,
			&entry.UpdatedOn,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fields, err := repo.getFields(entry.ID)
		if err != nil {
			return nil, err
		}
		entry.Fields = fields
	}

	return entries, nil
}

// Get entry.
func (repo *PostgresRepository) Get(ID common.ID) (*ContentType, error) {
	entry := &ContentType{}
	row := repo.db.QueryRow(singleContentType, ID)

	err := row.Scan(
		&entry.ID,
		&entry.Name,
		&entry.Fields,
		&entry.CreatedOn,
		&entry.UpdatedOn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		panic(err)
	}

	fields, err := repo.getFields(entry.ID)
	if err != nil {
		return nil, err
	}
	entry.Fields = fields

	return entry, nil
}

// Create entry.
func (repo *PostgresRepository) Create(c *ContentType) (common.ID, error) {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := tx.Exec(
			insertContentType,
			c.ID,
			c.Name,
			c.CreatedOn,
			c.UpdatedOn,
		)

		if err != nil {
			return err
		}

		for _, field := range c.Fields {
			_, err = tx.Exec(
				insertContentTypeField,
				c.ID,
				field.Key,
				field.Name,
				field.Length,
				field.CreatedOn,
				field.UpdatedOn,
			)

			if err != nil {
				return err
			}
		}

		return err
	})

	return c.ID, err
}

// Update updates entry.
func (repo *PostgresRepository) Update(c *ContentType) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := tx.Exec(
			updateContent,
			c.Name,
			c.UpdatedOn,
			c.ID,
		)

		if err != nil {
			return err
		}

		fields, err := repo.getFields(c.ID)

		for _, field := range c.Fields {
			exists := hasField(field.Key, fields)

			if exists {
				_, err := tx.Exec(
					updateContentTypeField,
					field.Name,
					field.Length,
					field.UpdatedOn,
					c.ID,
				)

				if err != nil {
					return err
				}
			} else {
				_, err = tx.Exec(
					insertContentTypeField,
					c.ID,
					field.ID,
					field.Key,
					field.Name,
					field.Length,
					field.CreatedOn,
					field.UpdatedOn,
				)

				if err != nil {
					return err
				}
			}
		}

		return err
	})

	return err
}

// Delete deletes entry
func (repo *PostgresRepository) Delete(ID common.ID) error {
	_, err := repo.db.Exec(deleteContentType, ID)
	return err
}

func hasField(key string, fields []*Field) bool {
	for _, field := range fields {
		if key == field.Key {
			return true
		}
	}
	return false
}

// NewPostgresRepository create repo
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db,
	}
}
