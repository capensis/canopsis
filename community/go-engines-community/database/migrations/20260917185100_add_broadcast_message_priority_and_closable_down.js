db.broadcast_message.updateMany(
    {},
    {
        $unset: {
            priority: "",
            closable: ""
        }
    }
);
