db.instruction_execution.updateMany({status: 0}, {$set: {status: 3}});
db.job_history.updateMany({status: 0}, {$set: {status: 3}});
