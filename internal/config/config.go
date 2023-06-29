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
	"errors"
	"fmt"
	h "hirohito/internal/helpers"
	m "hirohito/internal/models"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Public
	Configuration m.Config

	// private
	envBinds []string = []string{
		"loglevel",
		"discord_token",
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

	prefix := viper.GetString("prefix")
	if prefix == "" {
		Configuration.Global.Logger.Warn("no discord prefix provided, setting the default: !")
		Configuration.Discord.Prefix = "!"
	} else {
		Configuration.Discord.Prefix = prefix
	}

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

func init() {

	// Build config
	err := initEnv()
	if err != nil {
		log.Fatal("unable to init config. Bye.")
	}

	// Configure logger
	initLogging()

	// Configure Discord Client
	err = initDiscord()
	if err != nil {
		Configuration.Global.Logger.Fatalf("no discord client could be created: %v. The bot cannot function. Exiting..", err)
	}

}
