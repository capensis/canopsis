db.dynamic_infos.updateMany({}, {$set: {"infos.$[].type": "set_to_info"}})
