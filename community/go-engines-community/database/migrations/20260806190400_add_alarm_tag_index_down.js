db.alarm_tag.getIndexes().forEach(function (index) {
    if (index.name === "label_1") {
        db.alarm_tag.dropIndex(index.name);
    }
});