db.userpreferences.updateMany({"content.dense": 1}, {
    $set: {
        "content.dense": true,
    },
});
db.userpreferences.updateMany({"content.dense": 0}, {
    $set: {
        "content.dense": false,
    },
});
db.widgets.updateMany({"parameters.dense": 1}, {
    $set: {
        "parameters.dense": true,
    },
});
db.widgets.updateMany({"parameters.dense": 0}, {
    $set: {
        "parameters.dense": false,
    },
});
