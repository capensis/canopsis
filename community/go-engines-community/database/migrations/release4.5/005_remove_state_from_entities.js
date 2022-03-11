db.default_entities.updateMany(
    {},
    {
        $unset: {"state": 1}
    }
);
