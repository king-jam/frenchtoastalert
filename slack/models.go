package slack

import (
	"fmt"
	"github.com/king-jam/ft-alert-bot/models"
)

func SlackAlertFromTemplate(slackAlert models.SlackAlert) string {
	slackTemplate :=
		`
	{
		"blocks": [
			{
				"type": "divider"
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "*French Toast Alert:*"
				},
				"accessory": {
					"type": "image",
					"image_url": "https://www.seriouseats.com/images/2015/04/20140411-french-toast-recipe-09-edit.jpg",
					"alt_text": "French Toast Alerts"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "%s, %s"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "Toast Level: %d"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "\nExpected Snowfall:\t\t%f\n"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "\nLow End Snowfall:\t\t%f\n"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "\nHigh End Snowfall:\t\t%f\n"
				}
			},
			{
				"type": "divider"
			}
		]
	}
	`

	return fmt.Sprintf(slackTemplate, slackAlert.City, slackAlert.State, slackAlert.ToastLevel, slackAlert.ExpectedSnowfall, slackAlert.LowSnowfall, slackAlert.HighSnowfall)

}
