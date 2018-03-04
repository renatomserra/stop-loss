![Kraken](https://themerkle.com/wp-content/uploads/2016/05/Kraken-Exchange.jpg)

# Kraken Stop Loss bot
Kraken has disabled stop losses so i created this bot.
If the price for BTC to EUR is lower than your set stop loss, all your open orders are cancelled, otherwise kraken might reject the order as your funds can be locked to an open order and then a market sell is placed for the full amount of BTC held.
You will be notified at every step through telegram.
You can check the app is alive via telegram with /alive and an endpoint "/health" is setup for pinging (you can set alerting and if using a free heroku tier will keep your app alive)

## ENV vars to setup:

- KRAKEN_KEY `Your API key setup on kraken, remember to give necessary permissions`
- KRAKEN_SECRET `Your API key setup on kraken`
- STOP_LOSS_AT `A float value for your stop loss in EUR`
- TELEGRAM_TOKEN `A token you get when creating a telegram bot through BotFather`
- TELEGRAM_CHAT_ID `Start a chat with the created bot and use the /alive action to print your id and set it here`
- PORT `Default port the app should run in`

# Running locally

Just run `go run main.go` to run the app
App has dotenv setup so you can configure your env vars in a .env
