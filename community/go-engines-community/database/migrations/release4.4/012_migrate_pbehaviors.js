(function () {
    // Unset null tstop
    db.pbehavior.updateMany({
        "tstop": { "$type": 10 }, // BsonNull type
    }, {
        "$unset": { "tstop": "" },
    });
})();