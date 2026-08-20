db.default_entities.createIndex({ "$**": "text" }, {name:"full_text", default_language:"none"});
if (!db.getCollectionNames().includes('pattern_optimize_job')) {
    db.createCollection("pattern_optimize_job", {capped: true, size: 104857600});
}
