db.broadcast_message.updateMany(
    {
        priority: null
    },
    {
        $set: {
            priority: 1
        }
    }
);

db.broadcast_message.updateMany(
    {
        closable: null
    },
    {
        $set: {
            closable: true
        }
    }
);
