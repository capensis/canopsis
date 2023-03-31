db.userpreferences.updateMany({"content.dense": true}, {
    $set: {
        "content.dense": 1,
    },
});
db.userpreferences.updateMany({"content.dense": false}, {
    $set: {
        "content.dense": 0,
    },
});
db.widgets.updateMany({"parameters.dense": true}, {
    $set: {
        "parameters.dense": 1,
    },
});
db.widgets.updateMany({"parameters.dense": false}, {
    $set: {
        "parameters.dense": 0,
    },
});
