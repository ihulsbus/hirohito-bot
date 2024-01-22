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
	"fmt"
	c "hirohito/internal/config"
	h "hirohito/internal/helpers"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func createJoinableChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var permission []*discordgo.PermissionOverwrite
	var name, topic string

	guildInfo, err := checkGuildSetup(i.GuildID)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	permitted := h.PermissionChecker(guildInfo, i)
	if !permitted {
		h.SendInteractionResponse(s, i, h.InsufficientPermissions)
		return
	}

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "channelname":
			name = strings.ToLower(option.StringValue())
			name = strings.ReplaceAll(name, " ", "-")
		case "topic":
			topic = option.StringValue()
		default:
			h.SendInteractionResponse(s, i, h.UnknownOption)
			return
		}
	}

	if name == "" || topic == "" {
		h.SendInteractionResponse(s, i, "name or topic are empty. Both need to be between 2 and 100 characters.")
		return
	}

	guildChannels, err := s.GuildChannels(i.GuildID)
	if err != nil {
		h.SendInteractionResponse(s, i, fmt.Sprintf("Unable to retrieve list of guild channels to check uniqueness: %s", err))
		return
	}

	if _, found := h.FindChannel(guildChannels, name); found {
		h.SendInteractionResponse(s, i, "Requested channel name already exists.")
		return
	}

	roleData := discordgo.RoleParams{
		Name:        name,
		Hoist:       &falseBool,
		Mentionable: &falseBool,
	}

	role, err := c.Roles.CreateRole(i.GuildID, &roleData)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	permission = append(permission,
		&discordgo.PermissionOverwrite{
			ID:    role.ID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: 197632,
		},
		&discordgo.PermissionOverwrite{
			ID:   guildInfo.AnyoneRoleID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: 1024,
		},
		&discordgo.PermissionOverwrite{
			ID:    guildInfo.AdminRoleID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: 66560,
		},
		&discordgo.PermissionOverwrite{
			ID:    guildInfo.ModeratorRoleID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: 66560,
		},
	)

	channelData := discordgo.GuildChannelCreateData{
		Name:                 name,
		Topic:                topic,
		ParentID:             guildInfo.JoinableChannelsCategoryID,
		PermissionOverwrites: permission,
	}

	channel, err := c.Channels.CreateTextChannel(i.GuildID, channelData)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	err = c.Messages.JoinableChannelEmbed(i.GuildID, guildInfo.JoinChannelID, channel)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	h.SendInteractionResponse(s, i, fmt.Sprintf("Channel created: %v", channel.Mention()))
}

func deleteJoinableChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var name string

	guildInfo, err := checkGuildSetup(i.GuildID)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	permitted := h.PermissionChecker(guildInfo, i)
	if !permitted {
		h.SendInteractionResponse(s, i, h.InsufficientPermissions)
		return
	}

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "channelname":
			name = option.StringValue()
		default:
			h.SendInteractionResponse(s, i, h.UnknownOption)
			return
		}
	}

	if name == "" {
		h.SendInteractionResponse(s, i, "name is empty. name needs to be between 2 and 100 characters.")
		return
	}

	messages, err := c.Messages.GetMessagesInChannel(guildInfo.JoinChannelID)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	message, err := h.FindChannelEmbedMessage(messages, name)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	err = c.Messages.DeleteMessage(guildInfo.JoinChannelID, message.ID)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
	}

	roles, err := c.Roles.RetrieveRoles(guildInfo.GuildID)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
		return
	}

	pos, found := h.FindChannelRole(roles, name)
	if found {
		err = c.Roles.DeleteRole(guildInfo.GuildID, roles[pos].ID)
		if err != nil {
			h.SendInteractionResponse(s, i, err.Error())
		}
	}

	guildChannel, err := h.FindChannelInGuild(s, i.GuildID, name)
	if err != nil {
		h.SendInteractionResponse(s, i, err.Error())
	}

	err = c.Channels.DeleteTextChannel(guildChannel.ID)
	if err != nil {
		h.SendInteractionResponse(s, i, fmt.Sprintf("unable to delete channel: %s", err))
		return
	}

	err = h.SendInteractionResponse(s, i, "Channel deleted")
	if err != nil {
		logger.Errorf("unable to send response to guild: %s", err)
	}

}
