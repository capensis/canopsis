if (!db.permission.findOne({_id: "api_launch_event_recording"})) {
    db.permission.insertOne({
        _id: "api_launch_event_recording",
        name: "api_launch_event_recording",
        description: "Launch event recording",
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_launch_event_recording": 1
        }
    });
}

if (!db.permission.findOne({_id: "api_resend_events"})) {
    db.permission.insertOne({
        _id: "api_resend_events",
        name: "api_resend_events",
        description: "Event recorder resend events",
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_resend_events": 1
        }
    });
}

if (!db.permission.findOne({_id: "models_eventsRecord"})) {
    db.permission.insertOne({
        _id: "models_eventsRecord",
        name: "models_eventsRecord",
        description: "Events record",
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_eventsRecord": 1
        }
    });
}
