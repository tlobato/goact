package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [dir]")
		fmt.Println("`[dir]` refers to the folder you want your project to be, write it without brackets")
		os.Exit(1)
	}

	folder := os.Args[1]
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		fmt.Println("Failed to create or access folder:", err)
		os.Exit(1)
	}

	removingFiles := []string{"go.mod", "main.go", "readme.md"}
	for _, file := range removingFiles {
		os.RemoveAll(file)
	}

	err = os.Chdir(folder)
	if err != nil {
		fmt.Println("Failed to change directory:", err)
		os.Exit(1)
	}

	//SETUP GO
	fmt.Println("====== Goact ======")
	fmt.Println("Setup your Go and React project with ease!")

	fmt.Println("\nGO MODULE SETUP")
	reader := bufio.NewReader(os.Stdin)

	os.WriteFile("./.gitignore", []byte(".env\nnode_modules/\ndist/\n"), 0644)

	fmt.Print("Enter your github.com account name: ")
	account, _ := reader.ReadString('\n')
	account = strings.TrimSpace(account)

	fmt.Print("Enter your project's name: ")
	project, _ := reader.ReadString('\n')
	project = strings.TrimSpace(project)
	project = fmt.Sprintf("github.com/%s/%s", account, project)

	cmd := exec.Command("go", "mod", "init", project)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to create go.mod:", err)
		return
	}

	fmt.Println("Go module setup:", project)
	cmd = exec.Command("go", "get", "github.com/joho/godotenv")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}

	os.Mkdir("./server", 0755)
	fmt.Println("\nYou can choose fiber as a simpler way to setup your backend, if you want the standard library write n")
	fmt.Print("Use Fiber? [Y/n]: ")
	useFiber, _ := reader.ReadString('\n')
	useFiber = strings.TrimSpace(useFiber)
	if useFiber != "n" {
		cmd = exec.Command("go", "get", "github.com/gofiber/fiber/v2")
		if err = cmd.Run(); err != nil {
			fmt.Println("Failed to install Fiber: ", err)
		}
		handlersFile := `package server

import (
		"github.com/gofiber/fiber/v2"
)

func HandleHello(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "hello!"})
}`

		os.WriteFile("./server/handlers.go", []byte(handlersFile), 0644)

		mainFile := fmt.Sprintf(`package main

import (
		"log"

		"%s/server"

		"github.com/joho/godotenv"
		"github.com/gofiber/fiber/v2"
		"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
		godotenv.Load(".env")

		app := fiber.New()

		//handle frontend routes
		app.Static("/", "./client/dist")

		api := app.Group("/api", cors.New(cors.Config{
			AllowOrigins: "http://localhost:3000", //change this later
			AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		}))
		api.Get("/hello", server.HandleHello)

		log.Fatal(app.Listen(":3000"))
}
`, project)
		os.WriteFile("./main.go", []byte(mainFile), 0644)
	} else {
		handlersFile := `//Example code
package server

import (
		"net/http"
		"encoding/json"
)

func ApiRoutes() {
		http.HandleFunc("/api/hello", handleHello)
}

type HelloResponse struct {
		Message string ` + "`json:`" + `"message"` + `
}

func handleHello(w http.ResponseWriter, r *http.Request) {
		response := HelloResponse{Message: "Hello, World!"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
}`
		os.WriteFile("./server/handlers.go", []byte(handlersFile), 0644)

		mainFile := fmt.Sprintf(`package main

import (
		"log"
		"os"
		"net/http"
		"fmt"

		"%s/server"
		"github.com/joho/godotenv"
)

func main () {
		godotenv.Load(".env")

		//handle frontend routes
		fs := http.FileServer(http.Dir("./client/dist/"))
		http.Handle("/", fs)

		server.ApiRoutes()

		port := os.Getenv("PORT")
		if port == "" {
			port = ":3000"
		}
		fmt.Println("Server running on port", port)
		fmt.Println("http://localhost:3000 - go to /api/hello to test the Go API")
		log.Fatal(http.ListenAndServe(port, nil))
}`, project)
		os.WriteFile("./main.go", []byte(mainFile), 0644)
	}

	os.Mkdir("./client", 0755)
	cmd = exec.Command("npm", "create", "vite@latest", ".", "--", "--template", "react-ts", "-y")
	cmd.Dir = "./client"
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to create Vite project:", err)
		return
	}

	fmt.Println("\nFRONTEND SETUP")
	fmt.Printf("Use React Router? (Y/n): ")
	useRR, _ := reader.ReadString('\n')
	useRR = strings.TrimSpace(useRR)

	if useRR == "n" {
		os.WriteFile("./client/src/main.tsx", []byte(`import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
   <h1 className="text-2xl">Hello from Goact!</h1>
  </StrictMode>,
)`), 0644)
	} else {
		cmd = exec.Command("npm", "install", "react-router")
		cmd.Dir = "./client"
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to install React Router:", err)
		}

		os.WriteFile("./client/src/main.tsx", []byte(`import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { BrowserRouter, Link, Outlet, Route, Routes } from 'react-router'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={
          <div className='flex h-screen'>
            <aside className='h-full w-56 bg-zinc-900 text-white'>
              <nav className='flex flex-col gap-1.5 p-2 text-lg'>
                <Link to="/" className='w-full cursor-pointer hover:bg-zinc-700 rounded bg-zinc-800 p-1.5'>Home</Link>
                <Link to="/about" className='w-full cursor-pointer hover:bg-zinc-700 rounded bg-zinc-800 p-1.5'>About</Link>
                <Link to="/usage" className='w-full cursor-pointer hover:bg-zinc-700 rounded bg-zinc-800 p-1.5'>Usage</Link>
              </nav>
            </aside>
            <Outlet />
          </div>}>
          <Route index element={<main className='bg-zinc-700 flex-1 text-white flex flex-col items-center pt-16'><h1 className='text-6xl font-extrabold'>Home</h1></main>} />
          <Route path="/about" element={<main className='bg-zinc-700 flex-1 text-white flex flex-col items-center pt-16'><h1 className='text-6xl font-extrabold'>About</h1></main>} />
          <Route path="/usage" element={<main className='bg-zinc-700 flex-1 text-white flex flex-col items-center pt-16'><h1 className='text-6xl font-extrabold'>Usage</h1></main>} />
        </Route>
      </Routes>
    </BrowserRouter>
  </StrictMode>,
  )`), 0644)
	}

	fmt.Printf("Use Tailwind? (Y/n): ")
	useTailwind, _ := reader.ReadString('\n')
	useTailwind = strings.TrimSpace(useTailwind)

	if useTailwind != "n" {
		cmd = exec.Command("npm", "install", "tailwindcss", "autoprefixer", "postcss")
		cmd.Dir = "./client"
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to install Tailwind:", err)
		}
		cmd = exec.Command("npx", "tailwindcss", "init")
		cmd.Dir = "./client"
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to init Tailwind:", err)
		}

		tailwindConfig := `/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}`
		os.WriteFile("./client/tailwind.config.js", []byte(tailwindConfig), 0644)

		postcssConfig := `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}`
		os.WriteFile("./client/postcss.config.js", []byte(postcssConfig), 0644)

		indexCSS := `@tailwind base;
@tailwind components;
@tailwind utilities;`
		os.WriteFile("./client/src/index.css", []byte(indexCSS), 0644)
	}

	viteConfig := `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/api": {
        target: 'http://localhost:3000',
        changeOrigin: true
      }
    }
  }
})`

	os.WriteFile("./client/vite.config.ts", []byte(viteConfig), 0644)

	fmt.Println("\nSetup complete!")
	fmt.Println("\nNext steps:")
	fmt.Printf("1. cd %s/client && npm install\n", folder)

	fmt.Println("2. Run Golang dev mode with `air` or `go run .`")
	fmt.Println("3. And run Vite dev mode with `cd client && npm run dev`")

	fmt.Println("Happy coding!")
}
