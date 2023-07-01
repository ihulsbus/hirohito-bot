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

package users

import "github.com/bwmarrin/discordgo"

// Constructor
type DiscordClient interface{}

type Users struct {
	discordClient *discordgo.Session
}

func UsersConstructor(client *discordgo.Session) *Users {
	return &Users{
		discordClient: client,
	}
}

func (u Users) AssignUserToRole(guildID, userID, roleID string) error {

	err := u.discordClient.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

func (u Users) RemoveUserFromRole(guildID, userID, roleID string) error {

	err := u.discordClient.GuildMemberRoleRemove(guildID, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

func (u Users) GetUserRoles(guildID, userID string) ([]string, error) {

	user, err := u.discordClient.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}

	return user.Roles, nil
}
