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

package hirohito

import (
	"context"
	c "hirohito/internal/config"
	h "hirohito/internal/helpers"

	"github.com/bwmarrin/discordgo"
)

var (
	minLength = 2
	maxLength = 100
	trueBool  = true
	falseBool = false

	// less typing by referencing
	discordClient = c.Configuration.Discord.Client
	logger        = c.Configuration.Global.Logger

	// All commands and options must have a description
	// Commands/options without description will fail the registration
	// of the command.
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "ping the bot",
		},
		{
			Name:        "source",
			Description: "Get a link to the bot's source code",
		},
		{
			Name:         "createjoinablechannel",
			Description:  "Create a joinable channel",
			DMPermission: &falseBool,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channelname",
					Description: "name of the channel to be created",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "topic",
					Description: "topic of the channel to be created",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
			},
		},
		{
			Name:         "deletejoinablechannel",
			Description:  "Delete a joinable channel",
			DMPermission: &falseBool,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channelname",
					Description: "name of the channel to be deleted",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
			},
		},
		{
			Name:         "setuphirohito",
			Description:  "Setup hirohito for your guild",
			DMPermission: &falseBool,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "joinchannelid",
					Description: "id of the channel people use to join joinable channels",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "adminchannelid",
					Description: "id of the channel admins use to create joinable channels",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "joinablechannelscategoryid",
					Description: "id of the category under which joinable channels must be created",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "anyoneroleid",
					Description: `id of the "@everyone role"`,
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "adminroleid",
					Description: "id of the administrator role in the guild",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "moderatorroleid",
					Description: "id of the moderators role in the guild",
					MinLength:   &minLength,
					MaxLength:   maxLength,
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping":                  ping,
		"source":                source,
		"createjoinablechannel": createJoinableChannel,
		"deletejoinablechannel": deleteJoinableChannel,
		"setuphirohito":         setupGuild,
	}
)

func Hirohito(ctx context.Context) {
	logger.Info("starting emperor hirohito")

	hirohitoCtx, hirohitoCancel := context.WithCancel(ctx)
	defer hirohitoCancel()

	err := c.DataStore.SetupDatastore(hirohitoCtx)
	if err != nil {
		logger.Fatalf("error setting up datastore. Bot cannot function. Error: %s", err)
	}

	// Register handler for incoming commands
	discordClient.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	discordClient.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	discordClient.AddHandler(reactionHandler)

	err = discordClient.Open()
	if err != nil {
		logger.Fatalf("Error opening discord session: %v", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discordClient.ApplicationCommandCreate(discordClient.State.User.ID, "", v)
		if err != nil {
			logger.Errorf("Cannot create command '%v'. Error: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer discordClient.Close()

	logger.Infoln("Bot is now running. Press CTRL-C to exit.")

	// wait for the context to report done and then do a cleanup
	<-ctx.Done()

	/*
		We need to fetch the commands, since deleting requires the command ID.
		We are doing this from the returned commands on line 70, because using
		this will delete all the commands, which might not be desirable, so we
		are deleting only the commands that we added.
	*/
	logger.Info("deregistering commands")

	registeredCommands, err = discordClient.ApplicationCommands(discordClient.State.User.ID, "")
	if err != nil {
		logger.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := discordClient.ApplicationCommandDelete(discordClient.State.User.ID, "", v.ID)
		if err != nil {
			logger.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.SendInteractionPingResponse(s, i)
}
