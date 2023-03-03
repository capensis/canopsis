db.default_rights.updateMany({}, {
    $rename: {
        email: "mail",
        ui_groups_navigation_type: "groupsNavigationType",
        ui_tours: "tours",
    }
});
db.default_rights.updateMany({crecord_type: "action"}, {
    $rename: {
        description: "desc",
    }
});
