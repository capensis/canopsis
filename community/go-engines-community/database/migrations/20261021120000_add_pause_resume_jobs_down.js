if (db.getCollectionNames().includes("pause_job")) {
    db.pause_job.getIndexes().forEach(function (index) {
        if (index.name === "retry_at_1" || index.name === "source_timestamp_path_1") {
            db.pause_job.dropIndex(index.name);
        }
    });

    if (db.pause_job.countDocuments({}) === 0) {
        db.pause_job.drop();
    }
}

if (db.getCollectionNames().includes("resume_job")) {
    db.resume_job.getIndexes().forEach(function (index) {
        if (index.name === "retry_at_1") {
            db.resume_job.dropIndex(index.name);
        }
    });
}

if (db.getCollectionNames().includes("resume_job") && db.resume_job.countDocuments({}) === 0) {
    db.resume_job.drop();
}

db.default_entities.getIndexes().forEach(function (index) {
    if (index.name === "pause_source_timestamp_path_1") {
        db.default_entities.dropIndex(index.name);
    }
});
