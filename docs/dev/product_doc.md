# Infinite Experiment Discord Bot – Combined Product Specification

## 🧭 Overview
The Infinite Experiment Discord Bot empowers pilots and virtual airlines within the Infinite Flight simulator community to register, track, and manage flight operations directly from Discord. It integrates with Infinite Flight APIs, Airtable, and a Go-based backend to serve authenticated, role-based, real-time aviation data across both structured airline routes and flexible charter modes.

---

## ✈️ Core Capabilities

### 1. 🚀 Registration & Authentication
- `/register <community_username>`
  - Registers the user using their IF Community username and links them to the VA associated with the server.

- OTP-based Verification:
  - `/sync_if <username>` – Fetches latest IF profile + flight log.
  - `/verify_route <last_flight_route>` – Confirms identity via last known route.

- Modes supported:
  - Scheduled Flights
  - Charter Flights

### 2. 🧑‍✈️ Pilot Features
- `/profile` – Displays total hours, current rank, recent activity, linked VA.
- `/logbook` – Returns last 5–10 flights across both flight modes.
- `/my_rank` – Shows current rank, earned hours, and required hours for promotion.
- `/my_stats` – Personalized flight breakdown by mode, region, duration.
- `/flight_history [mode] [range]` – Filtered view of flights based on mode (scheduled/charter) and time range.

### 3. 🧑‍✈️ Crew Manager Features
- `/cm_stats` – Aggregates flight logs by mode per pilot. Supports flight count, hours, VA activity.
- `/va flights <day>` – Fetches all submitted PIREPs for a day (or recent N flights).
- `/va stats` – Summarized overview of the VA's total flights, active pilots, and recent activity.
- `/va roster` – Displays pilots with rank, hours, roles, and recent flights.

### 4. 🎉 Community Features
- `/event create` – Schedule community flights or events with full description, route, time.
- `/event list` – Lists upcoming events with RSVP tracking.
- `/event rsvp` – Enables pilot RSVP tracking for coordination.

### 5. 📈 Analytics & Leaderboards
- `/leaderboard top_hours` – Top pilots based on flight hours.
- `/leaderboard top_routes` – Most popular routes.
- `/leaderboard top_active` – Most active pilots by flights per week.
- `/flight_map <callsign>` – Shows recent flight trajectories on a rendered map.

### 6. 🔧 Admin/Manager Tools
- `/set_role <@user> <role>` – Assigns bot-level VA roles (pilot, CM, manager).
- `/refresh <@user>` – Re-syncs user data from backend/API.
- `/assign_rank <@user> <rank>` – Allows manual override in case of API mismatch.
- `/flight_mode set <charter/scheduled>` – Sets preferred default for logging.

### 7. 🧰 Quality-of-Life Features
- `/help` – Categorized command list with brief descriptions.
- `/health` – System and service diagnostics.
- Automatic callsign parsing from usernames.
- Auto-detect flight mode based on context.
- Bot DMs to notify about promotion eligibility.
- Discord buttons for RSVP, confirmations, and quick actions.

---

## 📦 Backend/Infra Design Notes
- Backend written in Go with REST + gRPC endpoints.
- PostgreSQL for core data.
- Airtable for per-VA dynamic configs and custom overrides.
- Real-time data from Infinite Flight APIs.
- Discord API integration via TypeScript bot.
- All requests authenticated via JWT (UI) or API Key (Bot).
- Users can belong to multiple VAs, each scoped by server ID.

---

## 🔮 Future Enhancements
- Visual dashboard with rank progress, flight calendar.
- Custom route validation engines for scheduled modes.
- Route recommendation engine per rank level.
- Voice and chat-based briefings before events.
- Integration with NocoDB for simplified config management.
- PIREP editing and resubmission features.
- Automatic time zone conversion for events.

---


