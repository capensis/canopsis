db.default_entities.updateMany({}, {$unset: {alarms_cumulative_data: ""}});

db.default_entities.find().forEach(function (doc) {
    if (!doc.output_template) {
        return;
    }

    var output = doc.output_template;
    output = output.replace(/\{\{([^{}]*)\.Alarms([^{}]*)}}/, "{{$1.Active$2}}")
    output = output.replace(/\{\{([^{}]*)\.State.Info([^{}]*)}}/, "{{$1.State.Ok$2}}")
    db.default_entities.updateOne({_id: doc._id}, {$set: {output_template: output}});
});

db.widget_filters.aggregate([
    {$match: {weather_service_pattern: {$ne: null}}},
    {
        $lookup: {
            from: "widgets",
            localField: "widget",
            foreignField: "_id",
            as: "w",
        }
    },
    {$unwind: "$w"},
    {$match: {"w.type": "ServiceWeather"}},
]).forEach(function (doc) {
    if (!doc.weather_service_pattern) {
        return;
    }

    var iconMap = {
        ok: "wb_sunny",
        minor: "person",
        major: "person",
        critical: "wb_cloudy",
        maintenance: "build",
        inactive: "brightness_3",
        pause: "pause",
    };
    var newPattern = [];
    var updated = false;
    doc.weather_service_pattern.forEach(function (group) {
        var newGroup = [];
        group.forEach(function (cond) {
            if (cond.field === "icon" || cond.field === "secondary_icon") {
                if (typeof cond.cond.value === "string") {
                    cond.cond.value = iconMap[cond.cond.value];
                    updated = true;
                } else if (Array.isArray(cond.cond.value)) {
                    for (var i = 0; i < cond.cond.value.length; i++) {
                        if (typeof cond.cond.value[i] === "string") {
                            cond.cond.value[i] = iconMap[cond.cond.value[i]];
                            updated = true;
                        }
                    }
                }
            }

            newGroup.push(cond);
        });

        newPattern.push(newGroup);
    });

    if (updated) {
        db.widget_filters.updateOne({_id: doc._id}, {$set: {weather_service_pattern: newPattern}});
    }
});
