(function () {
    // Uncomment current Europe/Paris timezone or set your own.
    // Summer time
    var tz = "+02:00"
    // Winter time
    // var tz = "+01:00"

    var now = new Date();
    var month = now.getMonth() + 1;
    if (month < 10) {
        month = "0" + month;
    }
    var firstDateOfMonth = new Date(now.getFullYear() + "-" + month + "-01T00:00:00.000" + tz);
    firstDateOfMonth = Math.floor(firstDateOfMonth / 1000);

    var cursor = db.instruction_execution.aggregate([
        {$match: {status: 2}},
        {
            $group: {
                _id: "$instruction",
                users: {$addToSet: "$user"},
                month_executions: {
                    $sum: {
                        $cond: {
                            if: {$gte: ["$completed_at", firstDateOfMonth]},
                            then: 1,
                            else: 0,
                        }
                    }
                }
            }
        }
    ]);

    while (cursor.hasNext()) {
        var doc = cursor.next();
        db.instruction.updateOne({_id: doc._id}, {
            $set: {
                users: doc.users,
                month_executions: {
                    month: firstDateOfMonth,
                    count: doc.month_executions,
                }
            }
        });
    }
})();
