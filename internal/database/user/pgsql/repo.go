package pgsql

import (
	"context"
	"database/sql"
	"fmt"
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
                    user_name,
					chat_id,
					first_name,
					last_name,
					language_code,
					is_bot,
					created_at
               FROM users
               `

//SaveUser - creating new user in the database
func (r *Repo) SaveUser(ctx context.Context, user *entity.User) error {
	mdl := NewUserModel(*user)
	const sqlstr = `
				INSERT INTO
								users(
								      	user_id,
								      	user_name,
										chat_id,
										first_name,
										last_name,
								      	language_code,
								      	is_bot,
										created_at
										)
				VALUES (
										$1,
										$2,
										$3,
										$4,
				        				$5,
				        				$6,
				        				$7,
				        				$8
								)
				RETURNING id
								`
	err := r.DB.QueryRowContext(
		ctx,
		sqlstr,
		mdl.UserID,
		mdl.UserName,
		mdl.ChatID,
		mdl.FirstName,
		mdl.LastName,
		mdl.LanguageCode,
		mdl.IsBot,
		mdl.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		log.Printf("Check SaveUser's QueryRowContext for %v", err)
		return err
	}
	return nil
}

func (r *Repo) FindByUserID(ctx context.Context, userID uint) (bool, error) {
	var sqlstr = sqlstr + `
        WHERE user_id = $1
          `
	var u UserModel
	if err := r.DB.QueryRowContext(ctx,
		sqlstr,
		userID,
	).Scan(
		&u.ID,
		&u.UserID,
		&u.UserName,
		&u.ChatID,
		&u.FirstName,
		&u.LastName,
		&u.LanguageCode,
		&u.IsBot,
		&u.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not exist in the database %v", err)
			return false, fmt.Errorf("User not exist")
		}
		log.Printf("Check FindByIDAndUserID QueryRowContext for %v", err)
		return false, err
	}

	return true, nil
}

func (r *Repo) FindByIDAndUserID(ctx context.Context, id, userID uint) (*entity.User, error) {
	var sqlstr = sqlstr + `
        WHERE id = $1
          AND user_id = $2`

	var u UserModel
	if err := r.DB.QueryRowContext(ctx,
		sqlstr,
		id,
		userID,
	).Scan(
		&u.ID,
		&u.UserID,
		&u.UserName,
		&u.ChatID,
		&u.FirstName,
		&u.LastName,
		&u.LanguageCode,
		&u.IsBot,
		&u.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not exist in the database %v", err)
			return nil, nil
		}
		log.Printf("Check FindByIDAndUserID QueryRowContext for %v", err)
		return nil, err
	}

	e := ConvertIntoUserEntity(u)
	return &e, nil
}

func (r *Repo) Update(ctx context.Context, user *entity.User) error {
	m := NewUserModel(*user)
	const sqlstr = `
                UPDATE users
                SET    
                    user_name     = $1,
                    chat_id       = $2,
                    first_name    = $3,
                    last_name     = $4,
                    language_code = $5,
                    is_bot        = $6
                    
                
                WHERE  id         = $7
                   `
	_, err := r.DB.ExecContext(
		ctx,
		sqlstr,
		m.UserName,
		m.ChatID,
		m.FirstName,
		m.LastName,
		m.LanguageCode,
		m.IsBot,

		m.ID,
	)
	if err != nil {
		log.Printf("Check Update's ExecContext for %v", err)
		return err
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, user *entity.User) error {
	const sqlstr = `
                    DELETE
                    FROM users
                    WHERE id = $1`
	_, err := r.DB.ExecContext(
		ctx,
		sqlstr,
		user.ID,
	)
	if err != nil {
		log.Printf("Check ExecContext Delete for %v", err)
		return err
	}
	return nil
}
