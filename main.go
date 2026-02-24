package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "move":
		handleMove()
	case "scroll":
		handleScroll()
	case "click":
		handleClick()
	case "press":
		handlePress()
	case "release":
		handleRelease()
	case "list", "list-buttons":
		printButtons()
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func handleMove() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: please provide movement")
		fmt.Fprintln(os.Stderr, "Usage: mouse-cli move 100")
		fmt.Fprintln(os.Stderr, "       mouse-cli move 100,200")
		os.Exit(1)
	}

	dx, dy, err := parseMovement(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	vm, err := NewVirtualMouse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create virtual mouse: %v\n", err)
		fmt.Fprintln(os.Stderr, "Hint: Need root or input group permission")
		os.Exit(1)
	}
	defer vm.Close()

	fmt.Printf("Moving mouse: dx=%d, dy=%d\n", dx, dy)
	if err := vm.Move(dx, dy); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func handleScroll() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: please provide scroll amount")
		fmt.Fprintln(os.Stderr, "Usage: mouse-cli scroll 3 (up) or -3 (down)")
		os.Exit(1)
	}

	lines, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	vm, err := NewVirtualMouse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create virtual mouse: %v\n", err)
		fmt.Fprintln(os.Stderr, "Hint: Need root or input group permission")
		os.Exit(1)
	}
	defer vm.Close()

	direction := "up"
	if lines < 0 {
		direction = "down"
	}
	fmt.Printf("Scrolling: %d lines %s\n", lines, direction)
	if err := vm.Scroll(lines); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func handleClick() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: please provide button")
		fmt.Fprintln(os.Stderr, "Usage: mouse-cli click left")
		os.Exit(1)
	}

	button := os.Args[2]

	vm, err := NewVirtualMouse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create virtual mouse: %v\n", err)
		fmt.Fprintln(os.Stderr, "Hint: Need root or input group permission")
		os.Exit(1)
	}
	defer vm.Close()

	fmt.Printf("Clicking: %s\n", button)
	if err := vm.Click(button); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func handlePress() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: please provide button")
		fmt.Fprintln(os.Stderr, "Usage: mouse-cli press left")
		os.Exit(1)
	}

	button := os.Args[2]

	vm, err := NewVirtualMouse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create virtual mouse: %v\n", err)
		fmt.Fprintln(os.Stderr, "Hint: Need root or input group permission")
		os.Exit(1)
	}
	defer vm.Close()

	fmt.Printf("Pressing: %s\n", button)
	if err := vm.Press(button); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done! (button held)")
}

func handleRelease() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: please provide button")
		fmt.Fprintln(os.Stderr, "Usage: mouse-cli release left")
		os.Exit(1)
	}

	button := os.Args[2]

	vm, err := NewVirtualMouse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create virtual mouse: %v\n", err)
		fmt.Fprintln(os.Stderr, "Hint: Need root or input group permission")
		os.Exit(1)
	}
	defer vm.Close()

	fmt.Printf("Releasing: %s\n", button)
	if err := vm.Release(button); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done! (button released)")
}

func printButtons() {
	fmt.Println("Supported mouse buttons:")
	fmt.Println("  left, l    - Left button")
	fmt.Println("  right, r   - Right button")
	fmt.Println("  middle, m  - Middle/Wheel button")
	fmt.Println("  side, s    - Side button")
	fmt.Println("  extra, e   - Extra button")
}

func printUsage() {
	fmt.Println("mouse-cli - Virtual mouse CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  mouse-cli move <dx>           Move mouse (dx or dx,dy)")
	fmt.Println("  mouse-cli scroll <lines>     Scroll wheel (positive=up, negative=down)")
	fmt.Println("  mouse-cli click <button>      Click a button")
	fmt.Println("  mouse-cli press <button>     Hold down a button")
	fmt.Println("  mouse-cli release <button>  Release a button")
	fmt.Println("  mouse-cli list-buttons      List supported buttons")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  mouse-cli move 100              Move right 100 pixels")
	fmt.Println("  mouse-cli move -100,-50         Move left 100, up 50")
	fmt.Println("  mouse-cli scroll 3              Scroll up 3 lines")
	fmt.Println("  mouse-cli click left            Left click")
	fmt.Println("  mouse-cli click right           Right click")
	fmt.Println("  mouse-cli press left            Hold left button")
	fmt.Println("  mouse-cli release left          Release left button")
	fmt.Println("")
	fmt.Println("Run 'mouse-cli list-buttons' for more information.")
}
