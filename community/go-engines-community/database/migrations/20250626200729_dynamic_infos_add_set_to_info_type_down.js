db.dynamic_infos.updateMany({}, {$unset: {"infos.$[].type": ""}})
