## Goact
Best Go + React setup for web development.
#### Featuring:
- Go's [Fiber](https://gofiber.io) library for easier routes and authentication.
- [React Router](https://reactrouter.com) configured by default.
- [Tailwind](https://tailwindcss.com) for easier styling.
- [React Compiler](https://react.dev/learn/react-compiler) for a better client performance. (not stable yet)
- Air + Vite development mode.
- All in one deployment (Go serves the front-end in production).

## Usage
- Clone the repo: `git clone https://github.com/tlobato/goact && cd goact`
- Use `go run . [dir]`, Go through the CLI and enjoy your personalized setup
---
- Dev server: `cd client && npm dev && cd .. && air .`<br/>
  It will automatically proxy API routes from port :3000 to :5173 through vite's proxy </br>
- Production: `cd client && npm build && cd .. && go run .`<br/>
  Go serves the built app at client/dist
