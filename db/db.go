package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	selectID    = `	SELECT id FROM zammad_vk WHERE vk = $1;`
	selectVK    = `	SELECT vk FROM zammad_vk WHERE id = $1;`
	deleteUser  = `	DELETE FROM zammad_vk WHERE vk = $1;`
	insertUser  = `	INSERT INTO zammad_vk (id, vk) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET vk = $2;`
	createTable = `	CREATE TABLE IF NOT EXISTS zammad_vk (
						id 	INTEGER PRIMARY KEY NOT NULL, 
						vk	INTEGER 			NOT NULL);`
)

type DB struct {
	DB  *sql.DB
	ctx context.Context
}

func NewDB(cfg config.DataBase, ctx context.Context) (s *DB, err error) {
	s = &DB{ctx: ctx}
	err = s.connect(cfg)
	return
}

func (s *DB) connect(cfg config.DataBase) (err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)

	if s.DB, err = sql.Open("postgres", dsn); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*2)
	defer cancel()

	if err = s.DB.PingContext(ctx); err != nil {
		return
	}

	if _, err = s.DB.ExecContext(ctx, createTable); err != nil {
		return
	}

	return
}

func (s *DB) InsertUser(vk, zammad int) (err error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err = s.DB.ExecContext(ctx, deleteUser, vk); err != nil {
		log.Error(err)
		return
	}

	if _, err = s.DB.ExecContext(ctx, insertUser, zammad, vk); err != nil {
		log.Error(err)
	}

	return
}

func (s *DB) DeleteUser(vk int) (err error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err = s.DB.ExecContext(ctx, deleteUser, vk); err != nil {
		log.Error(err)
	}

	return
}

func (s *DB) SelectZammad(vk int) (zammad int, err error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if err = s.DB.QueryRowContext(ctx, selectID, vk).Scan(&zammad); err != nil {
		log.Error(err)
	}

	return
}

func (s *DB) SelectVK(zammad int) (vk int, err error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if err = s.DB.QueryRowContext(ctx, selectVK, zammad).Scan(&vk); err != nil {
		log.Error(err)
	}

	return
}
