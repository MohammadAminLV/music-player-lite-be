# Simple Music Player app Back-end


## Quick start
Copy `data.json.example` to `data.json`
```bash
go mod tidy
```
```bash
go run .
```
The server starts on the PORT 8000 by default.


## Docker
Build and run:


```bash
docker-compose up --build
```


## API Endpoints
- `GET /api/tracks` - get all tracks
- `get /api/tracks/:index` - get a single track by its index (Because no ID exist in the data)

