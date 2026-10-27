if (!db.getCollectionNames().includes("resume_job")) {
    db.createCollection("resume_job");
}

var resumeJobIndexes = {};
db.resume_job.getIndexes().forEach(function (index) {
    resumeJobIndexes[index.name] = true;
});

if (!resumeJobIndexes["retry_at_1"]) {
    db.resume_job.createIndex({retry_at: 1}, {name: "retry_at_1"});
}

var entityIndexes = {};
db.default_entities.getIndexes().forEach(function (index) {
    entityIndexes[index.name] = true;
});

if (!entityIndexes["pause_source_timestamp_path_1"]) {
    db.default_entities.createIndex(
        {"pause.source": 1, "pause.timestamp": 1, "pause.path": 1},
        {name: "pause_source_timestamp_path_1", partialFilterExpression: {pause: {$exists: true}}}
    );
}

if (!db.getCollectionNames().includes("pause_job")) {
    db.createCollection("pause_job");
}

var pauseJobIndexes = {};
db.pause_job.getIndexes().forEach(function (index) {
    pauseJobIndexes[index.name] = true;
});

if (!pauseJobIndexes["retry_at_1"]) {
    db.pause_job.createIndex({retry_at: 1}, {name: "retry_at_1"});
}

if (!pauseJobIndexes["source_timestamp_path_1"]) {
    db.pause_job.createIndex(
        {"event.pause_source": 1, "event.timestamp": 1, "event.pause_path": 1},
        {name: "source_timestamp_path_1"}
    );
}
