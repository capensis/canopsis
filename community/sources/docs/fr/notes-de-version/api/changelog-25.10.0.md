# API Changelog 25.04.2 vs. 25.10.0

## GET /account/me
-  added the required property 'ui_theme_colors/colors/main/error_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/main/info_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/main/success_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/main/warning_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/table/hover_row' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/table/shift_row' to the response with the '200' status


## PUT /account/me
-  added the required property 'ui_theme_colors/colors/main/error_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/main/info_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/main/success_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/main/warning_background' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/table/hover_row' to the response with the '200' status
-  added the required property 'ui_theme_colors/colors/table/shift_row' to the response with the '200' status


## GET /active-broadcast-message
-  added the optional property '/items/views' to the response with the '200' status


## GET /alarm-counters
- :warning: deleted the 'query' request parameter 'instructions[]'
-  added the new optional 'query' request parameter 'instruction_filter_type'
-  added the new optional 'query' request parameter 'instruction_ids[]'
-  added the new optional 'query' request parameter 'instruction_statuses[]'
-  added the new optional 'query' request parameter 'instruction_type'


## POST /alarm-details
- :warning: the 'data/children/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the 'data/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the optional property 'data/children/data/items/entity/entity_pattern/items/items/alias' to the response with the '207' status
-  added the optional property 'data/children/data/items/v/initial_state' to the response with the '207' status
-  added the optional property 'data/children/data/items/v/max_state' to the response with the '207' status
-  added the optional property 'data/entity/entity_pattern/items/items/alias' to the response with the '207' status


## GET /alarm-display-names
-  added the new optional 'query' request parameter 'opened'


## POST /alarm-export
- :warning: removed the request property 'instructions'
- :warning: removed the request property 'tag'
-  added the new optional request property 'instruction_filter_type'
-  added the new optional request property 'instruction_ids'
-  added the new optional request property 'instruction_statuses'
-  added the new optional request property 'instruction_type'


## GET /alarm-tags
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## POST /alarm-tags
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status


## GET /alarm-tags/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## PUT /alarm-tags/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## GET /alarms
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: deleted the 'query' request parameter 'instructions[]'
-  added the new optional 'query' request parameter 'instruction_filter_type'
-  added the new optional 'query' request parameter 'instruction_ids[]'
-  added the new optional 'query' request parameter 'instruction_statuses[]'
-  added the new optional 'query' request parameter 'instruction_type'
-  added the optional property '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/initial_state' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/max_state' to the response with the '200' status


## GET /alarms/{id}
- :warning: the 'entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'v/initial_state' to the response with the '200' status
-  added the optional property 'v/max_state' to the response with the '200' status


## PUT /all-pbehavior-patterns
-  endpoint added


## GET /app-info
-  added the optional property 'notification_display_count' to the response with the '200' status
-  added the required property 'default_color_theme/colors/main/error_background' to the response with the '200' status
-  added the required property 'default_color_theme/colors/main/info_background' to the response with the '200' status
-  added the required property 'default_color_theme/colors/main/success_background' to the response with the '200' status
-  added the required property 'default_color_theme/colors/main/warning_background' to the response with the '200' status
-  added the required property 'default_color_theme/colors/table/hover_row' to the response with the '200' status
-  added the required property 'default_color_theme/colors/table/shift_row' to the response with the '200' status


## GET /broadcast-message
-  added the optional property '/allOf[subschema #2]/data/items/views' to the response with the '200' status


## POST /broadcast-message
-  added the new optional request property 'views'
-  added the optional property 'views' to the response with the '201' status


## GET /broadcast-message/{id}
-  added the optional property 'views' to the response with the '200' status


## PUT /broadcast-message/{id}
-  added the new optional request property 'views'
-  added the optional property 'views' to the response with the '200' status


## PUT /broadcast-message/{id}/read
-  endpoint added


## DELETE /bulk/entity-infos-properties
-  endpoint added


## POST /bulk/entityservices
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status


## PUT /bulk/entityservices
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status


## POST /bulk/eventfilters
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the new optional request property '/items/event_pattern/items/items/alias'
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/event_pattern/items/items/alias' to the response with the '207' status


## PUT /bulk/eventfilters
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the new optional request property '/items/event_pattern/items/items/alias'
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/event_pattern/items/items/alias' to the response with the '207' status


## POST /bulk/idle-rules
- :warning: the '/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/alarm_pattern/items/items/alias'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the optional property '/items/items/allOf[subschema #2]/item/alarm_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status


