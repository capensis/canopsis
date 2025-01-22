db.default_entities.updateMany({resolve_deleted_event_processed: {$ne: null}, soft_deleted: null}, {
    $unset: {resolve_deleted_event_processed: ""}
});

db.default_entities.updateMany({resolve_deleted_event_sent: {$ne: null}, soft_deleted: null}, {
    $unset: {resolve_deleted_event_sent: ""}
});
