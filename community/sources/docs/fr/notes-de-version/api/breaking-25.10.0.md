# API Changelog 25.04.2 vs. 25.10.0

## GET /alarm-counters
- :warning: deleted the 'query' request parameter 'instructions[]'


## POST /alarm-details
- :warning: the 'data/children/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the 'data/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## POST /alarm-export
- :warning: removed the request property 'instructions'
- :warning: removed the request property 'tag'


## GET /alarm-tags
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /alarm-tags
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /alarm-tags/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /alarm-tags/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /alarms
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: deleted the 'query' request parameter 'instructions[]'


## GET /alarms/{id}
- :warning: the 'entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /bulk/entityservices
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## PUT /bulk/entityservices
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## POST /bulk/eventfilters
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## PUT /bulk/eventfilters
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## POST /bulk/idle-rules
- :warning: the '/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## PUT /bulk/idle-rules
- :warning: the '/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## POST /bulk/pbehaviors
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## PUT /bulk/pbehaviors
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## POST /bulk/scenarios
- :warning: the '/items/actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## PUT /bulk/scenarios
- :warning: the '/items/actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'


## PUT /cat/account/paused-executions
- :warning: api path removed without deprecation


## GET /cat/declare-ticket-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /cat/declare-ticket-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the request property 'corporate_weather_service_pattern'
- :warning: removed the request property 'weather_service_pattern'


## GET /cat/declare-ticket-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /cat/declare-ticket-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the request property 'corporate_weather_service_pattern'
- :warning: removed the request property 'weather_service_pattern'


## GET /cat/dynamic-infos
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/alarm_update' from the response with the '200' status


## POST /cat/dynamic-infos
- :warning: added the new required request property 'infos/items/type'
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the optional property 'alarm_update' from the response with the '201' status


## GET /cat/dynamic-infos/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'alarm_update' from the response with the '200' status


## PUT /cat/dynamic-infos/{id}
- :warning: added the new required request property 'infos/items/type'
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'alarm_update' from the response with the '200' status


## GET /cat/event-records
- :warning: the '/allOf[subschema #2]/data/items/pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /cat/event-records-current
- :warning: the 'event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## GET /cat/instruction-stats/{id}/summary
- :warning: removed the optional property 'last_modified' from the response with the '200' status


## GET /cat/instructions
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/last_modified' from the response with the '200' status


## POST /cat/instructions
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the optional property 'last_modified' from the response with the '201' status


## GET /cat/instructions/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'last_modified' from the response with the '200' status


## PUT /cat/instructions/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'last_modified' from the response with the '200' status


## GET /cat/instructions/{id}/approval
- :warning: the 'original/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'original/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'updated/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'updated/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'original/last_modified' from the response with the '200' status
- :warning: removed the optional property 'updated/last_modified' from the response with the '200' status


## PUT /cat/instructions/{id}/approval
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'last_modified' from the response with the '200' status


## GET /cat/kpi-filters
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /cat/kpi-filters
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /cat/kpi-filters/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /cat/kpi-filters/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /cat/map-state/{id}
- :warning: the 'parameters/entities/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'parameters/entities/items/pinned_entities/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'parameters/expanded_entities/additionalProperties/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'parameters/points/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /cat/metaalarmrules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /cat/metaalarmrules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'total_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /cat/metaalarmrules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /cat/metaalarmrules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'total_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /cat/metrics-export/group
- :warning: the 'entity_patterns/items/pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## POST /cat/metrics/group
- :warning: the 'entity_patterns/items/pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## POST /cat/test-scenario-executions
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## POST /color-themes
- :warning: added the new required request property 'colors/main/error_background'
- :warning: added the new required request property 'colors/main/info_background'
- :warning: added the new required request property 'colors/main/success_background'
- :warning: added the new required request property 'colors/main/warning_background'
- :warning: added the new required request property 'colors/table/hover_row'
- :warning: added the new required request property 'colors/table/shift_row'