## PUT /bulk/idle-rules
- :warning: the '/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/alarm_pattern/items/items/alias'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the optional property '/items/items/allOf[subschema #2]/item/alarm_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status


## POST /bulk/pbehaviors
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the new optional request property '/items/exec_pattern'
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/exec_pattern' to the response with the '207' status


## PUT /bulk/pbehaviors
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/entity_pattern/items/items/alias'
-  added the new optional request property '/items/exec_pattern'
-  added the optional property '/items/items/allOf[subschema #2]/item/entity_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/exec_pattern' to the response with the '207' status


## POST /bulk/scenarios
- :warning: the '/items/actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/actions/items/alarm_pattern/items/items/alias'
-  added the new optional request property '/items/actions/items/entity_pattern/items/items/alias'
-  added the new optional request property '/items/actions/items/parameters/auth_token'
-  added the new optional request property '/items/actions/items/parameters/multiple_urls'
-  added the new optional request property '/items/actions/items/parameters/stop_on_fail'
-  added the new optional request property '/items/actions/items/parameters/stop_on_success'
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/alarm_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/entity_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/auth_token' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/multiple_urls' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/stop_on_fail' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/stop_on_success' to the response with the '207' status


## PUT /bulk/scenarios
- :warning: the '/items/actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
- :warning: the '/items/items/allOf[subschema #2]/item/actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '207'
-  added the new optional request property '/items/actions/items/alarm_pattern/items/items/alias'
-  added the new optional request property '/items/actions/items/entity_pattern/items/items/alias'
-  added the new optional request property '/items/actions/items/parameters/auth_token'
-  added the new optional request property '/items/actions/items/parameters/multiple_urls'
-  added the new optional request property '/items/actions/items/parameters/stop_on_fail'
-  added the new optional request property '/items/actions/items/parameters/stop_on_success'
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/alarm_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/entity_pattern/items/items/alias' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/auth_token' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/multiple_urls' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/stop_on_fail' to the response with the '207' status
-  added the optional property '/items/items/allOf[subschema #2]/item/actions/items/parameters/stop_on_success' to the response with the '207' status


## PUT /cat/account/executions
-  endpoint added


## PUT /cat/account/paused-executions
- :warning: api path removed without deprecation


## POST /cat/declare-ticket-rule-template-validate
-  endpoint added


## GET /cat/declare-ticket-rule-template-vars
-  endpoint added


## GET /cat/declare-ticket-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/webhooks/items/auth_token' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/webhooks/items/multiple_urls' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/webhooks/items/stop_on_success' to the response with the '200' status


## POST /cat/declare-ticket-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the request property 'corporate_weather_service_pattern'
- :warning: removed the request property 'weather_service_pattern'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the new optional request property 'webhooks/items/auth_token'
-  added the new optional request property 'webhooks/items/multiple_urls'
-  added the new optional request property 'webhooks/items/stop_on_success'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'webhooks/items/auth_token' to the response with the '201' status
-  added the optional property 'webhooks/items/multiple_urls' to the response with the '201' status
-  added the optional property 'webhooks/items/stop_on_success' to the response with the '201' status


## GET /cat/declare-ticket-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'webhooks/items/auth_token' to the response with the '200' status
-  added the optional property 'webhooks/items/multiple_urls' to the response with the '200' status
-  added the optional property 'webhooks/items/stop_on_success' to the response with the '200' status


## PUT /cat/declare-ticket-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the request property 'corporate_weather_service_pattern'
- :warning: removed the request property 'weather_service_pattern'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the new optional request property 'webhooks/items/auth_token'
-  added the new optional request property 'webhooks/items/multiple_urls'
-  added the new optional request property 'webhooks/items/stop_on_success'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'webhooks/items/auth_token' to the response with the '200' status
-  added the optional property 'webhooks/items/multiple_urls' to the response with the '200' status
-  added the optional property 'webhooks/items/stop_on_success' to the response with the '200' status


## GET /cat/dynamic-infos
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/alarm_update' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/infos/items/type' to the response with the '200' status


## POST /cat/dynamic-infos
- :warning: added the new required request property 'infos/items/type'
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the optional property 'alarm_update' from the response with the '201' status
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'infos/items/type' to the response with the '201' status


## GET /cat/dynamic-infos-copy-vars
-  endpoint added


## POST /cat/dynamic-infos-template-validate
-  endpoint added


