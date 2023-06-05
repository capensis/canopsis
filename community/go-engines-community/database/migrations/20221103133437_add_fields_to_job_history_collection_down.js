db.job_history.updateMany({}, {
    $unset: {
        name: "",
        retry_amount: "",
        retry_interval: "",
        multiple_executions: "",
        type: "",
        host: "",
        auth_username: "",
        auth_token: "",
        index: "",
    }
});
