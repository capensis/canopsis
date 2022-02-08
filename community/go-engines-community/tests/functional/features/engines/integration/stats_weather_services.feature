Feature: get stats with service weather
  I need to be able to get service weather

  Scenario: given service for one entity with stats in template should increment OK for entities without alarms
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-entityservice-stats-integration-connector",
      "connector_name": "test-entityservice-stats-integration-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-entityservice-stats-integration-component",
      "resource": "test-stats-ok-integration-resource-1",
      "state": 3,
      "output": "test-entityservice-stats-integration-resource-1"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-stats-integration",
      "name": "test-entityservice-stats-integration",
      "output_template": "test-entityservice-stats-integration-output",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
            "name": "test-stats-ok-integration-resource-1",
            "type" : "resource"
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing 
    When I wait the next periodical process
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"test-stats-ok-integration-resource-1/test-entityservice-stats-integration-component"}]}&correlation=false
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "entity": {
            "_id": "test-stats-ok-integration-resource-1/test-entityservice-stats-integration-component"
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
    When I do GET /api/v4/weather-services/test-entityservice-stats-integration
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-stats-ok-integration-resource-1/test-entityservice-stats-integration-component",
          "stats": {
              "ko": 1,
              "ok": 0
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
    Then the response key "data.0.stats.last_ko" should not be "null"
    When I send an event:
    """json
    {
      "connector": "test-entityservice-stats-integration-connector",
      "connector_name": "test-entityservice-stats-integration-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-entityservice-stats-integration-component",
      "resource": "test-stats-ok-integration-resource-1",
      "state": 2,
      "output": "test-entityservice-stats-integration-resource-1"
    }
    """
    When I wait the end of event processing
    When I wait the next periodical process
    When I do GET /api/v4/weather-services/test-entityservice-stats-integration
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-stats-ok-integration-resource-1/test-entityservice-stats-integration-component",
          "stats": {
              "ko": 2,
              "ok": 0
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
    Then the response key "data.0.stats.last_ko" should not be "null"
    When I do POST /api/v4/pbehaviors:
    """json
    {
        "_id": "test-entityservice-stats-integration-pbehavior-1", 
        "enabled": true,
        "author": "root",
        "name": "test-entityservice-stats-integration-pbehavior-1",
        "comments":[{"author":"root","message":"pause test-entityservice-stats-integration-pbehavior-1"}],
        "tstart": {{ now.Unix }},
        "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
        "type": "test-default-pause-type",
        "reason": "test-reason-to-engine",
        "filter": {
            "name": "test-stats-ok-integration-resource-1"
        }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do GET /api/v4/weather-services/test-entityservice-stats-integration
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-stats-ok-integration-resource-1/test-entityservice-stats-integration-component",
          "stats": {
              "ko": 0,
              "ok": 0
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
    Then the response key "data.0.stats.last_ko" should not be "null"
    Then the response key "data.0.stats.last_event" should not be "null"
    When I do DELETE /api/v4/pbehaviors/test-entityservice-stats-integration-pbehavior-1
    Then the response code should be 204
    When I wait the end of event processing
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-entityservice-stats-integration-connector",
      "connector_name": "test-entityservice-stats-integration-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-entityservice-stats-integration-component",
      "resource": "test-stats-ok-integration-resource-1",
      "state": 0,
      "output": "test-entityservice-stats-integration-resource-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services/test-entityservice-stats-integration
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-stats-ok-integration-resource-1/test-entityservice-stats-integration-component",
          "stats": {
              "ko": 0,
              "ok": 1
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
    Then the response key "data.0.stats.last_event" should not be "null"
    Then the response key "data.0.stats.last_ko" should not be "null"
