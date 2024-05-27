# Chatbot with Image Uploading and Retrieval

## Project Overview

Developed a chatbot to handle image uploading and retrieval, utilizing WebSocket communication and OpenAI integration for enhanced interactions

## Features

- Real-time chatting via WebSocket.
- Image upload and retrieval using REST api's.
- Conversation context management for OpenAI's GPT interactions.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

What things you need to install the software and how to install them:

```bash
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
go get github.com/go-sql-driver/mysql
go get github.com/joho/godotenv
go get github.com/sashabaranov/go-openai
```

## Installing

A step-by-step series of examples that tell you how to get a development environment running, before that please make sure that you are ran `setup.sql` script in the `chatbot` db and have specified all configurations along with openai key in `.env` file:
```bash
# Clone the repository
git clone https://github.com/yourusername/chatbot.git

# Navigate to the project directory
cd chatbot

# Build the project (optional)
go build

# Run the server
go run cmd/chatbot/main.go
```

## Usage
### Using Api's
* POST /upload
Upload an image with an associated identifier.
<img width="1313" alt="image" src="https://github.com/Indian-boult/chatbot/assets/64060864/bbe53852-feb6-482b-9e16-2cb14f3ed839">


* GET /retrieve
Retrieve the metadata for an uploaded image.
<img width="1289" alt="image" src="https://github.com/Indian-boult/chatbot/assets/64060864/d8cf96a3-6e85-471f-bc6f-d8de29134401">

### Using UI
Here I have created the simple UI to access the chatbot and use it's feature to upload image, retrieve image, talk about image.
<img width="1033" alt="image" src="https://github.com/Indian-boult/chatbot/assets/64060864/8b1f0d60-5141-4407-ac48-373cbd9f75ef">

## License
This project is licensed under the MIT License - see the LICENSE.md file for details