## GET /cat/dynamic-infos-template-vars
-  endpoint added


## GET /cat/dynamic-infos/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'alarm_update' from the response with the '200' status
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'infos/items/type' to the response with the '200' status


## PUT /cat/dynamic-infos/{id}
- :warning: added the new required request property 'infos/items/type'
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'alarm_update' from the response with the '200' status
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'infos/items/type' to the response with the '200' status


## GET /cat/entity-infos-log
-  endpoint added


## GET /cat/event-records
- :warning: the '/allOf[subschema #2]/data/items/pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/pattern/items/items/alias' to the response with the '200' status


## POST /cat/event-records-current
- :warning: the 'event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property 'event_pattern/items/items/alias'


## GET /cat/execution-statuses
-  endpoint added


## POST /cat/executions
-  added the optional property 'alarm' to the response with the '200' status
-  added the optional property 'entity' to the response with the '200' status


## GET /cat/executions/{id}
-  added the optional property 'alarm' to the response with the '200' status
-  added the optional property 'entity' to the response with the '200' status


## PUT /cat/executions/{id}/next
-  added the optional property 'alarm' to the response with the '200' status
-  added the optional property 'entity' to the response with the '200' status


## PUT /cat/executions/{id}/next-step
-  added the optional property 'alarm' to the response with the '200' status
-  added the optional property 'entity' to the response with the '200' status


## PUT /cat/executions/{id}/previous
-  added the optional property 'alarm' to the response with the '200' status
-  added the optional property 'entity' to the response with the '200' status


## PUT /cat/executions/{id}/resume
-  added the optional property 'alarm' to the response with the '200' status
-  added the optional property 'entity' to the response with the '200' status


## GET /cat/instruction-stats
-  added the new optional 'query' request parameter 'only_to_rate'


## GET /cat/instruction-stats/{id}/summary
- :warning: removed the optional property 'last_modified' from the response with the '200' status
-  added the optional property 'updated' to the response with the '200' status


## POST /cat/instruction-template-validate
-  endpoint added


## GET /cat/instruction-template-vars
-  endpoint added


## GET /cat/instructions
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/last_modified' from the response with the '200' status
-  added the new optional 'query' request parameter 'only_to_approve'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/repeat_triggers' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/retry_count' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/updated' to the response with the '200' status


## POST /cat/instructions
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the optional property 'last_modified' from the response with the '201' status
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'repeat_triggers'
-  added the new optional request property 'retry_count'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'repeat_triggers' to the response with the '201' status
-  added the optional property 'retry_count' to the response with the '201' status
-  added the optional property 'updated' to the response with the '201' status


## GET /cat/instructions/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'last_modified' from the response with the '200' status
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'repeat_triggers' to the response with the '200' status
-  added the optional property 'retry_count' to the response with the '200' status
-  added the optional property 'updated' to the response with the '200' status


## PUT /cat/instructions/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'last_modified' from the response with the '200' status
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'repeat_triggers'
-  added the new optional request property 'retry_count'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'repeat_triggers' to the response with the '200' status
-  added the optional property 'retry_count' to the response with the '200' status
-  added the optional property 'updated' to the response with the '200' status


## GET /cat/instructions/{id}/approval
- :warning: the 'original/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'original/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'updated/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'updated/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'original/last_modified' from the response with the '200' status
- :warning: removed the optional property 'updated/last_modified' from the response with the '200' status
-  added the optional property 'original/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'original/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'original/repeat_triggers' to the response with the '200' status
-  added the optional property 'original/retry_count' to the response with the '200' status
-  added the optional property 'original/updated' to the response with the '200' status
-  added the optional property 'updated/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'updated/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'updated/repeat_triggers' to the response with the '200' status
-  added the optional property 'updated/retry_count' to the response with the '200' status
-  added the optional property 'updated/updated' to the response with the '200' status


## PUT /cat/instructions/{id}/approval
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'last_modified' from the response with the '200' status
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'repeat_triggers' to the response with the '200' status
-  added the optional property 'retry_count' to the response with the '200' status
-  added the optional property 'updated' to the response with the '200' status


## POST /cat/job-template-validate
-  endpoint added


## GET /cat/job-template-vars
-  endpoint added


## GET /cat/kpi-filters
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## POST /cat/kpi-filters
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status


