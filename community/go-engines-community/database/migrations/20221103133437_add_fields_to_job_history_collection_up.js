db.job_history.aggregate([
    {
        $lookup: {
            from: "job",
            localField: "job",
            foreignField: "_id",
            as: "job"
        }
    },
    {$unwind: "$job"},
    {
        $lookup: {
            from: "job_config",
            localField: "job.config",
            foreignField: "_id",
            as: "config"
        }
    },
    {$unwind: "$config"},
    {
        $project: {
            name: "$job.name",
            retry_amount: "$job.retry_amount",
            retry_interval: "$job.retry_interval",
            multiple_executions: "$job.multiple_executions",
            type: "$config.type",
            host: "$config.host",
            auth_username: "$config.auth_username",
            auth_token: "$config.auth_token",
        }
    }
]).forEach(function (doc) {
    db.job_history.updateOne({_id: doc._id}, {
        $set: {
            name: doc.name,
            retry_amount: doc.retry_amount,
            retry_interval: doc.retry_interval,
            multiple_executions: doc.multiple_executions,
            type: doc.type,
            host: doc.host,
            auth_username: doc.auth_username,
            auth_token: doc.auth_token,
        }
    });
});
