db.dynamic_infos.updateMany({}, {
    $rename: {
        creation_date: "created",
        last_modified_date: "updated",
    }
});