## GET /cat/kpi-filters/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## PUT /cat/kpi-filters/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## GET /cat/map-state/{id}
- :warning: the 'parameters/entities/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'parameters/entities/items/pinned_entities/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'parameters/expanded_entities/additionalProperties/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'parameters/points/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'parameters/entities/items/entity/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'parameters/entities/items/pinned_entities/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'parameters/points/items/entity/entity_pattern/items/items/alias' to the response with the '200' status


## POST /cat/metaalarmrule-template-validate
-  endpoint added


## GET /cat/metaalarmrule-template-vars
-  endpoint added


## GET /cat/metaalarmrules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/total_entity_pattern/items/items/alias' to the response with the '200' status


## POST /cat/metaalarmrules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'total_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'total_entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'total_entity_pattern/items/items/alias' to the response with the '201' status


## GET /cat/metaalarmrules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'total_entity_pattern/items/items/alias' to the response with the '200' status


## PUT /cat/metaalarmrules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'total_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'total_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'total_entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'total_entity_pattern/items/items/alias' to the response with the '200' status


## POST /cat/metrics-export/group
- :warning: the 'entity_patterns/items/pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property 'entity_patterns/items/pattern/items/items/alias'


## POST /cat/metrics/group
- :warning: the 'entity_patterns/items/pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property 'entity_patterns/items/pattern/items/items/alias'


## POST /cat/test-declare-ticket-executions
-  added the new optional request property 'webhooks/items/auth_token'
-  added the new optional request property 'webhooks/items/multiple_urls'
-  added the new optional request property 'webhooks/items/stop_on_success'


## POST /cat/test-scenario-executions
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property 'actions/items/alarm_pattern/items/items/alias'
-  added the new optional request property 'actions/items/entity_pattern/items/items/alias'
-  added the new optional request property 'actions/items/parameters/auth_token'
-  added the new optional request property 'actions/items/parameters/multiple_urls'
-  added the new optional request property 'actions/items/parameters/stop_on_fail'
-  added the new optional request property 'actions/items/parameters/stop_on_success'


## POST /cat/webhook-token-rule-template-validate
-  endpoint added


## GET /cat/webhook-token-rule-template-vars
-  endpoint added


## GET /cat/webhook-token-rules
-  endpoint added


## POST /cat/webhook-token-rules
-  endpoint added


## POST /cat/webhook-token-rules-db-export
-  endpoint added


## DELETE /cat/webhook-token-rules/{id}
-  endpoint added


## GET /cat/webhook-token-rules/{id}
-  endpoint added


## PUT /cat/webhook-token-rules/{id}
-  endpoint added


## GET /color-themes
-  added the required property '/allOf[subschema #2]/data/items/colors/main/error_background' to the response with the '200' status
-  added the required property '/allOf[subschema #2]/data/items/colors/main/info_background' to the response with the '200' status
-  added the required property '/allOf[subschema #2]/data/items/colors/main/success_background' to the response with the '200' status
-  added the required property '/allOf[subschema #2]/data/items/colors/main/warning_background' to the response with the '200' status
-  added the required property '/allOf[subschema #2]/data/items/colors/table/hover_row' to the response with the '200' status
-  added the required property '/allOf[subschema #2]/data/items/colors/table/shift_row' to the response with the '200' status


## POST /color-themes
- :warning: added the new required request property 'colors/main/error_background'
- :warning: added the new required request property 'colors/main/info_background'
- :warning: added the new required request property 'colors/main/success_background'
- :warning: added the new required request property 'colors/main/warning_background'
- :warning: added the new required request property 'colors/table/hover_row'
- :warning: added the new required request property 'colors/table/shift_row'
-  added the required property 'colors/main/error_background' to the response with the '201' status
-  added the required property 'colors/main/info_background' to the response with the '201' status
-  added the required property 'colors/main/success_background' to the response with the '201' status
-  added the required property 'colors/main/warning_background' to the response with the '201' status
-  added the required property 'colors/table/hover_row' to the response with the '201' status
-  added the required property 'colors/table/shift_row' to the response with the '201' status


## GET /color-themes/{id}
-  added the required property 'colors/main/error_background' to the response with the '200' status
-  added the required property 'colors/main/info_background' to the response with the '200' status
-  added the required property 'colors/main/success_background' to the response with the '200' status
-  added the required property 'colors/main/warning_background' to the response with the '200' status
-  added the required property 'colors/table/hover_row' to the response with the '200' status
-  added the required property 'colors/table/shift_row' to the response with the '200' status


