Feature: get alarms 
  I need to be able to get alarms

  @concurrent
  Scenario: given alarm should get updated alarm from websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-axe-alarm-websocket-1",
      "connector_name": "test-connector-name-axe-alarm-websocket-1",
      "component": "test-component-axe-alarm-websocket-1",
      "resource": "test-resource-axe-alarm-websocket-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-alarm-websocket-1
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "alarms/test-widget":
    """json
    [
      "{{ .alarmId }}"
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-axe-alarm-websocket-1",
      "connector_name": "test-connector-name-axe-alarm-websocket-1",
      "component": "test-component-axe-alarm-websocket-1",
      "resource": "test-resource-axe-alarm-websocket-1",
      "source_type": "resource"
    }
    """
    Then I wait message from websocket room "alarms/test-widget" which contains:
    """json
    {
      "_id": "{{ .alarmId }}",
      "v": {
        "state": {
          "val": 3
        }
      }
    }
    """

  @concurrent
  Scenario: given alarm should get updated alarm details from websocket room
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-axe-alarm-websocket-2",
      "connector_name": "test-connector-name-axe-alarm-websocket-2",
      "component": "test-component-axe-alarm-websocket-2",
      "resource": "test-resource-axe-alarm-websocket-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-alarm-websocket-2
    Then the response code should be 200
    When I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I connect to websocket
    When I authenticate in websocket
    When I subscribe to websocket room "alarm-details/test-widget":
    """json
    [
      {
        "_id": "{{ .alarmId }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-axe-alarm-websocket-2",
      "connector_name": "test-connector-name-axe-alarm-websocket-2",
      "component": "test-component-axe-alarm-websocket-2",
      "resource": "test-resource-axe-alarm-websocket-2",
      "source_type": "resource"
    }
    """
    Then I wait message from websocket room "alarm-details/test-widget" which contains:
    """json
    {
      "_id": "{{ .alarmId }}",
      "steps": {
        "data": [
          {
            "_t": "stateinc",
            "val": 2
          },
          {
            "_t": "statusinc",
            "val": 1
          },
          {
            "_t": "stateinc",
            "val": 3
          }
        ]
      }
    }
    """
