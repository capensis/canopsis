db.widgets.find({"parameters.inlineLinksCount": {$gt: 0}}).forEach(function (doc) {
    if (!doc.parameters || !doc.parameters.inlineLinksCount || !doc.parameters.widgetColumns) {
        return;
    }

    var set = {};
    doc.parameters.widgetColumns.forEach(function (col, colIndex) {
        if (col.value === "links") {
            set["parameters.widgetColumns." + colIndex + ".inlineLinksCount"] = doc.parameters.inlineLinksCount;
        }
    });

    var update = {
        $unset: {
            "parameters.inlineLinksCount": "",
        },
    };
    if (Object.keys(set).length > 0) {
        update["$set"] = set;
    }
    db.widgets.updateOne({_id: doc._id}, update);
});
