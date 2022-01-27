Feature: Export alarms
  I need to be able to export a alarms

  Scenario: given export request should return alarms
    When I am admin
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "test-alarm-to-export",
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
    test-alarm-to-export-2,test-entity-to-export-resource-2/test-entity-to-export-component,,minor,ongoing,10 Aug 2020 05:30 CEST,,{},,
    test-alarm-to-export-1,test-entity-to-export-resource-1/test-entity-to-export-component,,critical,ongoing,09 Aug 2020 05:12 CEST,10 Aug 2020 06:17 CEST,"{""test-entity-to-export-resource-1-info-1"":{""name"":""test-entity-to-export-resource-1-info-1-name"",""description"":""test-entity-to-export-resource-1-info-1-description"",""value"":""test-entity-to-export-resource-1-info-1-value""}}",,

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
