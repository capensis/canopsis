(function () {
    db.periodical_alarm.find({"v.resolved": {"$exists": true}}).forEach(function (doc) {
        db.resolved_alarms.insertOne(doc)
    });
})();
