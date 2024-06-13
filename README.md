# Anna University Job Alerts

This project is a web scraper built in Go that fetches job listings from Anna University's events page. It checks for new job entries and sends notifications via Telegram if new jobs are found. The job details are stored in an SQLite database to prevent duplicate notifications. The project includes logging for error tracking and process verification.

## Features

- Scrapes job listings from Anna University's events page
- Checks for new job entries and sends notifications via Telegram
- Stores job details in an SQLite database
- Includes logging for error tracking and process verification

## Prerequisites

- Go 1.16 or later
- SQLite3
- Telegram Bot Token

## Installation

1. **Clone the repository**:

    ```sh
    git clone https://github.com/GoliathReaper/Anna_Job_Alert__Telegram_GO
    cd Anna_Job_Alert__Telegram_GO
    ```

2. **Install dependencies**:

    ```sh
    go get -u github.com/PuerkitoBio/goquery
    go get -u github.com/mattn/go-sqlite3
    go get -u gopkg.in/tucnak/telebot.v2
    ```

3. **Set up your environment**:
    - Create a Telegram bot and get the bot token.
    - Set the bot token and chat ID in the `botToken` and `chatIDStr` constants in `main.go`.

4. **Build the project**:

    For Windows 64-bit:
    ```sh
    set GOOS=windows
    set GOARCH=amd64
    set CGO_ENABLED=1
    go build -o annajobs.exe
    ```

    For Linux ARM (Raspberry Pi Zero):
    ```sh
    export GOOS=linux
    export GOARCH=arm
    export GOARM=6
    export CGO_ENABLED=1
    go build -o annajobs_pi
    ```

5. **Run the executable**:
    ```sh
    ./annajobs.exe  # For Windows
    ./annajobs_pi  # For Raspberry Pi
    ```

## Usage

- Ensure that the SQLite database (`AnnaUnivJobAlert.db`) is in the same directory as the executable.
- Run the executable to start scraping job listings and receiving Telegram notifications.

## Code Explanation

The main parts of the code include:

- **Scraping the job listings** from the specified URL using `goquery`.
- **Checking for new job entries** in the SQLite database to avoid duplicate notifications.
- **Sending notifications** via Telegram using the `tucnak/telebot.v2` package.
- **Logging** errors and actions to a log file for debugging and verification.

## Error Handling

- Any errors encountered during scraping, database operations, or sending notifications are logged to a `job_alerts.log` file.
- Telegram notifications are sent in case of critical errors.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
