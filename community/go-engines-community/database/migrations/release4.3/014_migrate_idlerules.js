(function () {
    db.idle_rule.updateMany({type: "last_event"}, {
        $set: {
            type: "alarm",
            alarm_condition: "last_event",
        }
    });
    db.idle_rule.updateMany({type: "last_update"}, {
        $set: {
            type: "alarm",
            alarm_condition: "last_update",
        }
    });
    var now = Math.ceil((new Date()).getTime() / 1000);
    db.idle_rule.updateMany({created: {$exists: false}}, {
        $set: {
            enabled: true,
            created: now,
            updated: now,
        }
    });

    var priority = -1;
    db.idle_rule.find().forEach(function (doc) {
        priority++;
        var val = parseInt(doc.duration);
        var unit = doc.duration.replace(val.toString(), "");
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
        db.idle_rule.updateOne({_id: doc._id}, {
            $set: {
                name: "Rule #" + (priority + 1),
                priority: priority,
                duration: {
                    seconds: seconds,
                    unit: unit,
                }
            },
        });
    });
})();
