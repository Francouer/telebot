package pgsql

import (
	"context"
	"database/sql"
	"log"

	"telebot/telebot/CA/internal/domain/entity"
)

type Repo struct {
	DB *sql.DB
}

const sqlstr = `
				SELECT 
					id,
					user_id,
					name,
					url,
					request_timeout,
					response_status,
					description,
					created_at
				FROM sites
               `

//Save -
func (r *Repo) Save(ctx context.Context, site *entity.Site) error {
	mdl := NewSiteModel(*site)
	const sqlstr = `
                INSERT INTO
                           sites(	
                                 user_id,
                                 name,
                                 url,
                                 request_timeout,
                                 response_status,
                                 description,
                                 created_at	
                                 )
                VALUES           (
                                 $1,
                                 $2,
                                 $3,
                                 $4,
                                 $5,
                                 $6,
                                 $7
                                 )
                RETURNING        id
					`
	err := r.DB.QueryRowContext(
		ctx,
		sqlstr,
		mdl.UserID,
		mdl.Name,
		mdl.URL,
		mdl.RequestTimeout, //TODO: service RequestTimeout = default = "120"
		mdl.ResponseStatus, //TODO: service ResponseStatus = default = "200"
		mdl.Description,    //TODO: service Description = default = "pending"
		mdl.CreatedAt,
	).Scan(&site.ID)
	if err != nil {
		log.Printf("Found an error in SaveSiteURLByName %v", err)
		return err
	}
	return nil
}

//FindByUserID -
func (r *Repo) FindByUserID(ctx context.Context, userID uint) ([]entity.Site, error) {
	var sqlstr = sqlstr + ` 
                          WHERE user_id = $1`
	var m SiteModel
	var models []SiteModel
	log.Printf("FindByUserID  %v", userID)
	rows, err := r.DB.QueryContext(
		ctx,
		sqlstr,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Check FindByIDAndUserID's ErrNoRows for %v", err)
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Name,
			&m.URL,
			&m.RequestTimeout,
			&m.ResponseStatus,
			&m.Description,
			&m.CreatedAt,
		)
		models = append(models, m)
	}
	e := ConvertSliceOfSiteModels(models)
	return e, nil
}

//FindByIDAndUserID -
func (r *Repo) FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.Site, error) {
	var sqlstr = sqlstr + `
		WHERE id = $1
		  AND user_id = $2`

	var m SiteModel
	log.Printf("FindByIDAndUserID %v, %v", id, userID)
	if err := r.DB.QueryRowContext(
		ctx,
		sqlstr,
		id,
		userID,
	).Scan(
		&m.ID,
		&m.UserID,
		&m.Name,
		&m.URL,
		&m.RequestTimeout,
		&m.ResponseStatus,
		&m.Description,
		&m.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Check FindByIDAndUserID's ErrNoRows for %v", err)
			return nil, nil
		}
		return nil, err
	}
	e := ConvertIntoSiteEntity(m)
	log.Printf("Site, err - %v ", &e)
	return &e, nil
}

//Update -
func (r *Repo) Update(ctx context.Context, site *entity.Site) error {
	m := NewSiteModel(*site)
	const sqlstr = `
                UPDATE sites
                SET    
                    name            = $1,
                    url             = $2,
                    request_timeout = $3,
                    response_status = $4,
                    description     = $5
                
                WHERE  id           = $6
                   `

	_, err := r.DB.ExecContext(
		ctx,
		sqlstr,
		m.Name,
		m.URL,
		m.RequestTimeout,
		m.ResponseStatus,
		m.Description,

		m.ID,
	)
	if err != nil {
		log.Printf("Check Update's ExecContext for %v", err)
		return err
	}
	return nil
}

//Delete -
func (r *Repo) Delete(ctx context.Context, id uint) error {
	const sqlstr = `
                    DELETE 
                    FROM  sites
                    WHERE id = $1
				`
	_, err := r.DB.ExecContext(
		ctx,
		sqlstr,

		id,
	)
	if err != nil {
		log.Printf("Check ExecContext Delete for %v", err)
		return err
	}
	return nil
}

//FindByUserIdURL - searching database's table "sites"
// for match incomming  URLs with database's URLs.
// True, nil if got a match.
func (r *Repo) FindByUserIdAndURL(ctx context.Context, UserID uint, URL string) (*entity.Site, error) {
	var sqlstr = sqlstr + `
		WHERE user_id = $1
		  AND url = $2`

	var m SiteModel
	if err := r.DB.QueryRowContext(ctx,
		sqlstr,
		UserID,
		URL,
	).Scan(
		&m.ID,
		&m.UserID,
		&m.Name,
		&m.URL,
		&m.RequestTimeout,
		&m.ResponseStatus,
		&m.Description,
		&m.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Print("ErrNoRows in FindByUserIdAndURL")
			return nil, nil
		}
		return nil, err
	}
	e := ConvertIntoSiteEntity(m)
	return &e, nil
}

//FindByUserIdName - searching database's table "sites"
// for match incomming  Names with database's Names.
// True, nil if got a match.
func (r *Repo) FindByUserIdAndName(ctx context.Context, UserID uint, Name string) (*entity.Site, error) {
	var sqlstr = sqlstr + `
		WHERE user_id = $1
		  AND name    = $2`
	var m SiteModel
	if err := r.DB.QueryRowContext(
		ctx,
		sqlstr,
		UserID,
		Name,
	).Scan(
		&m.ID,
		&m.UserID,
		&m.Name,
		&m.URL,
		&m.RequestTimeout,
		&m.ResponseStatus,
		&m.Description,
		&m.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Print("ErrNoRows in FindByUserIdAndURL")
			return nil, nil
		}
		return nil, err
	}
	e := ConvertIntoSiteEntity(m)
	return &e, nil
}
