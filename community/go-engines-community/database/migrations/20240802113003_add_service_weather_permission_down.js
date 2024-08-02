db.permission.deleteMany({
    _id: {
        $in: [
            "serviceweather_entityDeclareTicket",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.serviceweather_entityDeclareTicket": "",
    }
});
