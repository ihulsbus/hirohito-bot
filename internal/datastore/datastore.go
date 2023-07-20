/*
 Copyright (c) 2023 Ian Hulsbus

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package datastore

import (
	"context"
	"database/sql"
	"time"

	m "hirohito/internal/models"
)

// Constructor
type DataStoreClient interface{}

type DataStore struct {
	client *sql.DB
}

func DataStoreConstructor(client *sql.DB) *DataStore {
	return &DataStore{
		client: client,
	}
}

func (d DataStore) SetupDatastore(ctx context.Context) error {
	dsCtx, _ := context.WithTimeout(ctx, 60*time.Second)

	defer dsCtx.Done()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS "guildconfig" ("guildID" TEXT NOT NULL UNIQUE, "joinChannelID" TEXT, "adminChannelID" TEXT, "joinableChannelsCategoryID" TEXT, "anyoneRoleID" TEXT, "adminRoleID" TEXT, "moderatorRoleID" TEXT,  PRIMARY KEY("guildID"));
		CREATE TABLE IF NOT EXISTS "archiving" ("guildID"	TEXT NOT NULL UNIQUE, "auto" INTEGER NOT NULL DEFAULT 0, "interval"	INTEGER DEFAULT 60, PRIMARY KEY("guildID"));
	`
	tx, err := d.client.BeginTx(dsCtx, nil)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(sqlStmt); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Guild config
func (d DataStore) GetGuildInfo(guildID string) (*m.GuildInformation, error) {
	var data m.GuildInformation

	stmt, err := d.client.Prepare("SELECT guildID, joinChannelID, adminChannelID, joinableChannelsCategoryID, anyoneRoleID, adminRoleID, moderatorRoleID FROM guildconfig WHERE guildID = ?")
	if err != nil {
		return nil, err
	}

	if err := stmt.QueryRow(guildID).Scan(&data.GuildID, &data.JoinChannelID, &data.AdminChannelID, &data.JoinableChannelsCategoryID, &data.AnyoneRoleID, &data.AdminRoleID, &data.ModeratorRoleID); err != nil {
		return nil, err
	}

	return &data, nil
}

func (d DataStore) CreateGuildInfo(guildInfo m.GuildInformation) error {
	tx, err := d.client.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO guildconfig (guildID, joinChannelID, adminChannelID, joinableChannelsCategoryID, anyoneRoleID, adminroleid, moderatorroleid) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(guildInfo.GuildID, guildInfo.JoinChannelID, guildInfo.AdminChannelID, guildInfo.JoinableChannelsCategoryID, guildInfo.AnyoneRoleID, guildInfo.AdminRoleID, guildInfo.ModeratorRoleID); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (d DataStore) DeleteGuildInfo(guildID string) error {
	tx, err := d.client.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM guildconfig WHERE guildID = ?")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(guildID); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Archiving
func (d DataStore) GetArchivingInfo(guildID string) (*m.ArchivingInformation, error) {
	var data m.ArchivingInformation

	stmt, err := d.client.Prepare("Select guildID, auto, interval, archivingCategoryID FROM archiving WHERE guildID = ?")
	if err != nil {
		return nil, err
	}

	if err := stmt.QueryRow(guildID).Scan(&data.GuildID, &data.Auto, &data.Interval, &data.ArchivingCategoryID); err != nil {
		return nil, err
	}

	return &data, nil
}

func (d DataStore) CreateArchivingInfo(archivingInfo m.ArchivingInformation) error {
	tx, err := d.client.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO archiving (guildID, auto, interval, archivingCategoryID) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(archivingInfo.GuildID, archivingInfo.Auto, archivingInfo.Interval, archivingInfo.ArchivingCategoryID); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (d DataStore) DeleteArchivingInfo(guildID string) error {
	tx, err := d.client.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM archiving WHERE guildID = ?")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(guildID); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