## PUT /color-themes/{id}
- :warning: added the new required request property 'colors/main/error_background'
- :warning: added the new required request property 'colors/main/info_background'
- :warning: added the new required request property 'colors/main/success_background'
- :warning: added the new required request property 'colors/main/warning_background'
- :warning: added the new required request property 'colors/table/hover_row'
- :warning: added the new required request property 'colors/table/shift_row'
-  added the required property 'colors/main/error_background' to the response with the '200' status
-  added the required property 'colors/main/info_background' to the response with the '200' status
-  added the required property 'colors/main/success_background' to the response with the '200' status
-  added the required property 'colors/main/warning_background' to the response with the '200' status
-  added the required property 'colors/table/hover_row' to the response with the '200' status
-  added the required property 'colors/table/shift_row' to the response with the '200' status


## GET /component-alarms
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/initial_state' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/max_state' to the response with the '200' status


## PUT /contextgraph-import
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property '/items/entity_pattern/items/items/alias'


## PUT /contextgraph-import-partial
- :warning: the '/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property '/items/entity_pattern/items/items/alias'


## GET /data-storage
-  added the optional property 'config/entity' to the response with the '200' status
-  added the optional property 'config/entity_infos_log' to the response with the '200' status


## PUT /data-storage
-  added the new optional request property 'entity'
-  added the new optional request property 'entity_infos_log'
-  added the optional property 'config/entity' to the response with the '200' status
-  added the optional property 'config/entity_infos_log' to the response with the '200' status


## GET /entities
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## POST /entities/check-state-setting
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'inherited_entity_pattern/items/items/alias' to the response with the '200' status


## GET /entities/pbehaviors
- :warning: the '/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/items/pattern_exec_at' to the response with the '200' status
-  added the optional property '/items/pattern_ms' to the response with the '200' status


## GET /entities/state-setting
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'inherited_entity_pattern/items/items/alias' to the response with the '200' status


## POST /entity-export
-  added the new optional request property 'time_format'


## GET /entity-infos-dictionary/keys
-  added the optional property '/allOf[subschema #2]/data/items/type' to the response with the '200' status


## GET /entity-infos-dictionary/values
-  added the optional property '/allOf[subschema #2]/data/items/type' to the response with the '200' status


## GET /entity-infos-properties
-  endpoint added


## POST /entity-infos-properties
-  endpoint added


## DELETE /entity-infos-properties/{id}
-  endpoint added


## GET /entity-infos-properties/{id}
-  endpoint added


## PUT /entity-infos-properties/{id}
-  endpoint added


## GET /entitybasics
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## PUT /entitybasics
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## GET /entityservice-alarms/{id}
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/initial_state' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/max_state' to the response with the '200' status


## GET /entityservice-dependencies
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## GET /entityservice-impacts
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## POST /entityservice-template-validate
-  endpoint added


## GET /entityservice-template-vars
-  endpoint added


## POST /entityservices
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status


## GET /entityservices/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## PUT /entityservices/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## GET /eventfilter-copy-vars
-  endpoint added


## POST /eventfilter-template-validate
-  endpoint added


## GET /eventfilter-template-vars
-  endpoint added


## GET /eventfilter/rules
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/columns' from the response with the '200' status
-  added the new optional 'query' request parameter 'only_unread_failure'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/event_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/external_data/items/table/column_configs' to the response with the '200' status


## POST /eventfilter/rules
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the request property 'external_data/items/table/column_types'
- :warning: removed the request property 'external_data/items/table/columns'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '201' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '201' status
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'event_pattern/items/items/alias'
-  added the new optional request property 'external_data/items/table/column_configs'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'event_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'external_data/items/table/column_configs' to the response with the '201' status


## GET /eventfilter/rules/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'event_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'external_data/items/table/column_configs' to the response with the '200' status


## PUT /eventfilter/rules/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'event_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'event_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the request property 'external_data/items/table/column_types'
- :warning: removed the request property 'external_data/items/table/columns'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'event_pattern/items/items/alias'
-  added the new optional request property 'external_data/items/table/column_configs'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'event_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'external_data/items/table/column_configs' to the response with the '200' status


## POST /external-data-import/{id}
- :warning: removed the optional property 'columns' from the response with the '200' status
-  added the optional property 'column_configs' to the response with the '200' status
-  added the optional property 'error_info' to the response with the '200' status
-  added the optional property 'fail_reason' to the response with the '200' status


## PUT /external-data-import/{id}/complete
- :warning: added the new required request property 'column_tags'
- :warning: removed the request property 'column_types'


