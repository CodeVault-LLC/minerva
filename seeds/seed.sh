while IFS= read -r website; do
  curl -X POST http://localhost:3000/api/v1/scans -d "{\"url\":\"$website\"}" -H "Content-Type: application/json"
done < /Users/lukasolsen/Documents/codevault/humblebrag-api/seeds/websites-1000.txt
