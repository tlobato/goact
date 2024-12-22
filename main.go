package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	//GO MODULE
	fmt.Println("====== Goact ======")
	fmt.Println("Setup your Go and React project with ease!")

	fmt.Println("\nGO MODULE SETUP")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your github.com account name: ")
	account, _ := reader.ReadString('\n')
	account = strings.TrimSpace(account)

	fmt.Print("Enter your project's name: ")
	project, _ := reader.ReadString('\n')
	project = strings.TrimSpace(project)
	project = fmt.Sprintf("github.com/%s/%s", account, project)

	if err := os.Remove("./go.mod"); err != nil && !os.IsNotExist(err) {
		fmt.Println("Failed to remove go.mod:", err)
		return
	}

	cmd := exec.Command("go", "mod", "init", project)

	err := cmd.Run()
	if err != nil {
		fmt.Println("failed to create go.mod:", err)
		return
	}

	fmt.Println("Go module setup:", project)
	cmd = exec.Command("go", "get", "github.com/joho/godotenv")
	cmd.Run()

	//SETUP GO
	os.Mkdir("./server", 0755)
	handlersFile := `
	//Example code
	package server

	import (
	"net/http"
	"encoding/json"
	)

	func ApiRoutes() {
		http.HandleFunc("/api/hello", handleHello)
	}

	type HelloResponse struct {
		Message string ` + "`json:\"message\"`" + `
	}

	func handleHello(w http.ResponseWriter, r *http.Request) {
		response := HelloResponse{Message: "Hello, World!"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}`
	os.WriteFile("./server/handlers.go", []byte(handlersFile), 0644)

	os.Remove("./main.go")
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

	cmd = exec.Command("air", "init")
	err = cmd.Run()
	if err != nil {
		fmt.Println("failed to create air.toml:", err)
	}

	//SETUP FRONTEND
	fmt.Println("\nFRONTEND SETUP")
	fmt.Print("Your package manager (pnpm/bun): ")
	pkgManager, _ := reader.ReadString('\n')
	pkgManager = strings.TrimSpace(pkgManager)

	fmt.Print("Use TypeScript? (Y/n): ")
	useTs, _ := reader.ReadString('\n')
	useTs = strings.TrimSpace(useTs)

	template := "react-ts"
	if useTs == "n" {
		template = "react"
	}
	cmd = exec.Command(pkgManager, "create", "vite@latest", "./client", "--template", template)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to create vite app:", err)
		return
	}

	fmt.Printf("Use Tailwind? (Y/n): ")
	useTailwind, _ := reader.ReadString('\n')
	useTailwind = strings.TrimSpace(useTailwind)

	if useTailwind != "n" {
		cmd = exec.Command(pkgManager, "install", "-D", "tailwindcss", "postcss", "autoprefixer")
		cmd.Dir = "./client"
		err = cmd.Run()
		if err != nil {
			fmt.Println("failed to install tailwind:", err)
			return
		}

		cmd = exec.Command("npx", "tailwindcss", "init", "-p")
		cmd.Dir = "./client"
		err = cmd.Run()
		if err != nil {
			fmt.Println("failed to install tailwind:", err)
			return
		}

		os.WriteFile("./client/tailwind.config.js", []byte(`/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}`), 0644)

		os.WriteFile("./client/src/index.css", []byte(`@tailwind base;
@tailwind components;
@tailwind utilities;`), 0644)

	}

	os.Remove("./client/src/App.css")
	os.Remove("./client/src/assets")
	os.Remove("./client/src/App.tsx")

	os.WriteFile("./client/src/main.tsx", []byte(`import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
   <h1 className="text-2xl">Hello from Goact!</h1> 
  </StrictMode>,
)`), 0644)

	cmd = exec.Command(pkgManager, "install")
	cmd.Dir = "./client"
	err = cmd.Run()
	if err != nil {
		fmt.Println("failed to install dependencies:", err)
		return
	}

	fmt.Println("Setup complete!")

	fmt.Println("\n Use Guide:")
	fmt.Println("1. To run the Go dev server: `air .`")
	fmt.Printf("2. To run the React dev server: `cd client && %s run dev`\n", pkgManager)
	fmt.Printf("3. To build frontend: `cd client && %s build`\n", pkgManager)
	fmt.Printf("4. To run production: `cd client && %s build && cd .. && go run . \n", pkgManager)
}
