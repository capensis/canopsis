db.permission.deleteOne({_id: "listalarm_removeAssociatedTicket"});
db.role.updateMany({}, {
    $unset: {
        "permissions.listalarm_removeAssociatedTicket": "",
    }
});
db.role_template.updateMany({}, {
    $unset: {
        "permissions.listalarm_removeAssociatedTicket": "",
    }
});

db.periodical_alarm.updateMany(
    {"v.failed_ticket": {$ne: null}},
    {$unset: {"v.failed_ticket": null}},
);
