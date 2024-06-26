# delve-finder
This repository is inspired by Paul Graham, computer scientist, and co-founder of Y Combinator, regarding his concerns about those who "delve" too much when messaging him. The aim is to automatically detect and delete emails from these delvers for people who face similar problems as Paul. For context:

![Image Alt Text](Paul_Graham_Tweet.PNG)



## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Getting Started

### Prerequisites

You should have [Go installed on your machine](https://go.dev/doc/install). You need to set up a [google cloud project](https://console.cloud.google.com/projectcreate).  [Watch the video](https://www.youtube.com/watch?v=-rcRf7yswfM) on how to quickly set up your project. You only need to watch the video up until [6:56](https://youtu.be/-rcRf7yswfM?t=416).


### Installation 

After installing Go and setting up your project run the following commands:

1. You need to initialize your go module:
   ```bash
   go mod init [module name]
   ```
2. Get the Gmail API Go client library and OAuth2.0 package 
  ```bash
   go get google.golang.org/api/gmail/v1
   go get golang.org/x/oauth2/google
  ```

### Usage
1.Run your Go program:
   ```bash
   go run find_delver.go
   ```
2.Follow the insturctions printed on your terminal to get the authorization code

## Future Plans
Add a LLM detector in order to properly find this delvers lurking in emails.
## License

This project is licensed under the [MIT License](MIT-LICENSE.txt)



