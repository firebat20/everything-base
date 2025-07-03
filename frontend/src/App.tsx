import { Button } from "@/components/ui/button"
import {
  Menubar,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarTrigger,
} from "@/components/ui/menubar"
import { DemoComponents } from './components/DemoComponents' // Import the new component
import React from "react"

function App() {
  const [count, setCount] = React.useState(0)

  

  return (
    <div className="min-h-screen bg-white grid place-items-center mx-auto py-8">
      <div className="text-blue-900 text-2xl font-bold flex flex-col items-center space-y-4">
      <Menubar>
        <MenubarMenu>
          <MenubarTrigger>File</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>Option 1</MenubarItem>
            <MenubarItem>Option 2</MenubarItem>
            <MenubarItem>Option 3</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu>
          <MenubarTrigger>Debug</MenubarTrigger>
          <MenubarContent>
            <MenubarItem>Option 1</MenubarItem>
            <MenubarItem>Option 2</MenubarItem>
            <MenubarItem>Option 3</MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </Menubar>
      <DemoComponents />

        <h1>Vite + React + TS + Tailwind + shadcn/ui</h1>
        <Button onClick={() => setCount(count + 1)}>Count up ({count})</Button>
      </div>
    </div>
  )
}

export default App
