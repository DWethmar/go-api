package content

import (
	"database/sql"
	"log"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/database"
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

	getContentLocales = `
	SELECT locale FROM public.content_document
	WHERE content_id = $1
	`

	deleteContent = `
	DELETE FROM public.content WHERE id = $1
	`
)

// List entries.
func (repo *PostgresRepository) List() ([]*Content, error) {
	rows, err := repo.db.Query(allContent)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entries := make([]*Content, 0)

	for rows.Next() {
		entry := &Content{}
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

		// Postgres can only store timestamps in microsecond resolution.
		// Go supports nanosecond resolution and pq formats them in RFC3339Nano format, which includes nanoseconds.
		entry.UpdatedOn = common.DefaultTimePrecision(&entry.UpdatedOn)
		entry.CreatedOn = common.DefaultTimePrecision(&entry.CreatedOn)

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// Get one entry.
func (repo *PostgresRepository) Get(ID common.ID) (*Content, error) {
	entry := &Content{}
	row := repo.db.QueryRow(singleContent, ID)

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

	// Postgres can only store timestamps in microsecond resolution.
	// Go supports nanosecond resolution and pq formats them in RFC3339Nano format, which includes nanoseconds.
	entry.UpdatedOn = common.DefaultTimePrecision(&entry.UpdatedOn)
	entry.CreatedOn = common.DefaultTimePrecision(&entry.CreatedOn)

	return entry, nil
}

// Create add one entry.
func (repo *PostgresRepository) Create(c *Content) (common.ID, error) {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			insertContent,
			c.ID,
			c.Name,
			c.CreatedOn,
			c.UpdatedOn,
		)

		if err != nil {
			return err
		}

		for locale, fields := range c.Fields {
			_, err = tx.Exec(
				insertContentDocument,
				c.ID,
				locale,
				fields,
			)

			if err != nil {
				return err
			}
		}

		return err
	})

	return c.ID, err
}

func (repo *PostgresRepository) getLocales(ID common.ID) ([]string, error) {
	rows, err := repo.db.Query(getContentLocales, ID)

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
func (repo *PostgresRepository) Update(c *Content) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			updateContent,
			c.Name,
			c.UpdatedOn,
			c.ID,
		)

		if err != nil {
			return err
		}

		locales, err := repo.getLocales(c.ID)
		if err != nil {
			return err
		}

		for locale, fields := range c.Fields {
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
					c.ID,
					locale,
				)

				if err != nil {
					return err
				}
			} else {
				_, err = tx.Exec(
					insertContentDocument,
					c.ID,
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
func (repo *PostgresRepository) Delete(ID common.ID) error {
	r, err := repo.db.Exec(deleteContent, ID)
	if err != nil {
		return err
	}
	a, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if a == 0 {
		return ErrNotFound
	}
	return nil
}

// NewPostgresRepository create repo
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		db,
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
