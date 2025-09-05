db.external_data_table.updateMany(
    {column_configs: null},
    [
        {
            $set: {
                column_configs: {
                    $map: {
                        input: { $zip: { inputs: ["$columns", "$column_types"] } },
                        as: "pair",
                        in: {
                            name: { $arrayElemAt: ["$$pair", 0] },
                            type: 1,
                            tag: { $arrayElemAt: ["$$pair", 1] }
                        }
                    }
                }
            }
        },
        {
            $unset: ["columns", "column_types", "column_lengths"]
        }
    ]
);

db.external_data_import_worker.deleteMany({})
