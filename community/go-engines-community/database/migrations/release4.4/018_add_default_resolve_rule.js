(function () {
    var now = Math.ceil((new Date()).getTime() / 1000);
    db.resolve_rule.insertOne({
        "_id": "default_rule",
        "loader_id": "default_rule",
        "author": "root",
        "name": "Default rule",
        "duration": {
            "seconds": 60,
            "unit": "m"
        },
        "alarm_patterns": [],
        "entity_patterns": [],
        "created": now,
        "updated": now,
        "priority": 1
    })
})();
