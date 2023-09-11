Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given manual metaalarm and alarms should increase events_count and last_event_date in metaalarm
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-second-1-1"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response lastEventDate1={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-correlation-alarm-get-second-1",
      "comment": "test-metaalarm-correlation-alarm-get-second-1-comment",
      "alarms": ["{{ .alarmId1 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-correlation-alarm-get-second-1 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-correlation-alarm-get-second-1"
      }
    ]
    """
    Then I save response manualMetaAlarm={{ (index .lastResponse 0)._id }}
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "events_count": 1
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate1 }}"
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-second-1-2"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I save response alarmId2={{ (index .lastResponse.data 0)._id }}
    Then I save response lastEventDate2={{ (index .lastResponse.data 0).v.last_event_date }}
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-correlation-alarm-get-second-1",
      "connector_name": "test-correlation-alarm-get-second-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-correlation-alarm-get-second-1",
      "resource": "test-resource-correlation-alarm-get-second-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-correlation-alarm-get-second-1-3"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I save response alarmId3={{ (index .lastResponse.data 0)._id }}
    Then I save response lastEventDate3={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .manualMetaAlarm }}/add:
    """json
    {
      "comment": "test-metaalarm-correlation-alarm-get-second-1-comment",
      "alarms": ["{{ .alarmId2 }}", "{{ .alarmId3 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "events_count": 6
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate3 }}"
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .manualMetaAlarm }}/remove:
    """json
    {
      "comment": "test-metaalarm-correlation-alarm-get-second-1-comment",
      "alarms": ["{{ .alarmId2 }}", "{{ .alarmId3 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-alarm-get-second-1&correlation=true&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "events_count": 1
          }
        },
        {
        },
        {
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    Then the response key "data.0.v.last_event_date" should be "{{ .lastEventDate1 }}"
