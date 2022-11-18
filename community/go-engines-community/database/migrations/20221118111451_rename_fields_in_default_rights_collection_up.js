db.default_rights.updateMany({}, {
    $rename: {
        mail: "email",
        desc: "description",
        groupsNavigationType: "ui_groups_navigation_type",
        tours: "ui_tours",
    }
});