## GET /external-data-import/{id}/status
- :warning: removed the optional property 'columns' from the response with the '200' status
-  added the optional property 'column_configs' to the response with the '200' status
-  added the optional property 'error_info' to the response with the '200' status
-  added the optional property 'fail_reason' to the response with the '200' status


## GET /external-data-tables
- :warning: removed the optional property '/allOf[subschema #2]/data/items/column_types' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/columns' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/column_configs' to the response with the '200' status


## POST /external-data-tables
- :warning: removed the optional property 'column_types' from the response with the '201' status
- :warning: removed the optional property 'columns' from the response with the '201' status
-  added the optional property 'column_configs' to the response with the '201' status


## GET /external-data-tables/{id}
- :warning: removed the optional property 'column_types' from the response with the '200' status
- :warning: removed the optional property 'columns' from the response with the '200' status
-  added the optional property 'column_configs' to the response with the '200' status


## PUT /external-data-tables/{id}
- :warning: added the new required request property 'column_tags'
- :warning: removed the request property 'column_types'
- :warning: removed the optional property 'column_types' from the response with the '200' status
- :warning: removed the optional property 'columns' from the response with the '200' status
-  added the optional property 'column_configs' to the response with the '200' status


## GET /flapping-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## POST /flapping-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status


## GET /flapping-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## PUT /flapping-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## GET /idle-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## POST /idle-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status


## GET /idle-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## PUT /idle-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## POST /link-rule-template-validate
-  endpoint added


## GET /link-rule-template-vars
-  endpoint added


## GET /link-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property '/allOf[subschema #2]/data/items/external_data/items/table/columns' from the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/external_data/items/table/column_configs' to the response with the '200' status


## POST /link-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '201' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '201' status
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'external_data/items/table/column_configs' to the response with the '201' status


## GET /link-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'external_data/items/table/column_configs' to the response with the '200' status


## PUT /link-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: removed the optional property 'external_data/items/table/column_types' from the response with the '200' status
- :warning: removed the optional property 'external_data/items/table/columns' from the response with the '200' status
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'external_data/items/table/column_configs' to the response with the '200' status


## GET /notification
- :warning: api path removed without deprecation


## PUT /notification
- :warning: api path removed without deprecation


## GET /notification-settings
-  endpoint added


## PUT /notification-settings
-  endpoint added


## GET /notifications
-  endpoint added


## GET /open-alarms
- :warning: the 'entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'v/initial_state' to the response with the '200' status
-  added the optional property 'v/max_state' to the response with the '200' status


## GET /patterns
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /patterns
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the new optional request property 'weather_service_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'weather_service_pattern/items/items/alias' to the response with the '201' status


## POST /patterns-alarms-count
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/ms' to the response with the '200' status
-  added the optional property 'all/ms' to the response with the '200' status
-  added the optional property 'entities/ms' to the response with the '200' status
-  added the optional property 'entity_pattern/ms' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/ms' to the response with the '200' status


## POST /patterns-entities-count
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/ms' to the response with the '200' status
-  added the optional property 'all/ms' to the response with the '200' status
-  added the optional property 'entity_pattern/ms' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/ms' to the response with the '200' status


## GET /patterns/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'weather_service_pattern/items/items/alias' to the response with the '200' status


## PUT /patterns/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the new optional request property 'weather_service_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'weather_service_pattern/items/items/alias' to the response with the '200' status


## PUT /pbehavior-patterns
-  endpoint added


## GET /pbehaviors
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/pattern_exec_at' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/pattern_ms' to the response with the '200' status


## POST /pbehaviors
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'exec_pattern'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'pattern_exec_at' to the response with the '201' status
-  added the optional property 'pattern_ms' to the response with the '201' status


## GET /pbehaviors/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pattern_exec_at' to the response with the '200' status
-  added the optional property 'pattern_ms' to the response with the '200' status


## PATCH /pbehaviors/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pattern_exec_at' to the response with the '200' status
-  added the optional property 'pattern_ms' to the response with the '200' status


## PUT /pbehaviors/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'exec_pattern'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pattern_exec_at' to the response with the '200' status
-  added the optional property 'pattern_ms' to the response with the '200' status


## GET /pbehaviors/{id}/entities
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## GET /resolve-rules
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status


## POST /resolve-rules
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  the request property 'priority' became optional
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status


## GET /resolve-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## PUT /resolve-rules/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  the request property 'priority' became optional
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status


