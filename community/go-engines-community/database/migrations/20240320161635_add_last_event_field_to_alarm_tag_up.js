db.alarm_tag.updateMany({type: 0}, [
    {
        $set: {
            last_event_date: "$updated"
        }
    }
]);
