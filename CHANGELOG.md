# Changelog

## [0.1.0] - 2026-05-30
### Added
- **Frontend-Backend Integration**: fully connected the React frontend with the Go Gin backend.
- **Game Engine**: wired up the Dice Roller (`useGameStore.roll`) with backend API (`/api/game/roll`) to perform remote geospatial random rolling.
- **AI Chat Module**: connected the `RightDrawer` chat UI with the `/api/chat` endpoint.
- **Cyber Check-in**: added image upload and `api.generateImage` and `api.createCheckin` call to `CityDetailPanel.tsx` for Cyberpunk check-ins.
- **Achievement Wall**: connected `AchievementPage.tsx` with `/api/users/:id/achievements` and integrated it as an overlay via `Navbar`.
- **Backend Architecture**: implemented Go + Gin + GORM backend supporting MySQL and dynamic seeding.
- **Frontend Architecture**: built React + Vite + Zustand + TailwindCSS app with AMap 3D and Pannellum 360 viewer.

### Fixed
- Fixed TypeScript errors related to mismatched fields between Mock data (`city.figures`) and real API (`city.characters`).
- Fixed missing arguments in `roll()` calls within `Sidebar.tsx` and `DiceConsole.tsx`.
