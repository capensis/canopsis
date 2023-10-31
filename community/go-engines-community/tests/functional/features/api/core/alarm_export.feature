Feature: Export alarms
  I need to be able to export a alarms

  Scenario: given export request should return alarms
    When I am admin
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "test-resource-to-alarm-export",
      "fields": [
        {"name": "_id", "label": "ID"},
        {"name": "entity._id", "label": "Entity"},
        {"name": "entity.notexist", "label": "Not found field"},
        {"name": "v.connector", "label": "Connector"},
        {"name": "v.connector_name", "label": "Connector name"},
        {"name": "v.component", "label": "Component"},
        {"name": "v.resource", "label": "Resource"},
        {"name": "v.state.val", "label": "State"},
        {"name": "v.status.val", "label": "Status"},
        {"name": "t", "label": "Date"},
        {"name": "v.ack.t", "label": "Ack date"},
        {"name": "v.infos", "label": "Infos"},
        {"name": "v.infos.test-dynamic-infos-to-alarm_export.test-dynamic-infos-to-alarm-export-info-name", "label": "Info 1"},
        {"name": "v.infos.notexist.notexist", "label": "Not found info"},
        {"name": "entity.infos", "label": "Entity infos"},
        {"name": "entity.infos.test-resource-to-alarm-export-1-info-1.value", "label": "Entity info 1"},
        {"name": "entity.infos.notexist.value", "label": "Not found entity info"}
      ],
      "time_format": "DD MMM YYYY hh:mm:ss ZZ"
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/alarm-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID,Entity,Not found field,Connector,Connector name,Component,Resource,State,Status,Date,Ack date,Infos,Info 1,Not found info,Entity infos,Entity info 1,Not found entity info
    test-alarm-to-export-2,test-resource-to-alarm-export-2/test-component-default,,test-connector-default,test-connector-default-name,test-component-default,test-resource-to-alarm-export-2,minor,ongoing,10 Aug 2020 05:30 CEST,,,,,,,
    test-alarm-to-export-1,test-resource-to-alarm-export-1/test-component-default,,test-connector-default,test-connector-default-name,test-component-default,test-resource-to-alarm-export-1,critical,ongoing,09 Aug 2020 05:12 CEST,10 Aug 2020 06:17 CEST,test-dynamic-infos-to-alarm-export-info-name: test-dynamic-infos-to-alarm-export-info-value,test-dynamic-infos-to-alarm-export-info-value,,test-resource-to-alarm-export-1-info-1: test-resource-to-alarm-export-1-info-1-value,test-resource-to-alarm-export-1-info-1-value,

    """

  Scenario: given export request should return empty response
    When I am admin
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "not found",
      "fields": [
        {"name": "_id", "label": "ID"},
        {"name": "entity._id", "label": "Entity"},
        {"name": "entity.notexist", "label": "Not found field"},
        {"name": "v.state.val", "label": "State"},
        {"name": "v.status.val", "label": "Status"},
        {"name": "t", "label": "Date"},
        {"name": "v.ack.t", "label": "Ack date"},
        {"name": "entity.infos", "label": "Infos"},
        {"name": "entity.infos.datecustom.value", "label": "Not found infos"},
        {"name": "links.notexist", "label": "Not found links"}
      ],
      "time_format": "DD MMM YYYY hh:mm:ss ZZ"
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/alarm-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID,Entity,Not found field,State,Status,Date,Ack date,Infos,Not found infos,Not found links

    """

  Scenario: given export request with links should return alarms with links
    When I am admin
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "test-resource-to-alarm-export",
      "fields": [
        {"name": "_id", "label": "ID"},
        {"name": "links", "label": "Links"},
        {"name": "links.test-category-to-alarm-export-1", "label": "Link 1"},
        {"name": "links.notexist", "label": "Not found links"}
      ],
      "time_format": "DD MMM YYYY hh:mm:ss ZZ"
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/alarm-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID,Links,Link 1,Not found links
    test-alarm-to-export-2,,,
    test-alarm-to-export-1,test-link-rule-to-alarm-export-link-label: http://test-link-rule-to-alarm-export-link-url.com?user=root&ids[]=test-resource-to-alarm-export-1/test-component-default&,test-link-rule-to-alarm-export-link-label: http://test-link-rule-to-alarm-export-link-url.com?user=root&ids[]=test-resource-to-alarm-export-1/test-component-default&,

    """

  Scenario: given export request with instructions should return alarms with instructions
    When I am admin
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "test-resource-to-alarm-export",
      "fields": [
        {"name": "_id", "label": "ID"},
        {"name": "assigned_instructions", "label": "Instructions"}
      ]
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/alarm-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID,Instructions
    test-alarm-to-export-2,"test-instruction-to-alarm-export-1-name,test-instruction-to-alarm-export-2-name"
    test-alarm-to-export-1,test-instruction-to-alarm-export-2-name

    """

  Scenario: given export request with templates should return alarms with executed templates
    When I am admin
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "test-resource-to-alarm-export",
      "fields": [
        {"name": "_id", "label": "ID"},
        {
          "name": "v.infos",
          "label": "Template field",
          "template": "{{ `{{ range .Alarm.Value.Infos }}{{ range . }}{{ . }}{{ end }};{{ end }}` }}"
        }
      ]
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/alarm-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID,Template field
    test-alarm-to-export-2,
    test-alarm-to-export-1,test-dynamic-infos-to-alarm-export-info-value;

    """

  Scenario: given not exit export task should return not found error
    When I am admin
    When I do GET /api/v4/alarm-export/not-found
    Then the response code should be 404
    When I do GET /api/v4/alarm-export/not-found/download
    Then the response code should be 404

  Scenario: given unauth request should not allow access
    When I do POST /api/v4/alarm-export
    Then the response code should be 401
    When I do GET /api/v4/alarm-export/not-found
    Then the response code should be 401
    When I do GET /api/v4/alarm-export/not-found/download
    Then the response code should be 401

  Scenario: given request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/alarm-export
    Then the response code should be 403
    When I do GET /api/v4/alarm-export/not-found
    Then the response code should be 403
    When I do GET /api/v4/alarm-export/not-found/download
    Then the response code should be 403
