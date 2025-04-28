# API Changelog 24.10.0 vs. 25.04.0

## GET /account/me
- :warning: removed the optional property 'permissions/items/description' from the response with the '200' status


## PUT /account/me
- :warning: removed the success response with the status '201'


## POST /alarm-details
- :warning: removed the optional property 'data/children/data/items/v/infos_rule_version' from the response with the '207' status


## GET /alarm-tags
- :warning: removed the success response with the status '201'


## POST /alarm-tags
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /alarms
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status


## GET /alarms/{id}
- :warning: removed the optional property 'v/infos_rule_version' from the response with the '200' status


## PUT /alarms/{id}/ack
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /alarms/{id}/ackremove
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /alarms/{id}/assocticket
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /alarms/{id}/cancel
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /alarms/{id}/changestate
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /alarms/{id}/comment
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /alarms/{id}/snooze
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /alarms/{id}/uncancel
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /app-info
- :warning: removed the optional property 'stack' from the response with the '200' status


## POST /associativetable
- :warning: removed the optional property 'error' from the response with the '400' status


## POST /broadcast-message
- :warning: removed the optional property 'error' from the response with the '400' status


## POST /bulk/eventfilters
- :warning: the '/items/external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the '/items/items/allOf[subschema #2]/item/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '207'


## PUT /bulk/eventfilters
- :warning: the '/items/external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the '/items/items/allOf[subschema #2]/item/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '207'


## GET /cat/event-records/{id}
- :warning: the '/allOf[subschema #2]/data/items/event' response's property type/format changed from 'array'/'' to 'string'/'base64' for status '200'


## GET /cat/instruction-stats/{id}/executions
- :warning: removed the optional property '/allOf[subschema #2]/data/items/alarm' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/executed_on' from the response with the '200' status


## POST /color-themes
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /component-alarms
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status


## GET /data-storage
- :warning: removed the optional property 'history/alarm' from the response with the '200' status
- :warning: removed the optional property 'history/alarm_external_tag' from the response with the '200' status
- :warning: removed the optional property 'history/entity_cleaned' from the response with the '200' status
- :warning: removed the optional property 'history/entity_disabled' from the response with the '200' status
- :warning: removed the optional property 'history/entity_unlinked' from the response with the '200' status
- :warning: removed the optional property 'history/event_filter_failure' from the response with the '200' status
- :warning: removed the optional property 'history/event_records' from the response with the '200' status
- :warning: removed the optional property 'history/health_check' from the response with the '200' status
- :warning: removed the optional property 'history/junit' from the response with the '200' status
- :warning: removed the optional property 'history/pbehavior' from the response with the '200' status
- :warning: removed the optional property 'history/remediation' from the response with the '200' status
- :warning: removed the optional property 'history/webhook' from the response with the '200' status


## PUT /data-storage
- :warning: added the new required request property 'event_records/delete_after/enabled'
- :warning: removed the optional property 'history/alarm' from the response with the '200' status
- :warning: removed the optional property 'history/alarm_external_tag' from the response with the '200' status
- :warning: removed the optional property 'history/entity_cleaned' from the response with the '200' status
- :warning: removed the optional property 'history/entity_disabled' from the response with the '200' status
- :warning: removed the optional property 'history/entity_unlinked' from the response with the '200' status
- :warning: removed the optional property 'history/event_filter_failure' from the response with the '200' status
- :warning: removed the optional property 'history/event_records' from the response with the '200' status
- :warning: removed the optional property 'history/health_check' from the response with the '200' status
- :warning: removed the optional property 'history/junit' from the response with the '200' status
- :warning: removed the optional property 'history/pbehavior' from the response with the '200' status
- :warning: removed the optional property 'history/remediation' from the response with the '200' status
- :warning: removed the optional property 'history/webhook' from the response with the '200' status


## POST /entity-comments
- :warning: removed the optional property 'error' from the response with the '400' status


## PUT /entity-comments/{id}
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /entityservice-alarms/{id}
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status


## POST /event
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /eventfilter/rules
- :warning: the '/allOf[subschema #2]/data/items/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## POST /eventfilter/rules
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '201'
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /eventfilter/rules/{id}
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## PUT /eventfilter/rules/{id}
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## POST /flapping-rules
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /link-rules
- :warning: the '/allOf[subschema #2]/data/items/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## POST /link-rules
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '201'
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /link-rules/{id}
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## PUT /link-rules/{id}
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## GET /open-alarms
- :warning: removed the optional property 'v/infos_rule_version' from the response with the '200' status


## GET /permissions
- :warning: removed the optional property '/allOf[subschema #2]/data/items/playlist' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/view' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/view_group' from the response with the '200' status


## POST /resolve-rules
- :warning: removed the optional property 'error' from the response with the '400' status


## GET /resolved-alarms
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status


## GET /role-templates
- :warning: removed the optional property 'data/items/permissions/items/description' from the response with the '200' status


## GET /roles
- :warning: removed the optional property '/allOf[subschema #2]/data/items/permissions/items/description' from the response with the '200' status


## POST /roles
- :warning: added the new required request property 'type'
- :warning: removed the optional property 'permissions/items/description' from the response with the '201' status


## GET /roles/{id}
- :warning: removed the optional property 'permissions/items/description' from the response with the '200' status


## PUT /roles/{id}
- :warning: added the new required request property 'name'
- :warning: added the new required request property 'type'
- :warning: removed the optional property 'permissions/items/description' from the response with the '200' status


## POST /state-settings
- :warning: removed the optional property 'error' from the response with the '400' status


## POST /template-validator/declare-ticket-rules
- :warning: removed the optional property 'error' from the response with the '400' status


## POST /template-validator/event-filter-rules
- :warning: removed the optional property 'error' from the response with the '400' status


## POST /template-validator/scenarios
- :warning: removed the optional property 'error' from the response with the '400' status



