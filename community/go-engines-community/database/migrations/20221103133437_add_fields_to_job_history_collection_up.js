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
            execution: 1,
            next_exec: 1,
            name: "$job.name",
            retry_amount: "$job.retry_amount",
            retry_interval: "$job.retry_interval",
            multiple_executions: "$job.multiple_executions",
            type: "$config.type",
            host: "$config.host",
            auth_username: "$config.auth_username",
            auth_token: "$config.auth_token",
        }
    },
    {
        $group: {
            _id: "$execution",
            jobs: {$push: "$$ROOT"}
        }
    },
    {
        $lookup: {
            from: "instruction_execution",
            localField: "_id",
            foreignField: "_id",
            as: "execution"
        }
    },
    {$unwind: "$execution"},
    {
        $lookup: {
            from: "instruction",
            localField: "execution.instruction",
            foreignField: "_id",
            as: "instruction"
        }
    },
    {$unwind: "$instruction"},
    {
        $project: {
            execution: "$execution._id",
            instruction_type: "$instruction.type",
            jobs: 1
        }
    },
]).forEach(function (doc) {
    if (doc.instruction_type == 0) {
        doc.jobs.forEach(function (job) {
            var set = {
                name: job.name,
                retry_amount: job.retry_amount,
                retry_interval: job.retry_interval,
                multiple_executions: job.multiple_executions,
                type: job.type,
                host: job.host,
                auth_username: job.auth_username,
                auth_token: job.auth_token,
            };

            db.job_history.updateOne({_id: job._id}, {$set: set});
        });
        return;
    }

    var jobs = {};
    var hasRef = {};
    doc.jobs.forEach(function (job) {
        jobs[job._id] = job;

        if (job.next_exec) {
            hasRef[job.next_exec] = true;
        }
    });
    var jobId = null;
    doc.jobs.forEach(function (job) {
        if (!hasRef[job._id]) {
            jobId = job._id;
        }
    });
    var index = 0;
    while (true) {
        if (!jobId) {
            break;
        }
        var job = jobs[jobId];
        var set = {
            name: job.name,
            retry_amount: job.retry_amount,
            retry_interval: job.retry_interval,
            multiple_executions: job.multiple_executions,
            type: job.type,
            host: job.host,
            auth_username: job.auth_username,
            auth_token: job.auth_token,
            index: index,
        };

        db.job_history.updateOne({_id: job._id}, {$set: set});
        index++;
        jobId = job.next_exec;
    }
});
