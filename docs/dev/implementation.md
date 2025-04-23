# 🗺️ Infinite Experiment Bot – Sequential Implementation Plan (Updated)

## ✅ Step 1: Minimal Flight Stats (Quick Wins)
- `/my_stats` command (basic implementation)
- Fetch synced flight log summary (hours, routes, recent date)
- Ensure flight logs are storing properly in DB from Airtable/API

---

## ✅ Step 2: Interactive Logbook
- Implement `/logbook [@user]` command
- Add pagination support via Discord buttons (prev/next)
- Respond with rich embed per flight (date, mode, duration, route)
- Store paginated tokens (e.g., in memory or encoded cursor)

---

## ✅ Step 3: Advanced History Filtering
- Implement `/flight_history [mode] [@user]`
- Accept mode (scheduled/charter/all)
- Use similar pagination system
- Fallback to user context if no mention

---

## ✅ Step 4: Web-based Map Integration (Move Up in Sequence)
- Vue 3 map viewer at route `/map?route[]=...`
- Use Deck.gl for 3D rendering of flight paths on a basemap
- Provide one-time use route list via backend link
- Link from `/logbook` and `/flight_history` embeds: “View in 3D Map”
- Optional animated flights using TripsLayer or ArcLayer in Deck.gl

---

## ✅ Step 5: Register, Sync, and Profile
- `/register <community_username>`
- `/sync_if`, `/verify_route` flow
- `/profile`, `/my_rank`
- Web UI: server/VA setup, Airtable link form

---

## ✅ Step 6: Crew Manager Bot Tools
- `/cm_stats`
- `/va flights <day>`
- `/va stats`, `/va roster`
- Require user role manager/crew role validation

---

## ✅ Step 7: Crew Manager Web Panel
- Vue 3 dashboard: server config, roster table, audit logs
- Admin login (API Key or session token)
- Filters: user, mode, hours, route
- Manual override tools for rank/status

---

## ✅ Step 8: Event Engagement
- `/event create`
- `/event list`, `/event rsvp`
- Web RSVP viewer and calendar
- Notifications + reminders via bot

---

## ✅ Step 9: Leaderboards & Visual Analytics
- `/leaderboard top_hours`, `/top_routes`
- Backend flight aggregations
- Add web viewable charts/graphs (e.g., hours/week)
- Top pilots by role/rank

---

## ✅ Step 10: Admin Tools & QoL
- `/set_role`, `/assign_rank`, `/flight_mode set`
- Error formatting, health check, `/refresh`
- Embed templates + response consistency helpers

---

Let me know if you'd like a Vue map mockup next or want to break this into tasks for tracking. 🛩️

