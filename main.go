package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/otiai10/copy"
)

func main() {
	//GO
	reader := bufio.NewReader(os.Stdin)
	folder := "."
	if len(os.Args) < 2 {
		fmt.Print("Where will your project be?: (.) ")
		folder, _ = reader.ReadString('\n')
		folder = strings.TrimSpace(folder)
	} else {
		folder = os.Args[1]
	}

	os.Mkdir(folder, 0755)

	fmt.Print("Use Fiber? [Y/n]: ")
	fiber, _ := reader.ReadString('\n')
	fiber = strings.TrimSpace(fiber)

	if folder == "." {
		os.Remove("main.go")
		os.Remove("go.mod")
		os.Remove("go.sum")
		os.Remove("readme.md")
	}

	var err error
	if fiber == "n" {
		err = copy.Copy("./code/vanilla", folder)
		if err != nil {
			fmt.Println("Dir contains files:", err)
		}
	} else {
		err = copy.Copy("./code/fiber", folder)
		if err != nil {
			fmt.Println("Dir contains files:", err)
		}
	}

	os.RemoveAll("./code")

	os.Chdir(folder)

	fmt.Print("Give a name to your Go module (Goact): ")
	module, _ := reader.ReadString('\n')
	module = strings.TrimSpace(module)

	var cmd *exec.Cmd
	if module == "" {
		cmd = exec.Command("go", "mod", "init", "Goact")
	} else {
		cmd = exec.Command("go", "mod", "init", module)
	}
	cmd.Run()

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Run()

	//FRONTEND
	fmt.Print("JS package manager (NPM/pnpm/bun): ")
	pm, _ := reader.ReadString('\n')
	pm = strings.TrimSpace(pm)
	if pm != "pnpm" && pm != "bun" {
		pm = "npm"
	}

	if pm == "npm" {
		cmd = exec.Command("npx", "create-vite@latest", "./client", "--template", "react-ts")
	} else {
		cmd = exec.Command(pm, "create", "vite@latest", "./client", "--template", "react-ts")
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error creating vite application:", err)
		return
	}

	if _, err := os.Stat("client"); os.IsNotExist(err) {
		fmt.Println("Client directory was not created successfully")
		return
	}

	err = os.Chdir("client")
	if err != nil {
		fmt.Println("Error changing to client directory:", err)
		return
	}

	fmt.Print("Use React Router? [Y/n]: ")
	rr, _ := reader.ReadString('\n')
	rr = strings.TrimSpace(rr)

	if rr != "n" {
		cmd = exec.Command(pm, "i", "react-router@latest")
		cmd.Run()
	}

	fmt.Print("Use Tailwind? [Y/n]: ")
	tailwind, _ := reader.ReadString('\n')
	tailwind = strings.TrimSpace(tailwind)

	if tailwind != "n" {
		cmd = exec.Command(pm, "i", "-D", "tailwindcss", "autoprefixer", "postcss")
		cmd.Run()

		if pm == "npm" {
			cmd = exec.Command("npx", "tailwindcss", "init")
		} else if pm == "pnpm" {
			cmd = exec.Command("pnpm", "dlx", "tailwindcss", "init")
		} else if pm == "bun" {
			cmd = exec.Command("bunx", "tailwindcss", "init")
		}

		os.Remove("./src/App.css")

		os.WriteFile("./postcss.config.js", []byte(`export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  }
}`), 0644)

		os.WriteFile("./src/index.css", []byte(`@tailwind base;
@tailwind components;
@tailwind utilities;
`), 0644)

		os.WriteFile("./tailwind.config.js", []byte(`/** @type {import('tailwindcss').Config} */
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

		os.WriteFile("./src/App.tsx", []byte(`export default function App() {
  return (
    <>
      <h1 className="text-2xl">Hello from Goact!</h1>
    </>
  );
}`), 0644)
	}

	if rr != "n" {
		os.WriteFile("./src/main.tsx", []byte(`import { createRoot } from "react-dom/client"
import "./index.css"
import App from "./App.tsx"
import { BrowserRouter, Route } from "react-router"
import { Routes } from "react-router"

createRoot(document.getElementById("root")!).render(
  <BrowserRouter>
    <Routes>
      <Route index element={<App />} />
    </Routes>
  </BrowserRouter>,
)`), 0644)
	}

	os.Remove("./README.md")

	os.WriteFile("./vite.config.ts", []byte(`import { defineConfig } from 'vite'
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
})`), 0644)

	fmt.Println("All Done!")
	if folder != "." {
		fmt.Printf("Use `cd %s` to find your project\n", folder)
	}
	fmt.Println("Save go files to fix undefined errors")
	fmt.Print("\n")

	fmt.Println("Happy Coding!")
}
