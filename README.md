# Gofeed

A simple Go scheduling app that reads a `config.yml` file and schedules jobs based on the configuration.

## Configuration

```yaml
export:
  mongoUri: "mongodb://mongoadmin:secret@localhost:27017/mydatabase?authSource=admin"
  mongoCollection: "products"
  
import:
    postgresUri: "postgresql://postgres:secret@localhost:5432/products"
    postgresTable: "public.products_raw"
    
job:
    name: "test"

    # Initial start time
    start: "2023-12-01T00:00:00Z"
    
    # Schedule value
    schedule: 1
    # Schedule time unit - valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
    scheduleUnit: "m"
```

## Usage

### Start Mongo and Postgres

```bash
docker-compose up -d
```

### (Optionally) Seed Mongo with data

This `just` recipe will seed Mongo with data from `data/seed.json`.
```bash
just seed-mongo
``` 

or 

```bash
mongoimport \
    --uri MONGO_URI \
    --collection MONGO_COLLECTION \
    --file IMPORT_JSON_FILE \
    --jsonArray
``` 

### Build and run the app

```bash
go build -o gofeed ./src
./gofeed -config=config.yml
```
