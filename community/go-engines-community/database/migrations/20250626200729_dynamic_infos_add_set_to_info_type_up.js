db.dynamic_infos.updateMany(
    {},
    {
        $set: {
            "infos.$[tplValue].type": "set_to_info_from_template",
            "infos.$[constValue].type": "set_to_info"
        }
    },
    {
        arrayFilters: [
            {"tplValue.value": {$type: "string", $regex: /\{\{/}},
            {$nor: [{"constValue.value": {$type: "string", $regex: /\{\{/ }}]}
        ]
    }
)
