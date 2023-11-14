Feature: resolve alarm on resolve rule
  I need to be able to resolve alarm on resolve rule

  @concurrent
  Scenario: given resolve rule should resolve alarm
    Given I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-axe-resolverule-1",
      "name": "test-resolve-rule-axe-resolverule-1-name",
      "description": "test-resolve-rule-axe-resolverule-1-desc",
      "entity_pattern":[
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-resolverule-1"
            }
          }
        ]
      ],
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "priority": 1
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-resolverule-1",
      "connector_name" : "test-connector-name-axe-resolverule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-resolverule-1",
      "resource" : "test-resource-axe-resolverule-1",
      "state" : 2,
      "output" : "test-output-axe-resolverule-1"
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-resolverule-1",
      "connector_name" : "test-connector-name-axe-resolverule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-resolverule-1",
      "resource" : "test-resource-axe-resolverule-1",
      "state" : 0,
      "output" : "test-output-axe-resolverule-1"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "resolve_close",
      "connector" : "test-connector-axe-resolverule-1",
      "connector_name" : "test-connector-name-axe-resolverule-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-resolverule-1",
      "resource" : "test-resource-axe-resolverule-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-resolverule-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-resolverule-1",
            "connector": "test-connector-axe-resolverule-1",
            "connector_name": "test-connector-name-axe-resolverule-1",
            "resource": "test-resource-axe-resolverule-1",
            "state": {
              "val": 0
            },
            "status": {
              "val": 0
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
    Then the response key "data.0.v.resolved" should be greater than 0
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """
