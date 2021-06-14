Feature: Export alarms
  I need to be able to export a alarms

  Scenario: given export request should return alarms
    When I am admin
    When I do POST /api/v4/alarm-export?search=test-alarm-to-export&active_columns[]=_id&active_columns[]=entity._id&active_columns[]=v.state.val&active_columns[]=v.status.val&active_columns[]=t&active_columns[]=v.ack.t&active_columns[]=entity.infos&active_columns[]=entity.infos.datecustom.value&time_format=DD%20MMM%20YYYY%20hh:mm:ss%20ZZ
    Then the response code should be 200
    When I do GET /api/v4/alarm-export/{{ .lastResponse._id }}
    Then the response code should be 200
    When I wait 1s
    When I do GET /api/v4/alarm-export/{{ .lastResponse._id }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    _id,entity._id,v.state.val,v.status.val,t,v.ack.t,entity.infos,entity.infos.datecustom.value
    test-alarm-to-export-2,test-entity-to-export-resource-2/test-entity-to-export-component,minor,ongoing,10 Aug 2020 05:30 CEST,,{},
    test-alarm-to-export-1,test-entity-to-export-resource-1/test-entity-to-export-component,critical,ongoing,09 Aug 2020 05:12 CEST,10 Aug 2020 06:17 CEST,"{""test-entity-to-export-resource-1-info-1"":{""name"":""test-entity-to-export-resource-1-info-1-name"",""description"":""test-entity-to-export-resource-1-info-1-description"",""value"":""test-entity-to-export-resource-1-info-1-value""}}",

    """

  Scenario: given export request should return empty response
    When I am admin
    When I do POST /api/v4/alarm-export?search=not-found&active_columns[]=_id&active_columns[]=entity._id&active_columns[]=v.state.val&active_columns[]=v.status.val&active_columns[]=t&active_columns[]=v.ack.t&active_columns[]=entity.infos&time_format=DD%20MMM%20YYYY%20hh:mm:ss%20ZZ
    Then the response code should be 200
    When I do GET /api/v4/alarm-export/{{ .lastResponse._id }}
    Then the response code should be 200
    When I wait 1s
    When I do GET /api/v4/alarm-export/{{ .lastResponse._id }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    _id,entity._id,v.state.val,v.status.val,t,v.ack.t,entity.infos

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