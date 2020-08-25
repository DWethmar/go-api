package contenttype

import (
	"database/sql"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/models"
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

func (repo *PostgresRepository) getFields(contentModelID common.UUID) ([]*models.ContentTypeField, error) {
	rows, err := repo.db.Query(
		allContentTypeFields,
		contentModelID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entries := make([]*models.ContentTypeField, 0)

	for rows.Next() {
		entry := &models.ContentTypeField{}
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

// GetAll get all entries.
func (repo *PostgresRepository) GetAll() ([]*models.ContentType, error) {
	rows, err := repo.db.Query(allContentType)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entries := make([]*models.ContentType, 0)

	for rows.Next() {
		entry := &models.ContentType{}
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

// GetOne get one entry.
func (repo *PostgresRepository) GetOne(ID common.UUID) (*models.ContentType, error) {
	entry := &models.ContentType{}
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

// Add add one entry.
func (repo *PostgresRepository) Add(entry models.ContentType) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := tx.Exec(
			insertContentType,
			entry.ID,
			entry.Name,
			entry.CreatedOn,
			entry.UpdatedOn,
		)

		if err != nil {
			return err
		}

		for _, field := range entry.Fields {
			_, err = tx.Exec(
				insertContentTypeField,
				entry.ID,
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

	return err
}

// Update updates entry.
func (repo *PostgresRepository) Update(entry models.ContentType) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := tx.Exec(
			updateContent,
			entry.Name,
			entry.UpdatedOn,
			entry.ID,
		)

		if err != nil {
			return err
		}

		fields, err := repo.getFields(entry.ID)

		for _, field := range entry.Fields {
			exists := hasField(field.Key, fields)

			if exists {
				_, err := tx.Exec(
					updateContentTypeField,
					field.Name,
					field.Length,
					field.UpdatedOn,
					entry.ID,
				)

				if err != nil {
					return err
				}
			} else {
				_, err = tx.Exec(
					insertContentTypeField,
					entry.ID,
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
func (repo *PostgresRepository) Delete(ID common.UUID) error {
	_, err := repo.db.Exec(deleteContentType, ID)
	return err
}

func hasField(key string, fields []*models.ContentTypeField) bool {
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
