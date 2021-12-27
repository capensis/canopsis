Feature: send activation event on pbhleave
  I need to be able to trigger rule on alarm activation

  Scenario: given event for new alarm and maintenance pbehavior should not send event
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-1",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-1",
      "resource" : "test-resource-pbehavior-axe-activation-event-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-activation-event-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "10m" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-axe-activation-event-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-1",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-1",
      "resource" : "test-resource-pbehavior-axe-activation-event-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-pbehavior-axe-activation-event-1"},{"v.activation_date":{"$exists":false}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-axe-activation-event-1",
            "connector_name" : "test-connector-name-pbehavior-axe-activation-event-1",
            "component" : "test-component-pbehavior-axe-activation-event-1",
            "resource" : "test-resource-pbehavior-axe-activation-event-1",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "pbhenter"}
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

  Scenario: given event for new alarm and active pbehavior should send event
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-2",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-2",
      "resource" : "test-resource-pbehavior-axe-activation-event-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-activation-event-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "10m" }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-axe-activation-event-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-2",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-2",
      "resource" : "test-resource-pbehavior-axe-activation-event-2",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-pbehavior-axe-activation-event-2"},{"v.activation_date":{"$exists":true}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-axe-activation-event-2",
            "connector_name" : "test-connector-name-pbehavior-axe-activation-event-2",
            "component" : "test-component-pbehavior-axe-activation-event-2",
            "resource" : "test-resource-pbehavior-axe-activation-event-2",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "pbhenter"}
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

  Scenario: given event for new alarm and maintenance pbehavior should send event on pbhleave
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-3",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-3",
      "resource" : "test-resource-pbehavior-axe-activation-event-3",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-activation-event-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "3s" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-axe-activation-event-3"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-3",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-3",
      "resource" : "test-resource-pbehavior-axe-activation-event-3",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I wait 3s
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-pbehavior-axe-activation-event-3"},{"v.activation_date":{"$exists":true}},{"$expr":{"$ne":["$v.activation_date","$v.creation_date"]}}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-axe-activation-event-3",
            "connector_name" : "test-connector-name-pbehavior-axe-activation-event-3",
            "component" : "test-component-pbehavior-axe-activation-event-3",
            "resource" : "test-resource-pbehavior-axe-activation-event-3",
            "steps": [
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "pbhenter"},
              {"_t": "pbhleave"}
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
