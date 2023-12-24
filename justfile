default:
    @just --choose

mongo-uri := "mongodb://mongoadmin:secret@localhost:27017/mydatabase?authSource=admin"
mongo-collection := "products"
mongo-import-file := "seed.json"
mongo-export-file := "products.out.json"

seed-mongo:
    @mongoimport \
        --uri {{ mongo-uri }} \
        --collection {{ mongo-collection }} \
        --file {{ mongo-import-file }} \
        --jsonArray

# First 5 days based on seed.json
start := "2021-12-01T00:00:00Z"
end := "2023-12-06T00:00:00Z"

export-mongo :
	@echo "{ \"updatedAt\": { \"\$gte\": { \"\$date\": \"{{start}}\" }, \"\$lte\": { \"\$date\": \"{{end}}\" } } }"
	@mongoexport \
		--uri {{ mongo-uri }} \
		--collection {{ mongo-collection }} \
		--out {{ mongo-export-file }} \
		--query "{ \"updatedAt\": { \"\$gte\": { \"\$date\": \"{{start}}\" }, \"\$lte\": { \"\$date\": \"{{end}}\" } } }"
