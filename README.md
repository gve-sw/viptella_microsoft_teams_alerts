# viptella_microsoft_teams_alerts
This application acts as a web server to receive webhooks from Viptella and creates alerts in a Microsoft Teams channel via the webhook connector functionality.


## Contacts
* Charles Llewellyn

## Solution Components
* Go
*  Viptella
*  Microsoft Teams

## Related Sandbox Environment

Cisco Secure SD-WAN (Viptela) v20.9.1 / 17.9.1 - Single DCv1





## Installation/Configuration

1. Create a [Microsoft Teams Channel Webhook](https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook)
2. Edit the config.yaml file with your Microsoft Teams webhook url and the port that you want to run the Go http server on.
3. On Viptella -- Create a webhook for the alarms that you want Microsoft Teams to receive. (Select "Email Notifications" from "Monitor -> Alarms".)
4. Run Go Server (go run server.go)


## Usage

To launch the go webserver:


    $ go run server.go



# Screenshots

![/IMAGES/screenshot.png](/IMAGES/screenshot.png)

### LICENSE

Provided under Cisco Sample Code License, for details see [LICENSE](LICENSE.md)

### CODE_OF_CONDUCT

Our code of conduct is available [here](CODE_OF_CONDUCT.md)

### CONTRIBUTING

See our contributing guidelines [here](CONTRIBUTING.md)

#### DISCLAIMER:
<b>Please note:</b> This script is meant for demo purposes only. All tools/ scripts in this repo are released for use "AS IS" without any warranties of any kind, including, but not limited to their installation, use, or performance. Any use of these scripts and tools is at your own risk. There is no guarantee that they have been through thorough testing in a comparable environment and we are not responsible for any damage or data loss incurred with their use.
You are responsible for reviewing and testing any scripts you run thoroughly before use in any non-testing environment.
