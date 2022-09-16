(function () {
    db.instruction_mod_stats.updateMany(
        {successful: null, execution_count: {$ne: null}}, 
        [{$set: {successful: "$execution_count"}}]);
    db.instruction_week_stats.updateMany(
        {successful: null, execution_count: {$ne: null}}, 
        [{$set: {successful: "$execution_count"}}])
})();
