package contentitem

import (
	"database/sql"
	"errors"
	"log"

	"github.com/DWethmar/go-api/pkg/database"
)

type PostgresRepository struct {
	db *sql.DB
}

var (
	getAll = `
	SELECT 
	id, 
	name, 
	COALESCE(jsonb_object_agg(t.locale, t.attrs) FILTER (WHERE t.locale IS NOT NULL), '{}'::JSONB) as attrs,
	created_on, 
	updated_on
	FROM public.content_item c
	LEFT OUTER JOIN public.content_item_translation t ON c.id = t.content_item_id
	GROUP BY c.id
	ORDER BY updated_on ASC`

	getOne = `
	SELECT 
	id, 
	name, 
	COALESCE(jsonb_object_agg(t.locale, t.attrs) FILTER (WHERE t.locale IS NOT NULL), '{}'::JSONB) as attrs,
	created_on, 
	updated_on
	FROM public.content_item c
	LEFT OUTER JOIN public.content_item_translation t ON c.id = t.content_item_id
	WHERE c.id = $1
	GROUP BY c.id
	LIMIT 1`

	insertContentItem = `
	INSERT INTO public.content_item (id, name, created_on, updated_on)
	VALUES ($1, $2, $3, $4)`

	insertContentItemTrans = `
	INSERT INTO public.content_item_translation(content_item_id, locale, attrs) 
	VALUES($1, $2, $3)`

	updateContentItem = `
	UPDATE public.content_item SET (name, updated_on) = ($1, $2)
	WHERE id = $3`

	updateContentItemTrans = `
	UPDATE public.content_item_translation SET attrs = $1
	WHERE content_item_id = $2 AND locale = $3`

	getLocales = `
	SELECT locale FROM public.content_item_translation
	WHERE content_item_id = $1
	`

	deleteContentItem = `
	DELETE FROM public.content_item WHERE id = $1
	`
)

func deleteTranslationStatement() {

}

func (repo *PostgresRepository) GetAll() ([]*ContentItem, error) {
	rows, err := repo.db.Query(getAll)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	contentItems := make([]*ContentItem, 0)

	for rows.Next() {
		contentItem := &ContentItem{}
		err := rows.Scan(
			&contentItem.ID,
			&contentItem.Name,
			&contentItem.Attrs,
			&contentItem.CreatedOn,
			&contentItem.UpdatedOn,
		)
		if err != nil {
			return nil, err
		}
		contentItems = append(contentItems, contentItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contentItems, nil
}

func (repo *PostgresRepository) GetOne(id ID) (*ContentItem, error) {
	contentItem := CreateContentItem()
	row := repo.db.QueryRow(getOne, id)

	var i string
	err := row.Scan(
		&i,
		&contentItem.Name,
		&contentItem.Attrs,
		&contentItem.CreatedOn,
		&contentItem.UpdatedOn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		panic(err)
	}

	contentItem.ID, err = ParseId(i)

	if err != nil {
		return nil, errors.New("Could not parse ID")
	}

	return &contentItem, nil
}

func (repo *PostgresRepository) Add(contentItem ContentItem) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			insertContentItem,
			contentItem.ID,
			contentItem.Name,
			contentItem.CreatedOn,
			contentItem.UpdatedOn,
		)

		if err != nil {
			return err
		}

		for locale, attrs := range contentItem.Attrs {
			_, err = tx.Exec(
				insertContentItemTrans,
				contentItem.ID,
				locale,
				attrs,
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
	rows, err := repo.db.Query(getLocales, id)

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

func (repo *PostgresRepository) Update(contentItem ContentItem) error {
	err := database.WithTransaction(repo.db, func(tx database.Transaction) error {
		_, err := repo.db.Exec(
			updateContentItem,
			contentItem.Name,
			contentItem.UpdatedOn,
			contentItem.ID,
		)

		if err != nil {
			return err
		}

		locales, err := repo.getLocales(contentItem.ID)
		if err != nil {
			return err
		}

		for locale, attrs := range contentItem.Attrs {
			hasLocale := false
			for _, n := range locales {
				if locale == n {
					hasLocale = true
				}
			}

			var err error
			if hasLocale {
				_, err = tx.Exec(
					updateContentItemTrans,
					attrs,
					contentItem.ID,
					locale,
				)

				if err != nil {
					return err
				}
			} else {
				_, err = tx.Exec(
					insertContentItemTrans,
					contentItem.ID,
					locale,
					attrs,
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
	_, err := repo.db.Exec(deleteContentItem, id)
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
