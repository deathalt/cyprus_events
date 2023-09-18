# Event Scraper and Notifier

A simple Go application that scrapes event details from "www.soldoutticketbox.com" every hour and inserts unique events into an SQLite database. If a new unique event is inserted into the database, a notification is sent via a Telegram bot.

## Table of Contents

- [Setup](#setup)
- [Usage](#usage)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Setup

### Prerequisites

1. **Go**: Ensure you have Go installed on your machine. If not, follow the official [Go installation guide](https://golang.org/doc/install).
2. **SQLite**: This application uses SQLite to store event details. Ensure SQLite is installed on your machine.
3. **Telegram Bot**: You'll need a Telegram bot token to send notifications and the chat ID of the destination.

### Installation

1. Clone the repository:

   ```bash
   git clone https://path-to-your-repository.git
   cd path-to-repository-folder
   ```

2. Install required Go packages:

   ```bash
   go get github.com/PuerkitoBio/goquery
   go get github.com/mattn/go-sqlite3
   ```

3. Set environment variables:

   - `TG_BOT_TOKEN`: Your Telegram Bot Token.
   - `TG_CHAT_ID`: The Chat ID where notifications will be sent.

   ```bash
   export TG_BOT_TOKEN='your-telegram-bot-token'
   export TG_CHAT_ID='your-chat-id'
   ```

## Usage

Simply run the Go application:

```bash
go run main.go
```

The application will scrape event details every hour and notify you on Telegram for any new unique events.

## Architecture

- **SQLite**: An embedded database used to store and deduplicate event details.
- **goquery**: A package that enables HTTP requests and HTML parsing in a jQuery-like syntax.
- **Telegram Bot API**: Used to send notifications for newly added events.

## Contributing

We welcome any contributions! If you find any issues or would like to add enhancements, please open an issue or a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.