db.instruction.updateMany({type: 1}, {$set: {triggers: ["create"]}});
db.instruction_execution.updateMany({type: 1}, {$set: {trigger: "create"}});
