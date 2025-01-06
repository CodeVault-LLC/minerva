
while IFS= read -r website; do
  curl -Method POST -Uri http://127.0.0.1:3000/api/v1/scans -ContentType 'application/json' -Body '{\"url\":\"$website\"}' -UseBasicParsing
done < websites-1000.txt
