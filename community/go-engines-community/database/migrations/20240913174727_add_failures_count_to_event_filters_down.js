db.eventfilter.updateMany({}, {
    $unset: {
        failures_count: "",
        unread_failures_count: "",
    }
});
