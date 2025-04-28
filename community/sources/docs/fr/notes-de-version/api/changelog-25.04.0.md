# API Changelog 24.10.0 vs. 25.04.0

## GET /account/me
- :warning: removed the optional property 'permissions/items/description' from the response with the '200' status
-  added the optional property 'idp_fields' to the response with the '200' status
-  added the optional property 'idp_roles' to the response with the '200' status
-  added the optional property 'permissions/items/view' to the response with the '200' status


## PUT /account/me
- :warning: removed the success response with the status '201'
-  added the success response with the status '200'


## POST /alarm-details
- :warning: removed the optional property 'data/children/data/items/v/infos_rule_version' from the response with the '207' status
-  added the optional property 'data/children/data/items/v/close_delay' to the response with the '207' status
-  added the optional property 'data/children/data/items/v/close_delay_value' to the response with the '207' status
-  added the optional property 'data/children/data/items/v/comments' to the response with the '207' status


## GET /alarm-tag-labels
-  endpoint added


## GET /alarm-tags
- :warning: removed the success response with the status '201'
-  added the success response with the status '200'


## POST /alarm-tags
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /alarms
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay_value' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/comments' to the response with the '200' status


## GET /alarms/{id}
- :warning: removed the optional property 'v/infos_rule_version' from the response with the '200' status
-  added the optional property 'v/close_delay' to the response with the '200' status
-  added the optional property 'v/close_delay_value' to the response with the '200' status
-  added the optional property 'v/comments' to the response with the '200' status


## PUT /alarms/{id}/ack
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /alarms/{id}/ackremove
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /alarms/{id}/assocticket
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /alarms/{id}/cancel
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /alarms/{id}/changestate
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /alarms/{id}/comment
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /alarms/{id}/snooze
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /alarms/{id}/uncancel
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /app-info
- :warning: removed the optional property 'stack' from the response with the '200' status
-  added the optional property 'file_import_max_size' to the response with the '200' status
-  added the optional property 'user_timezones' to the response with the '200' status
-  added the optional property 'version_description' to the response with the '200' status


