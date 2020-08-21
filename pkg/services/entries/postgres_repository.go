package entries

import (
	"database/sql"
	"errors"
	"log"

	"github.com/dwethmar/go-api/pkg/common"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/models"
)

// PostgresRepository repository for operating on entry data.
type PostgresRepository struct {
	db *sql.DB
}

var (
	getAll = `
	SELECT 
		id, 
		name, 
		COALESCE(
			jsonb_object_agg(t.locale, t.fields) FILTER (WHERE t.locale IS NOT NULL), 
			'{}'::JSONB
		) as fields,
		created_on, 
		updated_on
	FROM public.entry c
	LEFT OUTER JOIN public.entry_translation t ON c.id = t.entry_id
	GROUP BY c.id
	ORDER BY updated_on ASC`

	getOne = `
	SELECT 
		id, 
		name, 
		COALESCE(
			jsonb_object_agg(t.locale, t.fields) FILTER (WHERE t.locale IS NOT NULL), 
			'{}'::JSONB
		) as fields,
		created_on, 
		updated_on
	FROM public.entry c
	LEFT OUTER JOIN public.entry_translation t ON c.id = t.entry_id
	WHERE c.id = $1
	GROUP BY c.id
	LIMIT 1`

	insertEntry = `
	INSERT INTO public.entry (id, name, created_on, updated_on)
	VALUES ($1, $2, $3, $4)`

	insertEntryTrans = `
	INSERT INTO public.entry_translation(entry_id, locale, fields) 
	VALUES($1, $2, $3)`

	updateEntry = `
	UPDATE public.entry SET (name, updated_on) = ($1, $2)
	WHERE id = $3`

	updateEntryTrans = `
	UPDATE public.entry_translation SET fields = $1
	WHERE entry_id = $2 AND locale = $3`

	getEntryLocales = `
	SELECT locale FROM public.entry_translation
	WHERE entry_id = $1
	`

	deleteentry = `
	DELETE FROM public.entry WHERE id = $1
	`
)

// GetAll get all entries.
func (repo *PostgresRepository) GetAll() ([]*models.Entry, error) {
	rows, err := repo.db.Query(getAll)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entrys := make([]*models.Entry, 0)

	for rows.Next() {
		entry := &models.Entry{}
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
func (repo *PostgresRepository) GetOne(id common.UUID) (*models.Entry, error) {
	entry := models.CreateEntry()
	row := repo.db.QueryRow(getOne, id)

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

	return &entry, nil
}

// Add add one entry.
func (repo *PostgresRepository) Add(entry models.Entry) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			insertEntry,
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
				insertEntryTrans,
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
	rows, err := repo.db.Query(getEntryLocales, id)

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
func (repo *PostgresRepository) Update(entry models.Entry) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			updateEntry,
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
					updateEntryTrans,
					fields,
					entry.ID,
					locale,
				)

				if err != nil {
					return err
				}
			} else {
				_, err = tx.Exec(
					insertEntryTrans,
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
	_, err := repo.db.Exec(deleteentry, id)
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