## PUT /color-themes/{id}
- :warning: added the new required request property 'colors/main/error_background'
- :warning: added the new required request property 'colors/main/info_background'
- :warning: added the new required request property 'colors/main/success_background'
- :warning: added the new required request property 'colors/main/warning_background'
- :warning: added the new required request property 'colors/table/hover_row'
- :warning: added the new required request property 'colors/table/shift_row'


## GET /component-alarms
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /contextgraph-import
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## PUT /contextgraph-import-partial
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## GET /entities
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /entities/check-state-setting
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /entities/pbehaviors
- :warning: the '/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /entities/state-setting
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /entitybasics
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /entitybasics
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /entityservice-alarms/{id}
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /entityservice-dependencies
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /entityservice-impacts
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /entityservices
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /entityservices/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /entityservices/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /eventfilter/rules
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/columns' from the response with the '200' status


## POST /eventfilter/rules
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the request property 'external_data/items/table/column_types'
- :warning: removed the request property 'external_data/items/table/columns'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '201' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '201' status


## GET /eventfilter/rules/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status


## PUT /eventfilter/rules/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the request property 'external_data/items/table/column_types'
- :warning: removed the request property 'external_data/items/table/columns'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status


## POST /external-data-import/{id}
- :warning: removed the optional property 'columns' from the response with the '200' status


## PUT /external-data-import/{id}/complete
- :warning: added the new required request property 'column_tags'
- :warning: removed the request property 'column_types'


## GET /external-data-import/{id}/status
- :warning: removed the optional property 'columns' from the response with the '200' status


## GET /external-data-tables
- :warning: removed the optional property '/allOf[subschema #2]/data/items/column_types' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/columns' from the response with the '200' status


## POST /external-data-tables
- :warning: removed the optional property 'column_types' from the response with the '201' status
- :warning: removed the optional property 'columns' from the response with the '201' status


## GET /external-data-tables/{id}
- :warning: removed the optional property 'column_types' from the response with the '200' status
- :warning: removed the optional property 'columns' from the response with the '200' status


## PUT /external-data-tables/{id}
- :warning: added the new required request property 'column_tags'
- :warning: removed the request property 'column_types'
- :warning: removed the optional property 'column_types' from the response with the '200' status
- :warning: removed the optional property 'columns' from the response with the '200' status


## GET /flapping-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /flapping-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /flapping-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /flapping-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /idle-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /idle-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /idle-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /idle-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /link-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/columns' from the response with the '200' status


## POST /link-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '201' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '201' status


## GET /link-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status


## PUT /link-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status


## GET /notification
- :warning: api path removed without deprecation


## PUT /notification
- :warning: api path removed without deprecation


## GET /open-alarms
- :warning: the 'entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /patterns
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /patterns
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## POST /patterns-alarms-count
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## POST /patterns-entities-count
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## GET /patterns/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /patterns/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /pbehaviors
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /pbehaviors
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /pbehaviors/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PATCH /pbehaviors/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /pbehaviors/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /pbehaviors/{id}/entities
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /resolve-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /resolve-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /resolve-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /resolve-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /resolved-alarms
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /scenarios
- :warning: the '/allOf[subschema #2]/data/items/actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /scenarios
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /scenarios/{id}
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /scenarios/{id}
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /state-settings
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /state-settings
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'inherited_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /state-settings/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /state-settings/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'inherited_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /template-validator/declare-ticket-rules
- :warning: api path removed without deprecation


## POST /template-validator/event-filter-rules
- :warning: api path removed without deprecation


## POST /template-validator/scenarios
- :warning: api path removed without deprecation


## PUT /user-preferences
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /user-preferences/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /view-copy/{id}
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## POST /view-export
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## GET /view-groups
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /view-groups
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /view-groups/{id}
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /view-groups/{id}
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /view-import
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''


## POST /view-tab-copy/{id}
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## POST /view-tabs
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /view-tabs/{id}
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /view-tabs/{id}
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /views
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /views/{id}
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /views/{id}
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /widget-copy/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /widget-filters
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /widget-filters
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /widget-filters/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /widget-filters/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## POST /widgets
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'


## GET /widgets/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'


## PUT /widgets/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'



