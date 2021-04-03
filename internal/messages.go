package internal

const (
	MessageNoAdmin string = `No admin has been set for this bot! This bot will not function until an admin is set!

If you are the administrator of this bot, launch the bot with the following environment variable set:

` + "`AB_ADMIN_ID=%d`" + `

Once you have done so, run /start again for further setup information.`

	MessageInitialStart string = `Congratulations! You are registered as the admin for this bot! This message will eventually tell you how to setup command completion and give you a short overview of how the bot works. The bot is currently under development, so watch this space.`
)
