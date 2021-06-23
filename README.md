# Mattermost SplitHooks Plugin
This plugin sends webhook notifications from Split.io to Mattermost.

## Installation

1. Go to the [releases page of this GitHub repository](https://github.com/stylianosrigas/mattermost-plugin-splithooks/releases) and download the latest release for your Mattermost server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin, and enable it. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

Next, to configure the plugin, follow these steps:

3. After you've uploaded the plugin in **System Console > Plugins > Management**, go to the plugin's settings page at **System Console > Plugins > SplitHooks**.
4. Specify the team and channel to send messages to.
5. Generate the Token that will be use to validate the requests.
6. Select which EnvirontNames to accept notifications for (comma separated string).
7. Hit **Save**.
8. Next, copy the **Token** above the **Save** button, which is used to configure the plugin for your Split.io account.
9. Go to your Split.io account, paste the following webhook URL and specify the token you copied in step 8.

```
https://SITEURL/plugins/com.mattermost.plugin-splithooks/webhook&token=TOKEN
```
