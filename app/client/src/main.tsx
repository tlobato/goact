import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
   <h1 className="text-2xl">Hello from Goact!</h1>
  </StrictMode>,
)