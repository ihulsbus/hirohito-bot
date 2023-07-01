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

package helpers

import (
	"errors"
	"fmt"
	m "hirohito/internal/models"
	"regexp"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	InsufficientPermissions string = "Insufficient permissions for the command"
	UnknownOption           string = "unknown option provided to command. Aborting."
)

func SetupLogLevels() map[string]log.Level {
	logLevels := make(map[string]log.Level, 7)

	logLevels["PANIC"] = log.PanicLevel
	logLevels["FATAL"] = log.FatalLevel
	logLevels["ERROR"] = log.ErrorLevel
	logLevels["WARN"] = log.WarnLevel
	logLevels["INFO"] = log.InfoLevel
	logLevels["DEBUG"] = log.DebugLevel
	logLevels["TRACE"] = log.TraceLevel

	return logLevels
}

func FindChannelRole(roles []*discordgo.Role, channelName string) (int, bool) {
	for i, role := range roles {
		if role.Name == channelName {
			return i, true
		}
	}
	return -1, false
}

func FindRoleID(roleIDs []string, roleID string) (int, bool) {
	for i := range roleIDs {
		if roleIDs[i] == roleID {
			return i, true
		}
	}
	return -1, false
}

func FindChannel(channels []*discordgo.Channel, channelName string) (int, bool) {
	for i, channel := range channels {
		if channel.Name == channelName {
			return i, true
		}
	}
	return -1, false
}

func sendInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, resp *discordgo.InteractionResponse) error {
	err := s.InteractionRespond(i.Interaction, resp)
	if err != nil {
		return err
	}
	return nil
}

func SendInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}

	return sendInteraction(s, i, &resp)
}

func SendInteractionAwaitResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}

	return sendInteraction(s, i, &resp)
}

func SendInteractionAwaitUpdate(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}

	return sendInteraction(s, i, &resp)
}

func SendInteractionPingResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	}

	sendInteraction(s, i, &resp)
}

func PermissionChecker(guildInfo *m.GuildInformation, i *discordgo.InteractionCreate) bool {
	_, admin := FindRoleID(i.Member.Roles, guildInfo.AdminRoleID)
	_, mod := FindRoleID(i.Member.Roles, guildInfo.ModeratorRoleID)

	if admin || mod {
		return i.ChannelID == guildInfo.AdminChannelID
	}

	return false
}

func FindChannelInGuild(d *discordgo.Session, guildID, channelName string) (*discordgo.Channel, error) {
	guildChannels, err := d.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve guild channels: %s", err)
	}

	pos, found := FindChannel(guildChannels, channelName)
	if !found {
		return nil, errors.New("given channel name not found in guild's channel list")
	}

	return guildChannels[pos], nil
}

func FindChannelEmbedMessage(messages []*discordgo.Message, channelName string) (*discordgo.Message, error) {

	exp, err := regexp.Compile(`"(\w.*)"$`)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		if len(messages[i].Embeds) > 0 {
			for _, Embed := range messages[i].Embeds {
				result := exp.FindStringSubmatch(Embed.Title)
				if len(result) == 2 && result[1] == channelName {
					return messages[i], nil
				}
			}
		}
	}

	return nil, errors.New("no embed found for channel")
}
