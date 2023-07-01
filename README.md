# hirohito-bot
---

This is the main bot for the DAGC Discord server. 

This bot currently fulfills the following purposes:

* Allow users to join/leave optional channels
* Allow admins to automatically archive inactive channels
* Allow admins to create joinable channels

Additional functionality will be added when I need it.

## Building and packaging
TBD

## Running the bot

The bot uses env variables for it's configuration. The table blow shows the supported environment variables.

|varname                        |type     |required |
|---                            |---      |---      |
|HIROHITO_LOGLEVEL              |`string` |no       |
|HIROHITO_DATASTORE_PATH        |`string` |no       |
|HIROHITO_DISCORD_TOKEN         |`string` |yes      |

You'll need to create a discord bot via discord's developer portal to generate a token. 
After you've generated the token, you'll need to ensure that the bot has the "MESSAGE CONTENT INTENT" active.

Then invite the bot with the following permissions: 

* Administrator

This should give the following bot permission integer: 
`8`

Unfortunately the bot does not currentl support more narrowly scoped permissions (I tried).

## Datastore

The bot uses a sqlite DB to persistently store discord snowflake configuration information per guild. This includes the following information: 
* guildID
* ID of the channel where join embeds are posted
* ID of the channel where the bot's admin commands are executed
* ID of the category under which the joinable channels should be created
* ID of the "@anyone" role that needs to have their permissions manipulated on the newly created channels
* ID of the admin role that always has access to the joinable channel
* ID of the mods role that always has access to the joinable channel
