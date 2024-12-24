## Goact
Best Go + React setup for web development.
#### Featuring:
- Go's [Fiber](https://gofiber.io) library for easier routes and authentication.
- [React Router](https://reactrouter.com) configured by default.
- [Tailwind](https://tailwindcss.com) for easier styling.
- [React Compiler](https://react.dev/learn/react-compiler) for a better client performance.
- Air + Vite development mode.
- All in one deployment (Go serves the front-end in production).

## Usage
- Clone the repo: `git clone https://github.com/tomaslobato/goact`
- Use `go run . [dir]`, Go through the CLI and enjoy your personalized setup
---
- Dev server: `cd client && pnpm dev && cd .. && air .`
  It will automatically proxy API routes from port :3000 to :5173 through vite's proxy 
- Production: `cd client && pnpm build && cd .. && go run .`
  Go serves the built app at client/dist
