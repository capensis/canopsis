Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  Scenario: given pbehavior should create alarm with pbeahvior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "name": "test_pbehavior_1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test_post_resource_pbehavior_1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
      {
        "connector" : "test_post_connector_pbehavior_1",
        "connector_name" : "test_post_connector_name_pbehavior_1",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test_post_component_pbehavior_1",
        "resource" : "test_post_resource_pbehavior_1",
        "state" : 1,
        "output" : "noveo alarm"
      }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_post_resource_pbehavior_1"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test_pbehavior_1"
            },
            "connector" : "test_post_connector_pbehavior_1",
            "connector_name" : "test_post_connector_name_pbehavior_1",
            "component" : "test_post_component_pbehavior_1",
            "resource" : "test_post_resource_pbehavior_1"
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

  Scenario: given pbehavior and alarm should update alarm pbeahvior info
    Given I am admin
    When I send an event:
    """
      {
        "connector" : "test_post_connector_pbehavior_2",
        "connector_name" : "test_post_connector_name_pbehavior_2",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test_post_component_pbehavior_2",
        "resource" : "test_post_resource_pbehavior_2",
        "state" : 1,
        "output" : "noveo alarm"
      }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "name": "test_pbehavior_2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test_post_resource_pbehavior_2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test_post_resource_pbehavior_2"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test_pbehavior_2"
            },
            "connector" : "test_post_connector_pbehavior_2",
            "connector_name" : "test_post_connector_name_pbehavior_2",
            "component" : "test_post_component_pbehavior_2",
            "resource" : "test_post_resource_pbehavior_2"
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
