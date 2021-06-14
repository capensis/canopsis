(function () {
    var now = Math.ceil((new Date()).getTime() / 1000);

    db.view_playlist.find().forEach(function (doc) {
        var val = doc.interval.interval;
        var unit = doc.interval.unit;
        var seconds = val;
        switch (unit) {
            case "s":
                seconds = val;
                break
            case "m":
                seconds = val * 60;
                break
            case "h":
                seconds = val * 60 * 60;
                break
        }

        db.view_playlist.updateOne({_id: doc._id}, {
            $set: {
                interval: {
                    seconds: seconds,
                    unit: unit
                },
                created: now,
                updated: now,
            },
        })
    });
})();
