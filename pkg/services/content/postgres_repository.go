package content

import (
	"database/sql"
	"errors"
	"log"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/models"
)

// PostgresRepository repository for operating on content data.
type PostgresRepository struct {
	db *sql.DB
}

var (
	allContent = `
	SELECT 
		id, 
		name, 
		COALESCE(
			jsonb_object_agg(t.locale, t.fields) FILTER (WHERE t.locale IS NOT NULL), 
			'{}'::JSONB
		) as fields,
		created_on, 
		updated_on
	FROM public.content c
	LEFT OUTER JOIN public.content_document t ON c.id = t.content_id
	GROUP BY c.id
	ORDER BY created_on ASC`

	singleContent = `
	SELECT 
		id, 
		name, 
		COALESCE(
			jsonb_object_agg(t.locale, t.fields) FILTER (WHERE t.locale IS NOT NULL), 
			'{}'::JSONB
		) as fields,
		created_on, 
		updated_on
	FROM public.content c
	LEFT OUTER JOIN public.content_document t ON c.id = t.content_id
	WHERE c.id = $1
	GROUP BY c.id
	LIMIT 1`

	insertContent = `
	INSERT INTO public.content (id, name, created_on, updated_on)
	VALUES ($1, $2, $3, $4)`

	insertContentDocument = `
	INSERT INTO public.content_document(content_id, locale, fields) 
	VALUES($1, $2, $3)`

	updateContent = `
	UPDATE public.content SET (name, updated_on) = ($1, $2)
	WHERE id = $3`

	updateContentDocument = `
	UPDATE public.content_document SET fields = $1
	WHERE content_id = $2 AND locale = $3`

	getContentDocument = `
	SELECT locale FROM public.content_document
	WHERE content_id = $1
	`

	deleteContent = `
	DELETE FROM public.content WHERE id = $1
	`
)

// GetAll get all entries.
func (repo *PostgresRepository) GetAll() ([]*models.Content, error) {
	rows, err := repo.db.Query(allContent)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entrys := make([]*models.Content, 0)

	for rows.Next() {
		entry := &models.Content{}
		err := rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.Fields,
			&entry.CreatedOn,
			&entry.UpdatedOn,
		)
		if err != nil {
			return nil, err
		}
		entrys = append(entrys, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entrys, nil
}

// GetOne get one entry.
func (repo *PostgresRepository) GetOne(id common.UUID) (*models.Content, error) {
	entry := models.NewContent()
	row := repo.db.QueryRow(singleContent, id)

	var i string
	err := row.Scan(
		&i,
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

	entry.ID, err = common.ParseUUID(i)

	if err != nil {
		return nil, errors.New("Could not parse ID")
	}

	return entry, nil
}

// Add add one entry.
func (repo *PostgresRepository) Add(entry models.Content) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			insertContent,
			entry.ID,
			entry.Name,
			entry.CreatedOn,
			entry.UpdatedOn,
		)

		if err != nil {
			return err
		}

		for locale, fields := range entry.Fields {
			_, err = tx.Exec(
				insertContentDocument,
				entry.ID,
				locale,
				fields,
			)

			if err != nil {
				return err
			}

		}

		return err
	})

	return err
}

func (repo *PostgresRepository) getLocales(id common.UUID) ([]string, error) {
	rows, err := repo.db.Query(getContentDocument, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	locales := []string{}

	for rows.Next() {
		var locale string
		err := rows.Scan(&locale)
		if err != nil {
			return nil, err
		}
		locales = append(locales, locale)
	}

	return locales, nil
}

// Update updates entry.
func (repo *PostgresRepository) Update(entry models.Content) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			updateContent,
			entry.Name,
			entry.UpdatedOn,
			entry.ID,
		)

		if err != nil {
			return err
		}

		locales, err := repo.getLocales(entry.ID)
		if err != nil {
			return err
		}

		for locale, fields := range entry.Fields {
			hasLocale := false
			for _, n := range locales {
				if locale == n {
					hasLocale = true
				}
			}

			var err error
			if hasLocale {
				_, err = tx.Exec(
					updateContentDocument,
					fields,
					entry.ID,
					locale,
				)

				if err != nil {
					return err
				}
			} else {
				_, err = tx.Exec(
					insertContentDocument,
					entry.ID,
					locale,
					fields,
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
func (repo *PostgresRepository) Delete(id common.UUID) error {
	_, err := repo.db.Exec(deleteContent, id)
	return err
}

// CreatePostgresRepository create repo
func CreatePostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db,
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
