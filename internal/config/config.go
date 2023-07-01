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

package config

import (
	"database/sql"

	"errors"
	"fmt"
	c "hirohito/internal/channels"
	d "hirohito/internal/datastore"
	h "hirohito/internal/helpers"
	m "hirohito/internal/messages"
	"hirohito/internal/models"
	r "hirohito/internal/roles"
	u "hirohito/internal/users"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Public
	Configuration models.Config
	Channels      *c.Channels
	Roles         *r.Roles
	Users         *u.Users
	Messages      *m.Messages
	DataStore     *d.DataStore

	// private
	envBinds []string = []string{
		"loglevel",
		"discord_token",
		"datastore_path",
	}
)

func initEnv() error {
	viper.SetEnvPrefix("hirohito")

	for i := range envBinds {
		err := viper.BindEnv(envBinds[i])
		if err != nil {
			return fmt.Errorf("error binding to env var '%s': %s", envBinds[i], err.Error())
		}
	}

	return nil
}

func initLogging() {
	Configuration.Global.Logger = log.New()

	Configuration.Global.Logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logLevels := h.SetupLogLevels()

	if i, found := logLevels[strings.ToUpper(viper.GetString("loglevel"))]; found {
		Configuration.Global.Logger.SetLevel(i)

	} else {
		Configuration.Global.Logger.Warn("no or invalid loglevel specified. Assuming default Valid loglevels are: PANIC FATAL ERROR WARN INFO DEBUG TRACE")
		Configuration.Global.Logger.SetLevel(logLevels["INFO"])
	}

}

func initDiscord() error {
	var err error

	token := viper.GetString("discord_token")
	if token == "" {
		return errors.New("discord token not provided or not found")
	}

	log.Debug("creating Discord client")
	Configuration.Discord.Client, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	return err
}

func initDatastore() error {

	path := viper.GetString("datastore_path")
	if path == "" {
		return errors.New("datastore path not provided or not found")
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Unable to open the datastore: %v", err)
	}
	Configuration.DataStore.Client = db

	return nil
}

func init() {

	// Build config
	err := initEnv()
	if err != nil {
		log.Fatal("unable to init config. Bye.")
	}

	// Configure logger
	initLogging()

	// Configure datastore
	err = initDatastore()
	if err != nil {
		Configuration.Global.Logger.Fatalf("datastore could not be opened: %v. The bot cannot function. Exiting..", err)
	}

	// Configure Discord Client
	err = initDiscord()
	if err != nil {
		Configuration.Global.Logger.Fatalf("no discord client could be created: %v. The bot cannot function. Exiting..", err)
	}

	// init libraries
	DataStore = d.DataStoreConstructor(Configuration.DataStore.Client)
	Channels = c.ChannelsConstructor(Configuration.Discord.Client)
	Roles = r.RolesConstructor(Configuration.Discord.Client)
	Users = u.UsersConstructor(Configuration.Discord.Client)
	Messages = m.MessagesConstructor(Configuration.Discord.Client)

}
