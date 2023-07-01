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
	c "hirohito/internal/config"
	h "hirohito/internal/helpers"

	"regexp"

	"github.com/bwmarrin/discordgo"
)

func reactionHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	guildInfo, err := checkGuildSetup(m.GuildID)
	if err != nil {
		return
	}

	if m.ChannelID != guildInfo.JoinChannelID || m.UserID == s.State.User.ID {
		return
	}

	err = s.MessageReactionRemove(m.ChannelID, m.MessageReaction.MessageID, m.Emoji.APIName(), m.UserID)
	if err != nil {
		logger.Errorf("error removing reaction from user %s: %s", m.UserID, err)
		return
	}

	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		logger.Errorf("error retrieving message from reaction: %s", err)
		return
	}

	exp, err := regexp.Compile(`^<#([0-9].*)>`)
	if err != nil {
		logger.Error("Unable to compile regex")
	}

	result := exp.FindStringSubmatch(message.Embeds[0].Fields[0].Value)
	if len(result) < 1 {
		logger.Error("unable to get channelid from embed")
		return
	}

	channel, err := s.Channel(result[1])
	if err != nil {
		logger.Errorf("Unable to get channel: %s", err)
	}

	roleList, err := c.Roles.RetrieveRoles(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	i, found := h.FindChannelRole(roleList, channel.Name)
	if !found {
		logger.Error("Unable to find role of mentioned channel")
		return
	}

	userRoles, err := c.Users.GetUserRoles(m.GuildID, m.UserID)
	if err != nil {
		logger.Errorf("Error getting user %s roles: %s", m.UserID, err)
	}

	switch m.Emoji.APIName() {
	case "▶️":

		if _, found := h.FindRoleID(userRoles, roleList[i].ID); found {
			return
		}

		err = c.Users.AssignUserToRole(m.GuildID, m.UserID, roleList[i].ID)
		if err != nil {
			logger.Errorf("error assignng role %s to user %s. Error: %s", roleList[i].Name, m.UserID, err)
			return
		}

		c.Messages.UserJoinedChannelMessage(m.GuildID, channel.ID, *m.Member.User)
	case "🚮":

		if _, found := h.FindRoleID(userRoles, roleList[i].ID); !found {
			return
		}

		err = c.Users.RemoveUserFromRole(m.GuildID, m.UserID, roleList[i].ID)
		if err != nil {
			logger.Errorf("error removing role %s from user %s. Error: %s", roleList[i].Name, m.UserID, err)
			return
		}

		c.Messages.UserLeftChannelMessage(m.GuildID, channel.ID, *m.Member.User)
	}

}
