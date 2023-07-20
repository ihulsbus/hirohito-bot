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
	"errors"
	c "hirohito/internal/config"
	h "hirohito/internal/helpers"
	m "hirohito/internal/models"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func checkGuildSetup(guildID string) (*m.GuildInformation, error) {
	var guildInfo *m.GuildInformation
	var err error
	var noSetupMsg string = `hirohito has no setup for this guild. use the "setuphirohito" command to get started`

	guildInfo, err = c.DataStore.GetGuildInfo(guildID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(noSetupMsg)
		}
		return nil, err
	}

	if guildInfo.GuildID == guildID && guildInfo.JoinChannelID != "" && guildInfo.JoinableChannelsCategoryID != "" {
		return guildInfo, nil
	}

	return nil, errors.New(noSetupMsg)
}

func setupGuild(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var guildInfo m.GuildInformation

	existingGuildInfo, err := checkGuildSetup(i.GuildID)
	if err != nil {
		if !strings.Contains(err.Error(), "hirohito has no setup for this guild") {
			h.SendInteractionResponse(s, i, err.Error())
			return
		}
	}

	if existingGuildInfo != nil && existingGuildInfo.GuildID != "" {
		permitted := h.PermissionChecker(existingGuildInfo, i)
		if !permitted {
			h.SendInteractionResponse(s, i, h.InsufficientPermissions)
			return
		}
	}

	guildInfo.GuildID = i.GuildID

	if len(i.ApplicationCommandData().Options) < 1 {
		h.SendInteractionResponse(s, i, "no values in command args present. Halting operation!")
		return
	}

	// the most nasty shit ever. Name based extraction? hellooo???
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "joinchannelid":
			guildInfo.JoinChannelID = option.StringValue()

		case "adminchannelid":
			guildInfo.AdminChannelID = option.StringValue()

		case "joinablechannelscategoryid":
			guildInfo.JoinableChannelsCategoryID = option.StringValue()

		case "anyoneroleid":
			guildInfo.AnyoneRoleID = option.StringValue()

		case "adminroleid":
			guildInfo.AdminRoleID = option.StringValue()

		case "moderatorroleid":
			guildInfo.ModeratorRoleID = option.StringValue()

		default:
			h.SendInteractionResponse(s, i, "unrecognised option! Halting operation!")
			return
		}
	}

	err = c.DataStore.CreateGuildInfo(guildInfo)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	h.SendInteractionResponse(s, i, "Guild setup completed")

}