## GET /resolved-alarms
- :warning: the '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/initial_state' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/v/max_state' to the response with the '200' status


## POST /scenario-template-validate
-  endpoint added


## GET /scenario-template-vars
-  endpoint added


## GET /scenarios
- :warning: the '/allOf[subschema #2]/data/items/actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/actions/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/actions/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/actions/items/parameters/auth_token' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/actions/items/parameters/multiple_urls' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/actions/items/parameters/stop_on_fail' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/actions/items/parameters/stop_on_success' to the response with the '200' status


## POST /scenarios
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'actions/items/alarm_pattern/items/items/alias'
-  added the new optional request property 'actions/items/entity_pattern/items/items/alias'
-  added the new optional request property 'actions/items/parameters/auth_token'
-  added the new optional request property 'actions/items/parameters/multiple_urls'
-  added the new optional request property 'actions/items/parameters/stop_on_fail'
-  added the new optional request property 'actions/items/parameters/stop_on_success'
-  added the optional property 'actions/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'actions/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'actions/items/parameters/auth_token' to the response with the '201' status
-  added the optional property 'actions/items/parameters/multiple_urls' to the response with the '201' status
-  added the optional property 'actions/items/parameters/stop_on_fail' to the response with the '201' status
-  added the optional property 'actions/items/parameters/stop_on_success' to the response with the '201' status


## GET /scenarios/{id}
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'actions/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'actions/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'actions/items/parameters/auth_token' to the response with the '200' status
-  added the optional property 'actions/items/parameters/multiple_urls' to the response with the '200' status
-  added the optional property 'actions/items/parameters/stop_on_fail' to the response with the '200' status
-  added the optional property 'actions/items/parameters/stop_on_success' to the response with the '200' status


## PUT /scenarios/{id}
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'actions/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'actions/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'actions/items/alarm_pattern/items/items/alias'
-  added the new optional request property 'actions/items/entity_pattern/items/items/alias'
-  added the new optional request property 'actions/items/parameters/auth_token'
-  added the new optional request property 'actions/items/parameters/multiple_urls'
-  added the new optional request property 'actions/items/parameters/stop_on_fail'
-  added the new optional request property 'actions/items/parameters/stop_on_success'
-  added the optional property 'actions/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'actions/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'actions/items/parameters/auth_token' to the response with the '200' status
-  added the optional property 'actions/items/parameters/multiple_urls' to the response with the '200' status
-  added the optional property 'actions/items/parameters/stop_on_fail' to the response with the '200' status
-  added the optional property 'actions/items/parameters/stop_on_success' to the response with the '200' status


## GET /state-settings
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/inherited_entity_pattern/items/items/alias' to the response with the '200' status


## POST /state-settings
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'inherited_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'inherited_entity_pattern/items/items/alias'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'inherited_entity_pattern/items/items/alias' to the response with the '201' status


## GET /state-settings/{id}
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'inherited_entity_pattern/items/items/alias' to the response with the '200' status


## PUT /state-settings/{id}
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'inherited_entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'inherited_entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'inherited_entity_pattern/items/items/alias'
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'inherited_entity_pattern/items/items/alias' to the response with the '200' status


## GET /template-data
-  endpoint added


## POST /template-data
-  endpoint added


## DELETE /template-data/{id}
-  endpoint added


## GET /template-data/{id}
-  endpoint added


## PUT /template-data/{id}
-  endpoint added


## GET /template-test
-  endpoint added


## POST /template-test
-  endpoint added


## DELETE /template-test/{id}
-  endpoint added


## GET /template-test/{id}
-  endpoint added


## PUT /template-test/{id}
-  endpoint added


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
-  added the optional property 'filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## GET /user-preferences/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /view-copy/{id}
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the optional property 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '201' status


## POST /view-export
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'groups/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'groups/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'groups/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'groups/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'groups/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## GET /view-groups
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /view-groups
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '201' status


## GET /view-groups/{id}
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## PUT /view-groups/{id}
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /view-import
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the '/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
-  added the new optional request property '/items/views/items/tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias'
-  added the new optional request property '/items/views/items/tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias'
-  added the new optional request property '/items/views/items/tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias'
-  added the new optional request property '/items/views/items/tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias'


## POST /view-tab-copy/{id}
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the optional property 'widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '201' status


## POST /view-tabs
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the optional property 'widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '201' status


