db.alarm_tag.updateMany({type: 0}, {
    $unset: {
        last_event_date: ""
    }
});
