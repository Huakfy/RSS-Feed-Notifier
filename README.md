# RSS-Feed-Notifier

A lightweight Go-based security automation tool that monitors RSS feeds for vulnerability intelligence (CVE feeds) and sends real-time alerts to a chat channel (Discord).<br>
Designed to run as a scheduled job (cron / Task Scheduler) for continuous DevSecOps monitoring.

## 🚀 Features
- 🔎 Fetches vulnerability RSS feeds (CVE / security feeds)
- 🧠 Parses RSS using a robust Go parser
- 🔁 Deduplicates alerts using persistent GUID tracking
- 🧹 Converts HTML content into Markdown-safe formatting
- 📡 Sends formatted alerts to Discord via bot API
- 💾 Maintains lightweight local state per feed (file-based)
- ⏱ Designed for periodic execution (cron-friendly)

## 🏗 Architecture Overview
```
RSS Feed
   ↓
Parser (gofeed)
   ↓
HTML → Markdown conversion
   ↓
Deduplication (GUID tracking)
   ↓
Message formatting
   ↓
Discord API (Bot)
   ↓
Channel notifications
```

## 📦 Tech Stack
- Go (standard library heavy)
- RSS parsing: `github.com/mmcdole/gofeed`
- HTML → Markdown: `github.com/JohannesKaufmann/html-to-markdown`
- Discord REST API (v10)

## ⚙️ How it works
1. The program fetches a configured RSS feed
2. Each item is checked against the last stored GUID
3. New items are formatted into readable messages
4. Messages are sent to a Discord channel
5. The latest processed GUID is saved locally to avoid duplicates

## 💾 State Management

Each feed maintains a simple state file:
```
logs/
  cve-daily
```
This file stores the last processed GUID.

This allows:
- stateless execution (safe for cron jobs)
- deduplication across runs
- multi-feed extensibility

## 🔔 Example Output
```
CVE-2026-XXXX
Remote code execution vulnerability in ...
https://example.com/cve
Published: 2026-04-24
CVE Daily Feed
```

## ⏱ Running as a Cron Job
### Linux / WSL cron example
Run every 15 minutes:
```
*/15 * * * * /path/to/rss-bot
```
Make sure:
- The binary is executable
- Working directory is set correctly (for logs/ + .env)

## Windows (Task Scheduler alternative)
Run:
```
rss-bot.exe
```
Trigger:
- Every 15 minutes

## 🔐 Configuration

A simple `.env` file is used:
```
DISCORD_CHANNEL_ID
DISCORD_BOT_TOKEN
```

Example:
```
1399719606184050701
YOUR_BOT_TOKEN_HERE
```

## 📁 Project Structure
```
.
├── main.go
├── logs/                # persistent state per feed
├── .env                 # bot credentials
```
