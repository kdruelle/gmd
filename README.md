gmd — Go-based Docker Manager (TUI)

gmd is a fast, minimal, terminal-native Docker manager written in Go.
It provides a clean TUI experience for listing, inspecting, updating, and operating Docker containers and images — without needing Portainer or a browser.

This project uses:
	•	Bubble Tea (for TUI state machine)
	•	Lipgloss (for styling)
	•	Docker SDK for Go (container & image operations)
	•	Model Stack architecture (for navigation between screens)

⸻

📦 Features

Real-time Docker monitoring
	•	Initial snapshot of containers and images
	•	Live event listener for Docker (create, start, pull, die, destroy, etc.)
	•	Automatic UI refresh on changes

Images panel
	•	Displays images similarly to Portainer (grouped, sorted, tagged)
	•	Detects unused images
	•	Supports deletion with UI feedback
	•	Detailed rendering with Lipgloss styling

Containers panel
	•	Name, ShortID, status, and update availability flags
	•	Colored status indicators (running/exited/restarting/paused)
	•	Live refresh on events
	•	Trigger updates via keyboard (u)

Interactive container update workflow

Full update pipeline implemented in a dedicated model:
	1.	docker pull with per-layer progress bars
	2.	Stop container with spinner
	3.	Remove container
	4.	Recreate container from its previous inspect
	5.	Start container
	6.	Return to main UI when complete

Includes:
	•	bubbles/progress for per-layer bars
	•	Spinners for blocking steps
	•	Clean multi-line logs during update
	•	A spinUntilDone helper for long operations

Shell and logs from TUI
	•	Exit the alt-screen and open a real shell inside a container (exec)
	•	Tail logs directly using Bubble Tea subprocess integration

⸻

🚀 Installation

Install gmd with a single command:
curl -sSfL https://raw.githubusercontent.com/kdruelle/gmd/main/install.sh | sh

The installer automatically:
	•	detects your OS and CPU architecture
	•	fetches the latest GitHub release
	•	verifies the SHA256 checksum
	•	extracts the correct binary
	•	installs it into /usr/local/bin

⸻

🧪 Roadmap
	•	Popup confirmation boxes
	•	Configurable themes
	•	Log viewer with formatting
	•	Column sorting (CPU, MEM, Name)
	•	Podman support (maybe)
	•	Remote Docker host support
	•	Plugin system

⸻

🤝 Contributing

Contributions are welcome:
	1.	Fork the repo
	2.	Create a feature branch
	3.	Add tests when possible
	4.	Submit a PR with a clear description

⸻

📄 License

MIT License.
You’re free to use it, fork it, extend it, or integrate it.