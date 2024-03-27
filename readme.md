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
