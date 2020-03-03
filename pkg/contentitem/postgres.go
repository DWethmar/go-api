package contentitem

import "database/sql"

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
		contentItem := new(ContentItem)
		err := rows.Scan(&contentItem.ID, &contentItem.Name, &contentItem.Data, &contentItem.CreatedOn, &contentItem.UpdatedOn)
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

func (repo *PostgresRepository) GetOne(id int64) (*ContentItem, error) {
	var contentItem ContentItem
	row := repo.db.QueryRow(`
	SELECT * FROM public.content_item WHERE content_item.id = $1
	`, id)
	err := row.Scan(&contentItem.ID, &contentItem.Name, &contentItem.Data, &contentItem.CreatedOn, &contentItem.UpdatedOn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		panic(err)
	}
	return &contentItem, nil
}

func (repo *PostgresRepository) Create(contentItem *ContentItem) error {
	sqlStatement := `
	INSERT INTO public.content_item (name, data, created_on, updated_on)
	VALUES ($1, $2, $3, $4)`
	_, err := repo.db.Exec(sqlStatement, contentItem.Name, contentItem.Data, contentItem.CreatedOn, contentItem.UpdatedOn)
	return err
}

func (repo *PostgresRepository) Update(contentItem *ContentItem) error {
	sqlStatement := `
	UPDATE public.content_item SET (name, data, updated_on) = ($1, $2, $3)
	  WHERE id = $4`
	_, err := repo.db.Exec(sqlStatement, contentItem.Name, contentItem.Data, contentItem.UpdatedOn, contentItem.ID)
	return err
}

func (repo *PostgresRepository) Delete(id int64) error {
	sqlStatement := `DELETE FROM public.content_item WHERE id = $1`
	_, err := repo.db.Exec(sqlStatement, id)
	return err
}

func CreatePostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db,
	}
}