## POST /associativetable
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## POST /broadcast-message
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## POST /bulk/eventfilters
- :warning: the '/items/external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the '/items/items/allOf[subschema #2]/item/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '207'


## PUT /bulk/eventfilters
- :warning: the '/items/external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the '/items/items/allOf[subschema #2]/item/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '207'


## DELETE /bulk/external-data-tables/{table}/data
-  endpoint added


## POST /bulk/pbehaviors
-  added the new optional request property '/items/inherited'
-  added the new optional request property '/items/timezone'
-  added the optional property '/items/items/allOf[subschema #2]/item/inherited' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/timezone' to the response with the '207' status


## PUT /bulk/pbehaviors
-  added the new optional request property '/items/inherited'
-  added the new optional request property '/items/timezone'
-  added the optional property '/items/items/allOf[subschema #2]/item/inherited' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/timezone' to the response with the '207' status


## PUT /bulk/role-permissions
-  endpoint added


## PATCH /bulk/users
-  endpoint added


## GET /cat/alarm-executions/{id}
-  added the new optional 'query' request parameter 'ids[]'


## GET /cat/event-records/{id}
- :warning: the '/allOf[subschema #2]/data/items/event' response's property type/format changed from 'array'/'' to 'string'/'base64' for status '200'


## GET /cat/instruction-stats/{id}/changes
-  added the optional property '/allOf[subschema #2]/data/items/avg_alarm_ok_timeout' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/avg_successful' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/avg_successful_state_ok' to the response with the '200' status


## GET /cat/instruction-stats/{id}/executions
- :warning: removed the optional property '/allOf[subschema #2]/data/items/alarm' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/executed_on' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/_id' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_display_name' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_id' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_ok_at' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_ok_before_completed' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_ok_timeout' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_steps' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/completed_at' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/instruction_type' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/result_alarm_state' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/started_at' to the response with the '200' status


## GET /cat/instruction-stats/{id}/summary
-  added the optional property 'avg_alarm_ok_timeout' to the response with the '200' status
-  added the optional property 'avg_successful' to the response with the '200' status
-  added the optional property 'avg_successful_state_ok' to the response with the '200' status


## GET /cat/instructions
-  added the optional property '/allOf[subschema #2]/data/items/jobs/items/job/created' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/jobs/items/job/updated' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/steps/items/operations/items/jobs/items/created' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/steps/items/operations/items/jobs/items/updated' to the response with the '200' status


## POST /cat/instructions
-  added the optional property 'jobs/items/job/created' to the response with the '201' status
-  added the optional property 'jobs/items/job/updated' to the response with the '201' status
-  added the optional property 'steps/items/operations/items/jobs/items/created' to the response with the '201' status
-  added the optional property 'steps/items/operations/items/jobs/items/updated' to the response with the '201' status


## GET /cat/instructions/{id}
-  added the optional property 'jobs/items/job/created' to the response with the '200' status
-  added the optional property 'jobs/items/job/updated' to the response with the '200' status
-  added the optional property 'steps/items/operations/items/jobs/items/created' to the response with the '200' status
-  added the optional property 'steps/items/operations/items/jobs/items/updated' to the response with the '200' status


## PUT /cat/instructions/{id}
-  added the optional property 'jobs/items/job/created' to the response with the '200' status
-  added the optional property 'jobs/items/job/updated' to the response with the '200' status
-  added the optional property 'steps/items/operations/items/jobs/items/created' to the response with the '200' status
-  added the optional property 'steps/items/operations/items/jobs/items/updated' to the response with the '200' status


## GET /cat/instructions/{id}/approval
-  added the optional property 'original/jobs/items/job/created' to the response with the '200' status
-  added the optional property 'original/jobs/items/job/updated' to the response with the '200' status
-  added the optional property 'original/steps/items/operations/items/jobs/items/created' to the response with the '200' status
-  added the optional property 'original/steps/items/operations/items/jobs/items/updated' to the response with the '200' status
-  added the optional property 'updated/jobs/items/job/created' to the response with the '200' status
-  added the optional property 'updated/jobs/items/job/updated' to the response with the '200' status
-  added the optional property 'updated/steps/items/operations/items/jobs/items/created' to the response with the '200' status
-  added the optional property 'updated/steps/items/operations/items/jobs/items/updated' to the response with the '200' status


## PUT /cat/instructions/{id}/approval
-  added the optional property 'jobs/items/job/created' to the response with the '200' status
-  added the optional property 'jobs/items/job/updated' to the response with the '200' status
-  added the optional property 'steps/items/operations/items/jobs/items/created' to the response with the '200' status
-  added the optional property 'steps/items/operations/items/jobs/items/updated' to the response with the '200' status


## GET /cat/jobs
-  added the optional property '/allOf[subschema #2]/data/items/created' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/updated' to the response with the '200' status


## POST /cat/jobs
-  added the optional property 'created' to the response with the '201' status
-  added the optional property 'updated' to the response with the '201' status


## GET /cat/jobs/{id}
-  added the optional property 'created' to the response with the '200' status
-  added the optional property 'updated' to the response with the '200' status


## PUT /cat/jobs/{id}
-  added the optional property 'created' to the response with the '200' status
-  added the optional property 'updated' to the response with the '200' status


## POST /cat/meta-alarms
-  added the new optional request property 'component'
-  added the new optional request property 'infos'
-  added the new optional request property 'resource'
-  added the new optional request property 'tags'


## GET /cat/metaalarmrules
-  added the optional property '/allOf[subschema #2]/data/items/config/component_template' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/config/resource_template' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/infos' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/tags' to the response with the '200' status


## POST /cat/metaalarmrules
-  added the new optional request property 'config/component_template'
-  added the new optional request property 'config/resource_template'
-  added the new optional request property 'infos'
-  added the new optional request property 'tags'
-  the 'output_template' request property's maxLength was increased from '500' to '10000'
-  added the optional property 'config/component_template' to the response with the '201' status
-  added the optional property 'config/resource_template' to the response with the '201' status
-  added the optional property 'infos' to the response with the '201' status
-  added the optional property 'tags' to the response with the '201' status


## GET /cat/metaalarmrules/{id}
-  added the optional property 'config/component_template' to the response with the '200' status
-  added the optional property 'config/resource_template' to the response with the '200' status
-  added the optional property 'infos' to the response with the '200' status
-  added the optional property 'tags' to the response with the '200' status


## PUT /cat/metaalarmrules/{id}
-  added the new optional request property 'config/component_template'
-  added the new optional request property 'config/resource_template'
-  added the new optional request property 'infos'
-  added the new optional request property 'tags'
-  the 'output_template' request property's maxLength was increased from '500' to '10000'
-  added the optional property 'config/component_template' to the response with the '200' status
-  added the optional property 'config/resource_template' to the response with the '200' status
-  added the optional property 'infos' to the response with the '200' status
-  added the optional property 'tags' to the response with the '200' status


## POST /color-themes
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /component-alarms
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay_value' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/comments' to the response with the '200' status


## PUT /contextgraph-import
-  added the new optional request property '/items/tags'


## PUT /contextgraph-import-partial
-  added the new optional request property '/items/tags'


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
-  added the required property 'config/event_records/delete_after/enabled' to the response with the '200' status


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
-  added the required property 'config/event_records/delete_after/enabled' to the response with the '200' status


## GET /entities
-  added the new optional 'query' request parameter 'ids[]'


## GET /entities/pbehaviors
-  added the optional property '/items/inherited' to the response with the '200' status
-  added the optional property '/items/timezone' to the response with the '200' status


## GET /entity-categories
-  added the new optional 'query' request parameter 'ids[]'


## POST /entity-comments
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## PUT /entity-comments/{id}
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## POST /entity-export
-  added the new optional request property 'ids'


## GET /entityservice-alarms/{id}
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay_value' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/comments' to the response with the '200' status


## POST /event
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /eventfilter/rules
- :warning: the '/allOf[subschema #2]/data/items/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## POST /eventfilter/rules
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '201'
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /eventfilter/rules/{id}
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## PUT /eventfilter/rules/{id}
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## GET /external-data-export/{id}
-  endpoint added


## POST /external-data-export/{id}
-  endpoint added


## GET /external-data-export/{id}/download
-  endpoint added


## POST /external-data-import/{id}
-  endpoint added


## PUT /external-data-import/{id}/complete
-  endpoint added


## GET /external-data-import/{id}/data
-  endpoint added


## GET /external-data-import/{id}/status
-  endpoint added


## GET /external-data-tables
-  endpoint added


## POST /external-data-tables
-  endpoint added


## DELETE /external-data-tables/{id}
-  endpoint added


## GET /external-data-tables/{id}
-  endpoint added


## PUT /external-data-tables/{id}
-  endpoint added


## GET /external-data-tables/{table}/data
-  endpoint added


## POST /external-data-tables/{table}/data
-  endpoint added


## DELETE /external-data-tables/{table}/data/{id}
-  endpoint added


## GET /external-data-tables/{table}/data/{id}
-  endpoint added


## PUT /external-data-tables/{table}/data/{id}
-  endpoint added


## GET /external-data-tables/{table}/schema
-  endpoint added


## POST /flapping-rules
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /icons
-  added the optional property '/allOf[subschema #2]/data/items/fill_border' to the response with the '200' status


## POST /icons
-  added the optional property '/items/fill_border' to the response with the '201' status


## PATCH /icons/{id}
-  added the optional property '/items/fill_border' to the response with the '200' status


## PUT /icons/{id}
-  added the optional property '/items/fill_border' to the response with the '200' status


## PUT /internal/user_interface
-  added the new optional request property 'default_color_theme'
-  added the new optional request property 'version_description'
-  added the optional property 'default_color_theme' to the response with the '200' status
-  added the optional property 'version_description' to the response with the '200' status


## GET /link-rules
- :warning: the '/allOf[subschema #2]/data/items/external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## POST /link-rules
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '201'
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /link-rules/{id}
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## PUT /link-rules/{id}
- :warning: the 'external_data' request property type/format changed from 'object'/'' to 'array'/''
- :warning: the 'external_data' response's property type/format changed from 'object'/'' to 'array'/'' for status '200'


## GET /open-alarms
- :warning: removed the optional property 'v/infos_rule_version' from the response with the '200' status
-  added the optional property 'v/close_delay' to the response with the '200' status
-  added the optional property 'v/close_delay_value' to the response with the '200' status
-  added the optional property 'v/comments' to the response with the '200' status


## GET /pbehavior-reasons
-  added the new optional 'query' request parameter 'ids[]'


## POST /pbehavior-timespans
-  added the new optional request property 'timezone'


## GET /pbehavior-types
-  added the new optional 'query' request parameter 'ids[]'


## GET /pbehaviors
-  added the new optional 'query' request parameter 'ids[]'
-  added the optional property '/allOf[subschema #2]/data/items/inherited' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/timezone' to the response with the '200' status


## POST /pbehaviors
-  added the new optional request property 'inherited'
-  added the new optional request property 'timezone'
-  added the optional property 'inherited' to the response with the '201' status
-  added the optional property 'timezone' to the response with the '201' status


## GET /pbehaviors/{id}
-  added the optional property 'inherited' to the response with the '200' status
-  added the optional property 'timezone' to the response with the '200' status


## PATCH /pbehaviors/{id}
-  added the new optional request property 'inherited'
-  added the new optional request property 'timezone'
-  added the optional property 'inherited' to the response with the '200' status
-  added the optional property 'timezone' to the response with the '200' status


## PUT /pbehaviors/{id}
-  added the new optional request property 'inherited'
-  added the new optional request property 'timezone'
-  added the optional property 'inherited' to the response with the '200' status
-  added the optional property 'timezone' to the response with the '200' status


## GET /permissions
- :warning: removed the optional property '/allOf[subschema #2]/data/items/playlist' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/view' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/view_group' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/groups' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/title' to the response with the '200' status


## POST /resolve-rules
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /resolved-alarms
- :warning: removed the optional property '/allOf[subschema #2]/data/items/v/infos_rule_version' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/close_delay_value' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/comments' to the response with the '200' status


## GET /role-templates
- :warning: removed the optional property 'data/items/permissions/items/description' from the response with the '200' status
-  added the optional property 'data/items/permissions/items/view' to the response with the '200' status
-  added the optional property 'data/items/type' to the response with the '200' status


## GET /roles
- :warning: removed the optional property '/allOf[subschema #2]/data/items/permissions/items/description' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/permissions/items/view' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/type' to the response with the '200' status


## POST /roles
- :warning: added the new required request property 'type'
- :warning: removed the optional property 'permissions/items/description' from the response with the '201' status
-  added the optional property 'permissions/items/view' to the response with the '201' status
-  added the optional property 'type' to the response with the '201' status


## GET /roles/{id}
- :warning: removed the optional property 'permissions/items/description' from the response with the '200' status
-  added the optional property 'permissions/items/view' to the response with the '200' status
-  added the optional property 'type' to the response with the '200' status


## PUT /roles/{id}
- :warning: added the new required request property 'name'
- :warning: added the new required request property 'type'
- :warning: removed the optional property 'permissions/items/description' from the response with the '200' status
-  added the optional property 'permissions/items/view' to the response with the '200' status
-  added the optional property 'type' to the response with the '200' status


## POST /state-settings
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## POST /template-validator/declare-ticket-rules
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## POST /template-validator/event-filter-rules
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## POST /template-validator/scenarios
- :warning: removed the optional property 'error' from the response with the '400' status
-  added the optional property 'errors' to the response with the '400' status


## GET /users
-  added the new optional 'query' request parameter 'ids[]'
-  added the optional property '/allOf[subschema #2]/data/items/idp_fields' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/idp_roles' to the response with the '200' status


## POST /users
-  added the optional property 'idp_fields' to the response with the '201' status
-  added the optional property 'idp_roles' to the response with the '201' status


## GET /users/{id}
-  added the optional property 'idp_fields' to the response with the '200' status
-  added the optional property 'idp_roles' to the response with the '200' status


## PATCH /users/{id}
-  added the optional property 'idp_fields' to the response with the '200' status
-  added the optional property 'idp_roles' to the response with the '200' status


## PUT /users/{id}
-  added the optional property 'idp_fields' to the response with the '200' status
-  added the optional property 'idp_roles' to the response with the '200' status


## POST /view-copy/{id}
-  added the optional property 'tabs/items/widgets/items/parameters/table' to the response with the '201' status


## POST /view-export
-  added the optional property 'groups/items/views/items/tabs/items/widgets/items/parameters/table' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/parameters/table' to the response with the '200' status


## GET /view-groups
-  added the optional property '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/parameters/table' to the response with the '200' status


## POST /view-groups
-  added the optional property 'views/items/tabs/items/widgets/items/parameters/table' to the response with the '201' status


## GET /view-groups/{id}
-  added the optional property 'views/items/tabs/items/widgets/items/parameters/table' to the response with the '200' status


## PUT /view-groups/{id}
-  added the optional property 'views/items/tabs/items/widgets/items/parameters/table' to the response with the '200' status


## POST /view-import
-  added the new optional request property '/items/views/items/tabs/items/widgets/items/parameters/table'


## POST /view-tab-copy/{id}
-  added the optional property 'widgets/items/parameters/table' to the response with the '201' status


## POST /view-tabs
-  added the optional property 'widgets/items/parameters/table' to the response with the '201' status


## GET /view-tabs/{id}
-  added the optional property 'widgets/items/parameters/table' to the response with the '200' status


## PUT /view-tabs/{id}
-  added the optional property 'widgets/items/parameters/table' to the response with the '200' status


## POST /views
-  added the optional property 'tabs/items/widgets/items/parameters/table' to the response with the '201' status


## GET /views/{id}
-  added the optional property 'tabs/items/widgets/items/parameters/table' to the response with the '200' status


## PUT /views/{id}
-  added the optional property 'tabs/items/widgets/items/parameters/table' to the response with the '200' status


## POST /widget-copy/{id}
-  added the new optional request property 'parameters/table'
-  added the optional property 'parameters/table' to the response with the '201' status


## POST /widgets
-  added the new optional request property 'parameters/table'
-  added the optional property 'parameters/table' to the response with the '201' status


## GET /widgets/{id}
-  added the optional property 'parameters/table' to the response with the '200' status


## PUT /widgets/{id}
-  added the new optional request property 'parameters/table'
-  added the optional property 'parameters/table' to the response with the '200' status



