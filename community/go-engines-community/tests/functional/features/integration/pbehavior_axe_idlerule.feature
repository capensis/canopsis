Feature: update alarm on idle rule
  I need to be able to update alarm on idle rule

  Scenario: given pbehavior idle rule should update alarm
    Given I am admin
    When I do POST /api/v2/idle-rule:
    """
    {
      "name": "test-idlerule-axe-idlerule-name-3",
      "description": "test-idlerule-axe-idlerule-description-3",
      "type": "last_event",
      "duration": "3s",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-component-axe-idlerule-3",
            "state": {"val": 2}
          }
        }
      ],
      "operation": {
        "type": "pbehavior",
        "parameters": {
          "name": "test-pbehavior-axe-idlerule-3",
          "author": "test-author-axe-idlerule-3",
          "start_on_trigger": true,
          "duration": {
            "seconds": 600,
            "unit": "m"
          },
          "type": "test-maintenance-type-to-engine",
          "reason": "test-reason-to-engine"
        }
      },
      "author":"test-author-axe-idlerule-3"
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-idlerule-3",
      "connector_name" : "test-connector-name-axe-idlerule-3",
      "source_type" : "resource",
      "component" :  "test-component-axe-idlerule-3",
      "resource" : "test-resource-axe-idlerule-3",
      "state" : 2,
      "output" : "test-output-axe-idlerule-3",
      "long_output" : "test-long-output-axe-idlerule-3",
      "author" : "test-author-axe-idlerule-3",
      "timestamp": {{ (now.Add (parseDuration "-5s")).UTC.Unix }}
    }
    """
    When I wait the end of event processing
    When I wait 3s
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-axe-idlerule-3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-3",
            "connector": "test-connector-axe-idlerule-3",
            "connector_name": "test-connector-name-axe-idlerule-3",
            "resource": "test-resource-axe-idlerule-3",
            "pbehavior_info": {
              "canonical_type": "maintenance",
              "name": "test-pbehavior-axe-idlerule-3",
              "reason": "Test Engine",
              "type": "test-maintenance-type-to-engine",
              "type_name": "Engine maintenance"
            },
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
                "_t": "pbhenter",
                "a": "system",
                "m": "Pbehavior test-pbehavior-axe-idlerule-3. Type: Engine maintenance. Reason: Test Engine"
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
    # Change alarm state to not match idle rule condition anymore to not disturb another tests
    When I send an event:
    """
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-idlerule-3",
      "connector_name" : "test-connector-name-axe-idlerule-3",
      "source_type" : "resource",
      "component" :  "test-component-axe-idlerule-3",
      "resource" : "test-resource-axe-idlerule-3",
      "state" : 3,
      "output" : "test-output-axe-idlerule-3",
      "long_output" : "test-long-output-axe-idlerule-3",
      "author" : "test-author-axe-idlerule-3",
      "timestamp": {{ (now.Add (parseDuration "-5s")).UTC.Unix }}
    }
    """
    When I wait the end of event processing
