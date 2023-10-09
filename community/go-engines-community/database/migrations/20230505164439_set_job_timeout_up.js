db.job.find({}).forEach(function (doc) {
    var retryAmount = doc.retry_amount
    var retryInterval = doc.retry_interval

    // default value
    var jobWaitInterval = {
        "unit": "s",
        "value": 60
    }

    if (retryInterval !== undefined) {
        if (retryAmount !== undefined) {
            retryInterval.value = retryInterval.value * retryAmount
        }

        if (retryInterval.unit === "s" && retryInterval.value < 60) {
            retryInterval.value = 60
        }

        jobWaitInterval = retryInterval
    }

    db.job.updateOne(
        {
            _id: doc._id
        },
        {
            $set: {job_wait_interval: jobWaitInterval},
            $unset: {retry_amount: "", retry_interval: ""}
        },
    );
})
