# COPE GAMING (Cloud gaming platform)

This project allows users to play computer games on browsers, without download or installation, just at the push of a button.
This means that users can play AAA games on their low-spec computers or smartphones or even smart TVs anywhere and anytime.
The only requirement is fast and reliable internet.

The difference between this project and [other cloud-gaming projects](https://github.com/hxt365/cloud-gaming) is that it allows users to offer their computers (providers) to other users (players) and make profits.
That means games are run in computers of providers and streamed to browsers of players.

## Demo

![Demo](https://user-images.githubusercontent.com/16115992/159097057-be003206-8d2c-49f4-973b-7a27fccb8559.gif)

## Usage

### Prerequisite
- Golang >= 1.16
- NodeJs and npm (for Web UI)
- Docker and docker-compose (for providers)
- Current user on your computer has permissions to run Docker
- For now, we only support Linux providers

### How to run

- To run Coordinator service:
```bash
cd coordinator/
go run main.go
```

- To be a provider:
```
cd provider/
./run.sh
```

- To run UI:
```bash
cd web/
npm install
npm start
```

## Design

This project is inspired by [cloudmorph](https://github.com/giongto35/cloud-morph) and [drova.io](https://drova.io/).
The basic idea of the project is running games in Wine within Docker containers.
The video and audio from games are captured by Xvfb, Pulseaudio and processed by ffmpeg and then streamed to browsers of users using WebRTC.
Besides that, input from users (e.g. mouse clicks, keyboard events) are also captured and delivered to Syncinput using WebRTC Data channel.
Syncinput is a process that receives those input and simulate relevant events for the games using WinAPI.
