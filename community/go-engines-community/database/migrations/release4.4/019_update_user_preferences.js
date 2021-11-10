(function () {
    db.userpreferences.find().forEach(function (doc) {
        db.userpreferences.updateOne({_id: doc._id}, {
            $set: {
                user: doc.name,
                widget: doc.widget_id,
                content: doc.widget_preferences,
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
