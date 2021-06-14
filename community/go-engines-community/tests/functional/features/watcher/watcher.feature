Feature: update watcher on event
  I need to be able to see new watcher state on event

  Scenario: given watcher for entity should update watcher on entity event
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "_id": "test-watcher",
      "name": "test-watcher",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-watcher"}
      ],
      "output_template": "test-watcher"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" : "test-connector-watcher",
      "connector_name" : "test-connector-name-watcher",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-watcher",
      "resource" : "test-resource-watcher",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.component":"test-watcher"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "connector" : "watcher",
            "connector_name" : "watcher",
            "component" :  "test-watcher"
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
