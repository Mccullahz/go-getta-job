db = db.getSiblingDB('job_search_db');

// collections
db.createCollection('users');
db.createCollection('geo_results');
db.createCollection('businesses');
db.createCollection('jobs');
db.createCollection('results');
db.createCollection('starred_jobs');
db.createCollection('applied_jobs');

// indexes
db.users.createIndex({ email: 1 }, { unique: true });
db.geo_results.createIndex({ zip: 1, radius: 1 });
db.businesses.createIndex({ name: 1 });
db.jobs.createIndex({ title: "text" });
db.starred_jobs.createIndex({ user_id: 1, job_id: 1 }, { unique: true });
db.applied_jobs.createIndex({ user_id: 1, job_id: 1 }, { unique: true });

// load initial data
try {
  const geoData = cat('/docker-entrypoint-initdb.d/geo_results.json');
  if (geoData) {
    db.geo_results.insertMany(JSON.parse(geoData));
  }

  const resultsData = cat('/docker-entrypoint-initdb.d/results.json');
  if (resultsData) {
    db.results.insertMany(JSON.parse(resultsData));
  }

  print("Seed data loaded successfully");
} catch (e) {
  print("No seed data found or failed to load:", e);
}
