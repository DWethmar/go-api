package contententry

import (
	"database/sql"
	"errors"
	"log"

	"github.com/dwethmar/go-api/pkg/database"
)

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
	FROM public.content_entry c
	LEFT OUTER JOIN public.content_entry_field t ON c.id = t.content_entry_id
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
	FROM public.content_entry c
	LEFT OUTER JOIN public.content_entry_field t ON c.id = t.content_entry_id
	WHERE c.id = $1
	GROUP BY c.id
	LIMIT 1`

	insertEntry = `
	INSERT INTO public.content_entry (id, name, created_on, updated_on)
	VALUES ($1, $2, $3, $4)`

	insertEntryTrans = `
	INSERT INTO public.content_entry_field(content_entry_id, locale, fields) 
	VALUES($1, $2, $3)`

	updateEntry = `
	UPDATE public.content_entry SET (name, updated_on) = ($1, $2)
	WHERE id = $3`

	updateEntryTrans = `
	UPDATE public.content_entry_field SET fields = $1
	WHERE content_entry_id = $2 AND locale = $3`

	getEntryLocales = `
	SELECT locale FROM public.content_entry_field
	WHERE content_entry_id = $1
	`

	deleteentry = `
	DELETE FROM public.content_entry WHERE id = $1
	`
)

func (repo *PostgresRepository) GetAll() ([]*Entry, error) {
	rows, err := repo.db.Query(getAll)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	entrys := make([]*Entry, 0)

	for rows.Next() {
		entry := &Entry{}
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

func (repo *PostgresRepository) GetOne(id ID) (*Entry, error) {
	entry := CreateEntry()
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

	entry.ID, err = ParseId(i)

	if err != nil {
		return nil, errors.New("Could not parse ID")
	}

	return &entry, nil
}

func (repo *PostgresRepository) Add(entry Entry) error {
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

func (repo *PostgresRepository) getLocales(id ID) ([]string, error) {
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

func (repo *PostgresRepository) Update(entry Entry) error {
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

func (repo *PostgresRepository) Delete(id ID) error {
	_, err := repo.db.Exec(deleteentry, id)
	return err
}

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
