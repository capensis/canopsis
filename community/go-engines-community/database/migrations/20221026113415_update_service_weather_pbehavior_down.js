db.pbehavior.updateMany({}, {
    $unset: {
        origin: "",
        entity: ""
    }
});
