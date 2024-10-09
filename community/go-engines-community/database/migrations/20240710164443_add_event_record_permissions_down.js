db.permission.deleteMany({
    _id: {
        $in: [
            "api_launch_event_recording",
            "api_resend_events",
            "models_eventsRecord"
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_launch_event_recording": "",
        "permissions.api_resend_events": "",
        "permissions.models_eventsRecord": "",
    }
});
