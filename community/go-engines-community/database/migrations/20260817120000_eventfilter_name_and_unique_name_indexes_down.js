// Renamed duplicate names cannot be restored.
db.eventfilter.getIndexes().forEach(function (index) {
    if (index.name === "name_1") {
        db.eventfilter.dropIndex("name_1");
    }
});
db.dynamic_infos.getIndexes().forEach(function (index) {
    if (index.name === "name_1") {
        db.dynamic_infos.dropIndex("name_1");
    }
});
db.meta_alarm_rules.getIndexes().forEach(function (index) {
    if (index.name === "name_1") {
        db.meta_alarm_rules.dropIndex("name_1");
    }
});

db.eventfilter.updateMany({}, {$unset: {name: ""}});
