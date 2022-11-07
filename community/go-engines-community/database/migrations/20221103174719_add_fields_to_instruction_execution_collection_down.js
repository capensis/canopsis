db.instruction_execution.updateMany({}, {
    $unset: {
        name: "",
        type: "",
        timeout_after_execution: "",
        created_at: "",
        priority: "",
    }
});
