if (!db.permission.findOne({_id: "serviceweather_entityDeclareTicket"})) {
    db.permission.insertOne({
        _id: "serviceweather_entityDeclareTicket",
        name: "serviceweather_entityDeclareTicket",
        description: "Service weather: Access to 'Declare Ticket' action"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.serviceweather_entityDeclareTicket": 1
        }
    });
}
