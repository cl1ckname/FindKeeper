# Find Keeper Bot

Telegram bot for forwarding media from forwarded posts to a channel.

## Environment Variables

- `TELEGRAM_BOT_TOKEN` - Bot token from @BotFather
- `TELEGRAM_CHANNEL` - Channel ID (numeric value, e.g.: -1001234567890)

## Build and Run

### Local Build
```bash
make build
make run
```

### Docker
```bash
make docker-build
make docker-push
```

### Docker Compose (for Portainer)
1. Create `.env` file with your variables:
```
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHANNEL=-1001234567890
```

2. Run with docker-compose:
```bash
docker-compose up -d
```

### Install Dependencies
```bash
make install-deps
```

## Usage

1. Add the bot to your channel as an administrator
2. Get the channel ID (you can use @userinfobot)
3. Forward posts with media to the bot in private messages
4. The bot will automatically extract and send media to the channel

## Getting Channel ID

To get the channel ID:
1. Add @userinfobot to your channel
2. Forward any message from the channel to the bot
3. Copy the channel ID (starts with -100)
