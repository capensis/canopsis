db.user.updateMany({}, [
    {
        $set: {
            password: {$toLower: "$password"}
        }
    }
]);
