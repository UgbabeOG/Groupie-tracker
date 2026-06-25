# Groupie Tracker 

A robust, high-performance Go web application that consumes a multi-part entertainment REST API, processes relational data, and provides a highly optimized, user-friendly visualization interface. 

Built with an offline-first architectural mindset, this system guarantees high availability and zero runtime crashes by maintaining a thread-safe, hot-swappable in-memory cache synchronized with local JSON backups.

## Architecture & Features

- **Standard Library Only**: Built completely from scratch using only Go's native standard packages (`net/http`, `html/template`, `sync`, etc.).
- **In-Memory Cache (Boot Loader)**: The application loads all relational API entities into memory RAM at startup. Handlers read directly from RAM, eliminating runtime API network latency and disk I/O bottlenecks.
- **Resilient Offline Cache**: If the live API goes down, the server seamlessly falls back to reading local `json` backup storage files.
- **Thread-Safe Hot Swapping**: Includes an automated background synchronization mechanism that updates disk backups and safely swaps memory slices using a `sync.RWMutex` to prevent data-race crashes.
- **Safe HTML Template Buffer Rendering**: Executes template generation to an internal binary buffer before transmitting headers, ensuring zero server panics or half-rendered pages if structural faults occur.


### Architectural Diagram

```
                                +-----------------------------------+
                                |     Groupie Tracker Go App        |
                                +-----------------------------------+
                                                  |
                                       On Boot    v
                                +-----------------------------------+
                                |        init() / bootstrap         |
                                +-----------------------------------+
                                                  |
                         +------------------------+------------------------+
                         | Does Local Cache Exist?                         |
                         +------------------------+------------------------+
                                    |                                 |
                           No       v                                 v       Yes
               +--------------------------+               +--------------------------+
               |  Fetch Live REST API     |               |   Read Backup Files      |
               |  (artists/locations etc) |               |  (artists.json, etc.)    |
               +--------------------------+               +--------------------------+
                            |                                         |
            Save to Disk    v                                         |
               +--------------------------+                           |
               | Local JSON Cache Storage |                           |
               +--------------------------+                           |
                            |                                         |
                            +--------------------+--------------------+
                                                 |
                                                 v
                                  +------------------------------+
                                  |    json.NewDecoder().Decode()|
                                  +------------------------------+
                                                 |
                                                 v
                                  +------------------------------+
                                  |     Global Memory Cache      |
                                  |  (Protected by sync.RWMutex) |
                                  +------------------------------+
                                                 |
                        +------------------------+------------------------+
                        |                                                 |
         User HTTP Request Traffic                               Background Update Ticker
                        v                                                 v
         +------------------------------+                  +------------------------------+
         |      RWMutex.RLock()         |                  |       RWMutex.Lock()         |
         | (Concurrent Safe Read Access)|                  | (Exclusive Write/Memory Swap)|
         +------------------------------+                  +------------------------------+
                        |                                                 |
                        v                                                 v
         +------------------------------+                  +------------------------------+
         | Serve Handlers / UI Templates|                  | Overwrite RAM Cache Variables|
         +------------------------------+                  +------------------------------+

```

---

## Getting Started

### Prerequisites
- Go `1.26` or higher
- Docker

### Running Locally
```bash
# Clone the repository
git clone [https://acad.learn2earn.ng/git/gugbabe/groupie-tracker.git](https://acad.learn2earn.ng/git/gugbabe/groupie-tracker.git)
cd groupie-tracker

# Run the Go server
go run main.go

```

The server will boot up and bind to `http://localhost:8080`.

### Running with Docker

```bash
# Build the optimized production container
sudo docker build -t groupie-tracker .

# Start the background instance
sudo docker run -d -p 8080:8080 --name groupie-app groupie-tracker

# Monitor runtime server logs
sudo docker logs -f groupie-app

```

---

## 👥 Team & Task Division

* **Godwin Ugbabe (@gugbabe)**: Project Captain / Lead Integration & Submission
* **Osazuwa Matthew Ogbebor (@oogbebor)**: Core Infrastructure, Engine Core & Concurrency
* **Cynthia Isa Zoaka (@cisazoak)**: Front-End UI Layouts & Data Visualization Graphics

```

---

## Project Folder Structure

To support your specific views (`index`, `artistDetail`, and `error`) alongside our optimized caching backend engine, structure your repository cleanly like this:

```text
groupie-tracker/
├── .dockerignore
├── Dockerfile
├── go.mod
├── main.go                       # Web server entry point, router initialization
├── data/                         # Local JSON file cache persistence layer
│   ├── artists.json
│   ├── dates.json
│   ├── locations.json
│   └── relation.json
├── static/                       # Client assets handled by http.FileServer
│   ├── css/
│   │   └── style.css                 
│   └── js/
│       └── main.js
├── templates/                    # Core presentation layouts
│   ├── index.html                # Cards/Grid dashboard view for all bands
│   ├── artistDetail.html         # Rich profile view containing tables/maps
│   └── error.html                # Generic fail-safe message renderer (404/500)
└── utils/                        # Encapsulated application modules
    ├── cache.go                  # Global state management, loops, and Mutex swaps
    ├── handlers.go               # HTTP routing handlers (Home, Detail, Static)
    ├── models.go                 # Struct definitions with matched structural JSON tags
    └── template.go               # Crash-proof buffer rendering implementation

```

---

### Architectural Alignment Check

* **`utils/models.go`**: Contains the nested structs with metadata mapping rules for the 4 distinct entities.
* **`utils/cache.go`**: Declares your `sync.RWMutex`, reads files from the `data/` folder on boot, and runs a background routine to query fresh API data and write it back to disk.
* **`utils/handlers.go`**: Coordinates the routes and passes memory slices to `templates/index.html`, queries by ID parameters for `templates/artistDetail.html`, and redirects status failures to `templates/error.html`.
