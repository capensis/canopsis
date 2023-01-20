db.instruction.updateMany({type: 1}, {$unset: {triggers: ""}});
db.instruction_execution.updateMany({type: 1}, {$unset: {trigger: ""}});
