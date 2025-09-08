# OpenMeet (Side Project)

**OpenMeet** is a side project inspired by Zoom and Google Meet.  
It’s built on top of [LiveKit](https://github.com/livekit/livekit) to provide **real-time audio/video meetings** with **OAuth2 authentication**.

---

## ✨ Features

### Core Features
- [x] Secure authentication with Google Sign-In
- [x] Create and join meeting rooms
- [x] Real-time video and audio communication
- [x] Room persistence and management
- [x] Secure token-based access

### Host Controls
- [ ] Participant Management
  - [ ] Approve/remove participants
  - [ ] Mute/unmute participants
  - [ ] Assign co-hosts
  - [ ] View participant list
  - [ ] Kick participants
- [ ] Meeting Controls
  - [ ] Terminate meeting for all
  - [ ] Lock room to prevent new joins
  - [ ] End meeting and save recording
  - [ ] Control screen sharing permissions

### Participant Features
- [ ] Interactive Features
  - [ ] Raise hand feature
  - [ ] Emoji reactions
  - [ ] Chat messages
  - [ ] File sharing
  - [ ] Image sharing
  - [ ] Custom backgrounds
- [ ] Meeting Controls
  - [ ] Mute/unmute audio
  - [ ] Enable/disable video
  - [ ] Share screen
  - [ ] Change audio/video devices

### Calendar Integration
- [ ] Google Calendar Integration
  - [ ] Schedule meetings
  - [ ] Send calendar invites
  - [ ] Automatic meeting reminders
  - [ ] Join directly from calendar
- [ ] Apple Calendar Integration
  - [ ] iCalendar support
  - [ ] Add to calendar feature
  - [ ] Meeting notifications

### Session Management
- [ ] Long Session Support
  - [ ] Auto-reconnect on disconnect
  - [ ] Session persistence
  - [ ] Connection quality monitoring
  - [ ] Bandwidth optimization
- [ ] Recording
  - [ ] Cloud recording
  - [ ] Local recording
  - [ ] Recording management
  - [ ] Download recordings

### Mobile Support
- [ ] Responsive Web Interface
  - [ ] Mobile-first design
  - [ ] Touch-optimized controls
  - [ ] Portrait/landscape support
  - [ ] PWA support
- [ ] Native Mobile Apps
  - [ ] iOS App Store release
  - [ ] Android Play Store release
  - [ ] Push notifications
  - [ ] Background audio support

### Additional Features
- [ ] Breakout Rooms
- [ ] Virtual Backgrounds
- [ ] Noise Cancellation
- [ ] Live Captions
- [ ] Meeting Analytics
- [ ] Custom Branding Options

---

## 🛠️ Tech Stack
- **Backend:** Go (Golang)
- **Real-time:** LiveKit (WebRTC)
- **Auth:** OAuth2 (Google/GitHub provider)
- **Database:** PostgreSQL
- **Frontend:** React

---

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/open-meet.git

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run the server
go run cmd/main.go
```

## Environment Variables

```
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
ALLOWED_ORIGINS=your_allowed_origins
LIVEKIT_API_KEY=your_livekit_api_key
LIVEKIT_API_SECRET=your_livekit_api_secret
LIVEKIT_SERVER=your_livekit_server_url
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

### Q4 2025
- [ ] Host controls implementation
- [ ] Calendar integration
- [ ] Basic mobile responsive UI

### Q1 2026
- [ ] Native mobile app development
- [ ] Advanced session management
- [ ] Recording features

### Q2 2026
- [ ] App store submissions
- [ ] Enhanced interactive features
- [ ] Analytics and monitoring

## Contact

For questions and support, please [open an issue](https://github.com/yourusername/open-meet/issues) or contact the maintainers.
