db.instruction_execution.updateMany({}, {
    $unset: {
        name: "",
        description: "",
        type: "",
        timeout_after_execution: "",
        created_at: "",
        priority: "",
    }
});
