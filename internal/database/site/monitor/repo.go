package monitor

import (
	"context"
	"database/sql"
	"log"

	"telebot/telebot/CA/internal/database/site/pgsql"

	"telebot/telebot/CA/internal/domain/entity"
)

type Repo struct {
	DB *sql.DB
}

//AlertSave - saving incomming URL into database's table "sites"
func (r *Repo) AlertSave(ctx context.Context, site *entity.Site) error {
	mdl := pgsql.NewSiteModel(*site)
	const sqlstr = `
				INSERT INTO
							sites(
											url,
											description,
  											create_at	
								)
				VALUES (
											$1,
											$2,
											$3
						)
				RETURNING id
					`
	err := r.DB.QueryRowContext(
		ctx,
		sqlstr,
		mdl.URL,
		mdl.Description,
		mdl.CreatedAt,
	).Scan(&site.URL)
	if err != nil {
		log.Printf("Что-то пошло не так site (r *Repo) Save %v", err)
		return err
	}
	return nil
}
