db.alarm_tag.find({type: 1}).forEach(function (doc) {
    if (!doc.value || !doc.color) {
        return;
    }

    const label = doc.value.split(":")[0];
    const color = doc.color;
    db.alarm_tag_color.updateOne({_id: label}, {$setOnInsert: {color: color}}, {upsert: true});
});
