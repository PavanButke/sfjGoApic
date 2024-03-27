# How to run app
- `go run .` (at root folder)


## DB Conection
- postgres: assign db
- INSERT Query:
  - `CREATE TABLE job (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    duration INTERVAL,
    status VARCHAR(50)
    );`

## APIs
- http://localhost:8099/jobs [POST]
- Body :
    - `{
    "name": "Job Pavan",
    "duration": "5s",
    "status": "Pending"
    }`

- http://localhost:8099/jobs [GET]
    - `[
    {
        "id": 21,
        "name": "Job Pavan",
        "duration": "00:00:05",
        "status": "Pending"
    },
    {
        "id": 22,
        "name": "Job Pavan",
        "duration": "00:00:05",
        "status": "Pending"
    },
    {
        "id": 23,
        "name": "Job Pavan",
        "duration": "00:00:05",
        "status": "Pending"
    },
    {
        "id": 24,
        "name": "Job Pavan",
        "duration": "00:00:05",
        "status": "Pending"
    }]`

# Websocket 
- ws://localhost:8098/ws

# Project Structure
-project-root/
│
│   └── config.yaml            # Configuration files (e.g., database config)
│
├── handlers/                 
│   ├── handler.go            # Define Add Job and Update Job Logic
│   └── job_interfacce.go     # Declaring funcs
│
├── models/                   
│   └── db.go                # Define data models (e.g., Job struct)
│
├── server/                   
│   ├── router.go              # Define routes 
│   └── server.go              # Start HTTP server
    └── middleware.go          # middleware       
│
├── app/                 
│   └── sfj.go                  # Shortest-First Logic Impl
│
├── websocket/                    
│   └── ws.go                   # Websocket functions
│
│
├── main.go                   # Entry point of the application
│
└── README.md                 # Documentation about the project

# What Can Be Improved
- A seperate queries class can be developed
- Gin can also be used
- Config.yaml may be hidden
  
