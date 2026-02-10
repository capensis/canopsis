db.webhook_check_ticket_status.createIndex(
    {ticket_id: 1, ticket_system_name: 1},
    {name: "ticket_id_ticket_system_name_1", unique: true}
)
