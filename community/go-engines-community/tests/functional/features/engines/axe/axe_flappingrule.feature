Feature: update alarm status on flapping rule
  I need to be able to update alarm status on flapping rule

  @concurrent
  Scenario: given flapping rule should update alarm status
    Given I am admin
    When I do POST /api/v4/flapping-rules:
    """json
    {
      "_id": "test-flapping-rule-axe-flappingrule-1",
      "name": "test-flapping-rule-axe-flappingrule-1-name",
      "description": "test-flapping-rule-axe-flappingrule-1-desc",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-flappingrule-1"
            }
          }
        ]
      ],
      "freq_limit": 2,
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "priority": 10
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-flappingrule-1",
      "connector_name" : "test-connector-name-axe-flappingrule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-flappingrule-1",
      "resource" : "test-resource-axe-flappingrule-1",
      "state" : 1,
      "output" : "test-output-axe-flappingrule-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-flappingrule-1",
      "connector_name" : "test-connector-name-axe-flappingrule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-flappingrule-1",
      "resource" : "test-resource-axe-flappingrule-1",
      "state" : 2,
      "output" : "test-output-axe-flappingrule-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-flappingrule-1",
      "connector_name" : "test-connector-name-axe-flappingrule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-flappingrule-1",
      "resource" : "test-resource-axe-flappingrule-1",
      "state" : 3,
      "output" : "test-output-axe-flappingrule-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-flappingrule-1",
      "connector_name" : "test-connector-name-axe-flappingrule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-flappingrule-1",
      "resource" : "test-resource-axe-flappingrule-1",
      "state" : 2,
      "output" : "test-output-axe-flappingrule-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-flappingrule-1",
      "connector_name" : "test-connector-name-axe-flappingrule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-flappingrule-1",
      "resource" : "test-resource-axe-flappingrule-1",
      "state" : 3,
      "output" : "test-output-axe-flappingrule-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-flappingrule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-flappingrule-1",
            "connector": "test-connector-axe-flappingrule-1",
            "connector_name": "test-connector-name-axe-flappingrule-1",
            "resource": "test-resource-axe-flappingrule-1",
            "state": {
              "val": 3
            },
            "status": {
              "val": 3
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statedec",
                "val": 2
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 3
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 7
            }
          }
        }
      }
    ]
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "updatestatus",
      "connector" : "test-connector-axe-flappingrule-1",
      "connector_name" : "test-connector-name-axe-flappingrule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-flappingrule-1",
      "resource" : "test-resource-axe-flappingrule-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-flappingrule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-flappingrule-1",
            "connector": "test-connector-axe-flappingrule-1",
            "connector_name": "test-connector-name-axe-flappingrule-1",
            "resource": "test-resource-axe-flappingrule-1",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statedec",
                "val": 2
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 3
              },
              {
                "_t": "statusdec",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 8
            }
          }
        }
      }
    ]
    """
