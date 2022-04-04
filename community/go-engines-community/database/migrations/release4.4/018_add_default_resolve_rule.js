(function () {
    var id = "default_rule";
    var doc = db.resolve_rule.findOne({_id: id});
    if (!doc) {
        var now = Math.ceil((new Date()).getTime() / 1000);
        db.resolve_rule.insertOne({
            "_id": id,
            "loader_id": id,
            "author": "root",
            "name": "Default rule",
            "duration": {
                "seconds": 60,
                "unit": "m"
            },
            "created": now,
            "updated": now,
            "priority": 1
        })
    } else if (doc.alarm_patterns != null && doc.alarm_patterns.length == 0 &&
        doc.entity_patterns != null && doc.entity_patterns.length == 0) {
            // update previously inserted rule with empty alarm_paterns and entity_patterns
        db.resolve_rule.updateOne({_id: id}, {$unset: {"alarm_patterns": "", "entity_patterns": ""}});
    }
})();
