```
docker exec -it mongo mongosh -u admin -p admin --authenticationDatabase admin  

use predictions_db
db.createCollection("predictions")

db.predictions.find().pretty()

```

```
curl -X POST "http://localhost:8081/fetch-prediction" -H "Content-Type: application/json" -d '{"foo_id": 42}'
```