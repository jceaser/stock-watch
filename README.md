# Stock Watch

A command-line tool for managing stock portfolios, fetching current stock details, and generating reports.

## What is it for?
Stock Watch is designed to help investors track their stock holdings, manage inventory (quantity and cost basis), and receive alerts when prices hit specific thresholds. It allows you to:
- Maintain a database of stocks you own.
- Fetch real-time market data.
- Set automated alerts for high and low prices.
- Generate reports for easy analysis.

## Getting Started

### Prerequisites
- [Go](https://golang.org/) (version 1.18 or later) built with 1.25.14.
- If testing finnhub you will need a valid API key (configured in the source code).

### Build
To build the application, navigate to the root directory and run:

    go build stock-watch.go

## Usage Guide

The application supports two primary modes of operation: **Interactive Mode** (for manual
exploration) and **Batch Mode** (for automation or scripting).

### 1. Interactive Mode
Run the application without any flags to enter the interactive shell. This allows you to type commands and see immediate feedback. Use the `help` command to learn what actions are available. Additionally each command also takes a "help" parameter to get details about just that command.

	go run .

Once inside, you will see a > prompt. You can type commands like help, fetch, or report-full.

* Commands can be chained: You can run multiple commands at once using a semicolon (e.g., fetch; report-full).
* History: The application saves your command history to a temporary file for easy retrieval. Use the up and down arrow to pick commands.

### 2. Batch Mode

Use Batch Mode when you want to execute a single command and exit immediately. This is ideal for cron
jobs or scripts.

	go run . -batch-mode -initial "fetch"

-----

## Available Commands

### Portfolio Management

| Action | Description | Example |
| ------ | ----------- | ------- |
| create id desc note  | Add a new stock to your portfolio.       | create AAPL "Apple Inc" "Buy on dip" |
| read id              | View details for a specific stock.       | read AAPL
| update id desc note  | Update description or notes for a stock. | update AAPL "New Desc" "New Note"
| delete id            | Remove a stock from your records.        | delete AAPL


### Inventory & Alerts

| Command | Description | Example |
| ------- | ----------- | ------- |
| stock-update id row count value | Update a specific inventory row.  | stock-update AAPL 0 15 160.00
| alert id low high | Manually set price alert thresholds.            | alert AAPL 100 200
| assume-alerts id  | Automatically set alerts based on 52-week data. | assume-alerts AAPL
                                    
### Data & Reporting

Reports take a "filter" which can be either alerts or invitory
* alerts - only show records with an alert
* invitory - only show records with invitory

| Command | Description | Example |
| ------- | ----------- | ------- |
| fetch                   | Fetch the latest market prices for all stocks. | fetch
| report-full filters     | Display a full report in the terminal.         | report-full 
| report-csv filters      | Generate a CSV-formatted report.               | report-csv invitory
| report-markdown filters | Generate a Markdown-formatted report.          | report-markdown alerts
| report-shell filters    | Generate a stylized report for the terminal.   | report-shell alerts,invitory

----
Copyright (c) 2026 All rights reserved Thomas Cherry. Thomas.Cherry@gmail.com
