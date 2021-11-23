(function () {
    db.userpreferences.find().forEach(function (doc) {
        var userId = doc._id.replace(doc.widget_id + "_", "")

        db.userpreferences.updateOne({_id: doc._id}, {
            $set: {
                user: userId,
                widget: doc.widget_id,
                content: doc.widget_preferences,
                updated: doc.crecord_write_time,
            },
            $unset: {
                name: "",
                widget_id: "",
                widget_preferences: "",
                crecord_creation_time: "",
                crecord_name: "",
                crecord_type: "",
                crecord_write_time: "",
                enable: "",
                widgetXtype: "",
            }
        });
    });
})();
