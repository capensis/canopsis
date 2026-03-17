// Available global functions:
// genID returns a new string UUID
// isInt checks if a value is integer
// toInt transforms value to integer

if (db.widget_templates.countDocuments({title: "Default alarm sort columns", type: "alarm_sort_columns"}) === 0) {
    const now = Math.ceil((new Date()).getTime() / 1000);
    db.widget_templates.insertOne({
        _id: genID(),
        title: "Default alarm sort columns",
        type: "alarm_sort_columns",
        sort_columns: [{sort_by: "v.creation_date", sort: "desc"}],
        author: "root",
        created: now,
        updated: now
    });
}
