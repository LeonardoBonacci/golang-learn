// Switch to the database
db = db.getSiblingDB('predictions_db');

// Create the collection if it doesn't exist
db.createCollection("predictions");

// Optional: Create an index for faster lookups (example on fooId)

// Create a TTL index (e.g., expire records after 60 seconds)
db.predictions.createIndex({ timestamp: 1 }, { expireAfterSeconds: 60 });

db.predictions.createIndex({ fooId: 1 });

print("✅ Initialized MongoDB with 'predictions' collection and TTL index");