## GET /view-tabs/{id}
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## PUT /view-tabs/{id}
- :warning: the 'widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /views
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the optional property 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '201' status


## GET /views/{id}
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## PUT /views/{id}
- :warning: the 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'tabs/items/widgets/items/filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'tabs/items/widgets/items/filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'tabs/items/widgets/items/filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'tabs/items/widgets/items/filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /widget-copy/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'filters/items/alarm_pattern/items/items/alias'
-  added the new optional request property 'filters/items/entity_pattern/items/items/alias'
-  added the new optional request property 'filters/items/pbehavior_pattern/items/items/alias'
-  added the new optional request property 'filters/items/weather_service_pattern/items/items/alias'
-  added the optional property 'filters/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'filters/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'filters/items/pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'filters/items/weather_service_pattern/items/items/alias' to the response with the '201' status


## GET /widget-filters
- :warning: the '/allOf[subschema #2]/data/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the '/allOf[subschema #2]/data/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property '/allOf[subschema #2]/data/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /widget-filters
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the new optional request property 'weather_service_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'weather_service_pattern/items/items/alias' to the response with the '201' status


## GET /widget-filters/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'weather_service_pattern/items/items/alias' to the response with the '200' status


## PUT /widget-filters/{id}
- :warning: the 'alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'alarm_pattern/items/items/alias'
-  added the new optional request property 'entity_pattern/items/items/alias'
-  added the new optional request property 'pbehavior_pattern/items/items/alias'
-  added the new optional request property 'weather_service_pattern/items/items/alias'
-  added the optional property 'alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'weather_service_pattern/items/items/alias' to the response with the '200' status


## POST /widget-template-validate
-  endpoint added


## GET /widget-template-vars
-  endpoint added


## GET /widget-templates
-  added the optional property '/allOf[subschema #2]/data/items/actions' to the response with the '200' status
-  added the optional property '/allOf[subschema #2]/data/items/columns/items/isFilter' to the response with the '200' status


## POST /widget-templates
-  added the new optional request property 'actions'
-  added the new optional request property 'columns/items/isFilter'
-  added the new 'alarm_mass_quick_actions' enum value to the request property 'type'
-  added the new 'alarm_quick_actions' enum value to the request property 'type'
-  added the optional property 'actions' to the response with the '201' status
-  added the optional property 'columns/items/isFilter' to the response with the '201' status


## GET /widget-templates/{id}
-  added the optional property 'actions' to the response with the '200' status
-  added the optional property 'columns/items/isFilter' to the response with the '200' status


## PUT /widget-templates/{id}
-  added the new optional request property 'actions'
-  added the new optional request property 'columns/items/isFilter'
-  added the new 'alarm_mass_quick_actions' enum value to the request property 'type'
-  added the new 'alarm_quick_actions' enum value to the request property 'type'
-  added the optional property 'actions' to the response with the '200' status
-  added the optional property 'columns/items/isFilter' to the response with the '200' status


## POST /widgets
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '201'
-  added the new optional request property 'filters/items/alarm_pattern/items/items/alias'
-  added the new optional request property 'filters/items/entity_pattern/items/items/alias'
-  added the new optional request property 'filters/items/pbehavior_pattern/items/items/alias'
-  added the new optional request property 'filters/items/weather_service_pattern/items/items/alias'
-  added the optional property 'filters/items/alarm_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'filters/items/entity_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'filters/items/pbehavior_pattern/items/items/alias' to the response with the '201' status
-  added the optional property 'filters/items/weather_service_pattern/items/items/alias' to the response with the '201' status


## GET /widgets/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the optional property 'filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status


## PUT /widgets/{id}
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' request property type/format changed from ''/'' to 'object'/''
- :warning: the 'filters/items/alarm_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/entity_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/pbehavior_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
- :warning: the 'filters/items/weather_service_pattern/items/items/cond/value' response's property type/format changed from ''/'' to 'object'/'' for status '200'
-  added the new optional request property 'filters/items/alarm_pattern/items/items/alias'
-  added the new optional request property 'filters/items/entity_pattern/items/items/alias'
-  added the new optional request property 'filters/items/pbehavior_pattern/items/items/alias'
-  added the new optional request property 'filters/items/weather_service_pattern/items/items/alias'
-  added the optional property 'filters/items/alarm_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/entity_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/pbehavior_pattern/items/items/alias' to the response with the '200' status
-  added the optional property 'filters/items/weather_service_pattern/items/items/alias' to the response with the '200' status



