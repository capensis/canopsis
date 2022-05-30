Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event


  Scenario: given double ack events should update alarm with double ack
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-20",
      "connector_name" : "test-connector-name-axe-20",
      "source_type" : "resource",
      "component" :  "test-component-axe-20",
      "resource" : "test-resource-axe-20",
      "state" : 2,
      "output" : "test-output-axe-20",
      "long_output" : "test-long-output-axe-20",
      "author" : "test-author-axe-20"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-20",
      "connector_name" : "test-connector-name-axe-20",
      "source_type" : "resource",
      "component" :  "test-component-axe-20",
      "resource" : "test-resource-axe-20",
      "output" : "test-output-axe-20",
      "long_output" : "test-long-output-axe-20",
      "author" : "test-author-axe-20",
      "user_id": "test-author-id-20"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-20"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-20",
              "m": "test-output-axe-20",
              "user_id": "test-author-id-20",
              "val": 0
            },
            "component": "test-component-axe-20",
            "connector": "test-connector-axe-20",
            "connector_name": "test-connector-name-axe-20",
            "resource": "test-resource-axe-20",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "ack",
                "a": "test-author-axe-20",
                "m": "test-output-axe-20",
                "val": 0
              }
            ]
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
    When I send an event:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-20",
      "connector_name" : "test-connector-name-axe-20",
      "source_type" : "resource",
      "component" :  "test-component-axe-20",
      "resource" : "test-resource-axe-20",
      "output" : "new-test-output-axe-20",
      "long_output" : "test-long-output-axe-20",
      "author" : "test-author-axe-20",
      "user_id": "test-author-id-20"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-20"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-20",
              "m": "new-test-output-axe-20",
              "user_id": "test-author-id-20",
              "val": 0
            },
            "component": "test-component-axe-20",
            "connector": "test-connector-axe-20",
            "connector_name": "test-connector-name-axe-20",
            "resource": "test-resource-axe-20",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "ack",
                "a": "test-author-axe-20",
                "m": "test-output-axe-20",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-20",
                "m": "new-test-output-axe-20",
                "val": 0
              }
            ]
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
