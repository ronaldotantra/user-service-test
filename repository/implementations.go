package repository

import (
	"context"
	"database/sql"
)

func (r *repository) GetUserById(ctx context.Context, id int64) (*User, error) {
	output := &User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, name, phone, password FROM users WHERE id = $1", id).
		Scan(&output.Id, &output.Name, &output.Phone, &output.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return output, nil
}

func (r *repository) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	output := &User{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, name, phone, password FROM users WHERE phone = $1", phone).
		Scan(&output.Id, &output.Name, &output.Phone, &output.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return output, nil
}

func (r *repository) UpdateProfile(ctx context.Context, user User) error {
	query := `
	UPDATE users
	SET 
	name = $2,
	phone = $3,
	updated_at = NOW()
	WHERE id = $1;`

	_, err := r.Db.ExecContext(ctx, query,
		user.Id,
		user.Name,
		user.Phone,
	)
	return err
}

func (r *repository) InsertUser(ctx context.Context, user User) (*int64, error) {
	var id int64
	query := `
	INSERT INTO users(id, name, password, phone, created_at, updated_at) VALUES
	(DEFAULT, $1,$2,$3, NOW(), NOW()) RETURNING id;`

	err := r.Db.QueryRowContext(ctx, query,
		user.Name,
		user.Password,
		user.Phone,
	).Scan(&id)

	return &id, err
}

func (r *repository) GetUserToken(ctx context.Context, userId int64) (*UserToken, error) {
	output := &UserToken{}
	err := r.Db.QueryRowContext(ctx, "SELECT id, user_id, token, count_login FROM user_tokens WHERE user_id = $1", userId).
		Scan(&output.Id, &output.UserId, &output.Token, &output.CountLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return output, nil
}

func (r *repository) InsertToken(ctx context.Context, payload TokenPayloadInsert) error {
	query := `
	INSERT INTO user_tokens(id, user_id, token, count_login, created_at, updated_at) VALUES
	(DEFAULT, $1,$2,1, NOW(), NOW())`

	_, err := r.Db.ExecContext(ctx, query,
		payload.UserId,
		payload.Token,
	)
	return err

}

func (r *repository) UpdateToken(ctx context.Context, payload TokenPayloadUpdate) error {
	query := `
	UPDATE user_tokens
	SET 
	token = $2,
	count_login = count_login + 1,
	updated_at = NOW()
	WHERE id = $1;`

	_, err := r.Db.ExecContext(ctx, query,
		payload.Id,
		payload.Token,
	)
	return err
}
