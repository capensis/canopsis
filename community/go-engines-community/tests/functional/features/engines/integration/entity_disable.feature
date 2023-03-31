Feature: entity disable 

  @concurrent
  Scenario: given disable entity request should resolve alarm
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-1",
      "connector_name": "test-connector-name-entity-disable-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-1",
      "resource": "test-resource-entity-disable-1",
      "state": 1,
      "output": "test-output-entity-disable-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-entity-disable-1"
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
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-1/test-component-entity-disable-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-entity-disable-1/test-component-entity-disable-1:
    """json
    {
      "enabled": false,
      "impact_level": 3,
      "sli_avail_state": 1
    }
    """
    Then the response code should be 200
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "component": "test-component-entity-disable-1",
      "resource": "test-resource-entity-disable-1"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-1/test-component-entity-disable-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-1&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given disable component request should resolve alarms for its resources
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-2",
      "connector_name": "test-connector-name-entity-disable-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-2",
      "resource": "test-resource-entity-disable-2-1",
      "state": 1,
      "output": "test-output-entity-disable-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-2",
      "connector_name": "test-connector-name-entity-disable-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-2",
      "resource": "test-resource-entity-disable-2-2",
      "state": 1,
      "output": "test-output-entity-disable-2"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-2-1/test-component-entity-disable-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-2-2/test-component-entity-disable-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-2-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-entity-disable-2-1"
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
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-2-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-entity-disable-2-2"
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
    When I do PUT /api/v4/entitybasics?_id=test-component-entity-disable-2:
    """json
    {
      "enabled": false,
      "impact_level": 3,
      "sli_avail_state": 1
    }
    """
    Then the response code should be 200
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-2",
        "source_type": "component"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-2",
        "resource": "test-resource-entity-disable-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-2",
        "resource": "test-resource-entity-disable-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-component-entity-disable-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-2-1/test-component-entity-disable-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-2-2/test-component-entity-disable-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-2-1&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-2-2&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    
  @concurrent
  Scenario: given disable and then enable component requests shouldn't enable disabled resources
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-3",
      "connector_name": "test-connector-name-entity-disable-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-3",
      "resource": "test-resource-entity-disable-3-1",
      "state": 0,
      "output": "test-output-entity-disable-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-3",
      "connector_name": "test-connector-name-entity-disable-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-3",
      "resource": "test-resource-entity-disable-3-2",
      "state": 0,
      "output": "test-output-entity-disable-3"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-3-1/test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-3-2/test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-component-entity-disable-3:
    """json
    {
      "enabled": false,
      "impact_level": 3,
      "sli_avail_state": 1
    }
    """
    Then the response code should be 200
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-3",
        "source_type": "component"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-3",
        "resource": "test-resource-entity-disable-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-3",
        "resource": "test-resource-entity-disable-3-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-3-1/test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-3-2/test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-component-entity-disable-3:
    """json
    {
      "enabled": true,
      "impact_level": 3,
      "sli_avail_state": 1
    }
    """
    Then the response code should be 200
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-3",
        "source_type": "component"
      }
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-3-1/test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-3-2/test-component-entity-disable-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """

  @concurrent
  Scenario: given bulk disable component request should resolve alarms for its resources
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-4",
      "connector_name": "test-connector-name-entity-disable-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-4-1",
      "resource": "test-resource-entity-disable-4-1",
      "state": 1,
      "output": "test-output-entity-disable-4"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-4",
      "connector_name": "test-connector-name-entity-disable-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-4-1",
      "resource": "test-resource-entity-disable-4-2",
      "state": 1,
      "output": "test-output-entity-disable-4"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-4",
      "connector_name": "test-connector-name-entity-disable-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-4-2",
      "resource": "test-resource-entity-disable-4-3",
      "state": 1,
      "output": "test-output-entity-disable-4"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-4",
      "connector_name": "test-connector-name-entity-disable-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-4-2",
      "resource": "test-resource-entity-disable-4-4",
      "state": 1,
      "output": "test-output-entity-disable-4"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-1/test-component-entity-disable-4-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-2/test-component-entity-disable-4-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-3/test-component-entity-disable-4-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-4/test-component-entity-disable-4-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-entity-disable-4-1"
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
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-entity-disable-4-2"
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
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-3&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-entity-disable-4-3"
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
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-4&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-entity-disable-4-4"
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-component-entity-disable-4-1"
      },
      {
        "_id": "test-component-entity-disable-4-2"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-4-1",
        "source_type": "component"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-4-2",
        "source_type": "component"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-4-1",
        "resource": "test-resource-entity-disable-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-4-1",
        "resource": "test-resource-entity-disable-4-2",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-4-2",
        "resource": "test-resource-entity-disable-4-3",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-4-2",
        "resource": "test-resource-entity-disable-4-4",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-1/test-component-entity-disable-4-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-2/test-component-entity-disable-4-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-3/test-component-entity-disable-4-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-4-4/test-component-entity-disable-4-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-1&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-2&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-3&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-entity-disable-4-4&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given bulk disable and then bulk enable component requests shouldn't enable disabled resources
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-5",
      "connector_name": "test-connector-name-entity-disable-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-5-1",
      "resource": "test-resource-entity-disable-5-1",
      "state": 0,
      "output": "test-output-entity-disable-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-5",
      "connector_name": "test-connector-name-entity-disable-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-5-1",
      "resource": "test-resource-entity-disable-5-2",
      "state": 0,
      "output": "test-output-entity-disable-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-5",
      "connector_name": "test-connector-name-entity-disable-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-5-2",
      "resource": "test-resource-entity-disable-5-3",
      "state": 0,
      "output": "test-output-entity-disable-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-entity-disable-5",
      "connector_name": "test-connector-name-entity-disable-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-entity-disable-5-2",
      "resource": "test-resource-entity-disable-5-4",
      "state": 0,
      "output": "test-output-entity-disable-5"
    }
    """
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-component-entity-disable-5-1"
      },
      {
        "_id": "test-component-entity-disable-5-2"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-1",
        "source_type": "component"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-2",
        "source_type": "component"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-1",
        "resource": "test-resource-entity-disable-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-1",
        "resource": "test-resource-entity-disable-5-2",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-2",
        "resource": "test-resource-entity-disable-5-3",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-2",
        "resource": "test-resource-entity-disable-5-4",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-component-entity-disable-5-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-component-entity-disable-5-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-1/test-component-entity-disable-5-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-2/test-component-entity-disable-5-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-3/test-component-entity-disable-5-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-4/test-component-entity-disable-5-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-component-entity-disable-5-1"
      },
      {
        "_id": "test-component-entity-disable-5-2"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-1",
        "source_type": "component"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-entity-disable-5-2",
        "source_type": "component"
      }
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-component-entity-disable-5-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-component-entity-disable-5-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-1/test-component-entity-disable-5-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-2/test-component-entity-disable-5-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-3/test-component-entity-disable-5-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-entity-disable-5-4/test-component-entity-disable-5-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": false
    }
    """
