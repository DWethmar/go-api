package contentitem

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/DWethmar/go-api/pkg/database"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func (repo *PostgresRepository) GetAll() ([]*ContentItem, error) {
	rows, err := repo.db.Query(`
		SELECT 
			id, 
			name, 
			COALESCE(jsonb_object_agg(t.locale, t.attrs) FILTER (WHERE t.locale IS NOT NULL), '{}'::JSONB) as attrs,
			created_on, 
			updated_on
		FROM public.content_item c
		LEFT OUTER JOIN public.content_item_translation t ON c.id = t.content_item_id
		GROUP BY c.id
	`)
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
	row := repo.db.QueryRow(`
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
		`,
		id,
	)
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
	insertContentItem := `
		INSERT INTO public.content_item (id, name, created_on, updated_on)
		VALUES ($1, $2, $3, $4) RETURNING id
	`

	insertContentItemTrans := `
		INSERT INTO public.content_item_translation(content_item_id, locale, attrs) VALUES($1, $2, $3)
	`

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
			fmt.Printf("LIEP %v", err)
			if err != nil {
				return err
			}
		}

		return err
	})

	return err
}

func (repo *PostgresRepository) Update(contentItem ContentItem) error {
	sqlStatement := `
		UPDATE public.content_item SET (name, attrs, updated_on) = ($1, $2, $3)
		WHERE id = $4
	`
	_, err := repo.db.Exec(sqlStatement, contentItem.Name, contentItem.Attrs, contentItem.UpdatedOn, contentItem.ID)
	return err
}

func (repo *PostgresRepository) Delete(id ID) error {
	sqlStatement := `
		DELETE FROM public.content_item WHERE id = $1
	`
	_, err := repo.db.Exec(sqlStatement, id)
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
