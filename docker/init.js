db = db.getSiblingDB('job_search_db');

// collections schema and create statements
//===== users collection =====
db.createCollection('users', {
	validator: {
		$jsonSchema: {
			bsonType: 'object',
			required: ['username', 'email', 'password_hash'],
			properties: {
				username: { bsonType: 'string'},
				email: { bsonType: 'string'},
				password_hash: { bsonType: 'string'},
				created_at: { bsonType: 'date'}
			}
		}
	}
});
db.users.createIndex({ email: 1 }, { unique: true });

//===== geo results collection =====
db.createCollection('geo_results', {
	validator: {
		$jsonSchema: {
			bsonType: 'object',
			required: ['user_id', 'zip', 'radius'],
			properties: {
				user_id: { bsonType: 'objectId' },
				zip: { bsonType: 'string' },
				radius: { bsonType: 'int' },
				created_at: { bsonType: 'date' }
			}
		}
	}
});
db.geo_results.createIndex({ user_id: 1 })
db.geo_results.createIndex({ zip: 1, radius: 1 });

//===== business collection =====
db.createCollection('businesses', {
	validator: {
		$jsonSchema: {
			bsonType: 'object',
			required: ['geo_result_id', 'name', 'address', 'url', 'lat', 'lon'],
			properties: {
				geo_result_id: { bsonType: 'objectId' },
				name: { bsonType: 'string' },
				address: { bsonType: 'string' },
				url: { bsonType: 'string' },
				lat: { bsonType: 'double' },
				lon: { bsonType: 'double' },
			}
		}
	}
});
db.businesses.createIndex({ geo_result_id: 1 });
db.businesses.createIndex({ name: 1 });

//===== jobs collection =====
db.createCollection('jobs', {
	validator: {
		$jsonSchema: {
			bsonType: 'object',
			required: ['business_id', 'title', 'url'],
			properties: {
				business_id: { bsonType: 'objectId' },
				title: { bsonType: 'string' },
				description: { bsonType: 'string' },
				url: { bsonType: 'string' },
				posted_at: { bsonType: 'date' },
			}
		}
	}
});
db.jobs.createIndex({ business_id: 1 });
db.jobs.createIndex({ title: "text" });

//===== results collection =====
db.createCollection('job_results', {
	validator: {
		$jsonSchema: {
			bsonType: 'object',
			required: ['user_id', 'jobs', 'query_title'],
			properties: {
				user_id: { bsonType: 'objectId' },
				jobs: { bsonType: [array], items: { bsonType: 'objectId' } },
				query_title: { bsonType: 'string' },
				created_at: { bsonType: 'date' }
			}
		}
	}
});
db.job_results.createIndex({ user_id: 1 });
db.job_results.createIndex({ query_title: "text" });

//===== starred jobs collections =====
db.createCollection('starred_jobs', {
	validator: {
		$jsonSchema: {
			bsonType: 'object',
			required: ['user_id', 'job_id'],
			properties: {
				user_id: { bsonType: 'objectId' },
				job_id: { bsonType: 'objectId' },
				created_at: { bsonType: 'date' }
			}
		}
	}
});
db.starred_jobs.createIndex({ user_id: 1, job_id: 1 }, { unique: true });

//===== applied jobs collections =====
db.createCollection('applied_jobs', {
	validator: {
		$jsonSchema: {
			bsonType: 'object',
			required: ['user_id', 'job_id'],
			properties: {
				user_id: { bsonType: 'objectId' },
				job_id: { bsonType: 'objectId' },
				applied_at: { bsonType: 'date' }
			}
		}
	}
});


// debug print: no initial test data - system will use real scraped data
print("Database initialized - ready for real data");
