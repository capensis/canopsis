db.eventfilter.updateMany({}, {
    $unset: {
        next_resolved_start: "",
        next_resolved_stop: "",
    }
});
