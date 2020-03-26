package contentitem

import (
	"database/sql"
	"errors"
)

type PostgresRepository struct {
	db *sql.DB
}

func (repo *PostgresRepository) GetAll() ([]*ContentItem, error) {
	rows, err := repo.db.Query(`
	SELECT * FROM public.content_item ORDER BY created_on ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contentItems := make([]*ContentItem, 0)
	for rows.Next() {
		contentItem := &ContentItem{}
		err := rows.Scan(&contentItem.ID, &contentItem.Name, &contentItem.Attrs, &contentItem.CreatedOn, &contentItem.UpdatedOn)
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
	SELECT id, name, attrs, created_on, updated_on FROM public.content_item WHERE content_item.id = $1
	`, id)
	var i string
	err := row.Scan(&i, &contentItem.Name, &contentItem.Attrs, &contentItem.CreatedOn, &contentItem.UpdatedOn)

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
	sqlStatement := `
	INSERT INTO public.content_item (id, name, attrs, created_on, updated_on)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	_, err := repo.db.Exec(
		sqlStatement,
		contentItem.ID,
		contentItem.Name,
		contentItem.Attrs,
		contentItem.CreatedOn,
		contentItem.UpdatedOn,
	)
	return err
}

func (repo *PostgresRepository) Update(contentItem ContentItem) error {
	sqlStatement := `
	UPDATE public.content_item SET (name, attrs, updated_on) = ($1, $2, $3)
	  WHERE id = $4`
	_, err := repo.db.Exec(sqlStatement, contentItem.Name, contentItem.Attrs, contentItem.UpdatedOn, contentItem.ID)
	return err
}

func (repo *PostgresRepository) Delete(id ID) error {
	sqlStatement := `DELETE FROM public.content_item WHERE id = $1`
	_, err := repo.db.Exec(sqlStatement, id)
	return err
}

func CreatePostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db,
	}
}
