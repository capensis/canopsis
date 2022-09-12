db.dynamic_infos.updateMany({}, {
    $rename: {
        created: "creation_date",
        updated: "last_modified_date",
    }
});
