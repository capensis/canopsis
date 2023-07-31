Feature: get alarms 
  I need to be able to get alarms

  @concurrent
  Scenario: given meta alarm and updated child should get updated details from websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "connector": "test-connector-correlation-alarm-websocket-1",
        "connector_name": "test-connector-name-correlation-alarm-websocket-1",
        "component": "test-component-correlation-alarm-websocket-1",
        "resource": "test-resource-correlation-alarm-websocket-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "connector": "test-connector-correlation-alarm-websocket-1",
        "connector_name": "test-connector-name-correlation-alarm-websocket-1",
        "component": "test-component-correlation-alarm-websocket-1",
        "resource": "test-resource-correlation-alarm-websocket-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-websocket-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-correlation-alarm-websocket-1",
      "comment": "test-metaalarm-correlation-alarm-websocket-1-comment",
      "alarms": ["{{ .alarmId1 }}", "{{ .alarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-correlation-alarm-websocket-1 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-correlation-alarm-websocket-1"
      }
    ]
    """
    When I save response metaAlarmId={{ (index .lastResponse 0)._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "alarm-details/test-widget":
    """json
    [
      {
        "_id": "{{ .metaAlarmId }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "comment",
      "output": "test-comment-correlation-alarm-websocket-1",
      "connector": "test-connector-correlation-alarm-websocket-1",
      "connector_name": "test-connector-name-correlation-alarm-websocket-1",
      "component": "test-component-correlation-alarm-websocket-1",
      "resource": "test-resource-correlation-alarm-websocket-1-1",
      "source_type": "resource"
    }
    """
    Then I wait message from websocket room "alarm-details/test-widget" which contains:
    """json
    {
      "_id": "{{ .metaAlarmId }}",
      "children": {
        "data": [
          {
            "_id": "{{ .alarmId1 }}",
            "v": {
              "last_comment": {
                "m": "test-comment-correlation-alarm-websocket-1"
              }
            }
          },
          {
            "_id": "{{ .alarmId2 }}"
          }
        ]
      }
    }
    """
