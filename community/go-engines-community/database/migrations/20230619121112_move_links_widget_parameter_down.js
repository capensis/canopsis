db.widgets.find({"parameters.widgetColumns.value": "links"}).forEach(function (doc) {
    if (!doc.parameters || !doc.parameters.widgetColumns) {
        return;
    }

    var set = {};
    var unset = {};
    var inlineLinksCount = 0;
    doc.parameters.widgetColumns.forEach(function (col, colIndex) {
        if (col.value === "links" && col.inlineLinksCount > 0) {
            if (inlineLinksCount < col.inlineLinksCount) {
                inlineLinksCount = col.inlineLinksCount;
            }
            unset["parameters.widgetColumns." + colIndex + ".inlineLinksCount"] = "";
        }
    });

    if (inlineLinksCount > 0) {
        set["parameters.inlineLinksCount"] = inlineLinksCount;
    }

    var update = {};
    if (Object.keys(set).length > 0) {
        update["$set"] = set;
    }
    if (Object.keys(unset).length > 0) {
        update["$unset"] = unset;
    }
    if (Object.keys(update).length > 0) {
        db.widgets.updateOne({_id: doc._id}, update);
    }
});
