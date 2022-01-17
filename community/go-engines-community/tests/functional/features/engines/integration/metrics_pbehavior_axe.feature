Feature: SLI metrics should be added on alarm changes
  I need to be able to see SLI metrics.

  Scenario: given entity in maintenance pbehavior should add SLI maintenance metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-1-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-1",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-1",
      "resource" : "test-resource-metrics-pbehavior-axe-1",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-pbehavior-axe-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-metrics-pbehavior-axe-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.maintenance" is greater or equal than 2
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "downtime": 0
        }
      ]
    }
    """

  Scenario: given entity in pause pbehavior should add SLI downtime metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-2"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-2",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-2",
      "resource" : "test-resource-metrics-pbehavior-axe-2",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-pbehavior-axe-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
      "type": "test-pause-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-metrics-pbehavior-axe-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.downtime" is greater or equal than 2
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "maintenance": 0
        }
      ]
    }
    """

  Scenario: given entity in inactive pbehavior should add SLI maintenance metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-3-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-3"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-3",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-3",
      "resource" : "test-resource-metrics-pbehavior-axe-3",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-pbehavior-axe-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
      "type": "test-inactive-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-metrics-pbehavior-axe-3"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.maintenance" is greater or equal than 2
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "downtime": 0
        }
      ]
    }
    """

  Scenario: given entity in active pbehavior should add SLI maintenance metrics for outer intervals
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-4-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-4"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-4",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-4",
      "resource" : "test-resource-metrics-pbehavior-axe-4",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-pbehavior-axe-4",
      "tstart": {{ nowAdd "2s" }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-metrics-pbehavior-axe-4"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.maintenance" is greater or equal than 2
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "downtime": 0
        }
      ]
    }
    """

  Scenario: given entity in active pbehavior should not add SLI any metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-5-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-5"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-5",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-5",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-5",
      "resource" : "test-resource-metrics-pbehavior-axe-5",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-pbehavior-axe-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-metrics-pbehavior-axe-5"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I wait 1s
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": []
    }
    """

  Scenario: given alarm should add downtime metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-6-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-6"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-6",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-6",
      "resource" : "test-resource-metrics-pbehavior-axe-6",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-6",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-6",
      "resource" : "test-resource-metrics-pbehavior-axe-6",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.downtime" is greater or equal than 2
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "maintenance": 0
        }
      ]
    }
    """

  Scenario: given minor alarm with SLI minor state should not add downtime metrics for minor state
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-7-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-7"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-7",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-7",
      "resource" : "test-resource-metrics-pbehavior-axe-7",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-metrics-pbehavior-axe-7/test-component-metrics-pbehavior-axe-7:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 1
    }
    """
    Then the response code should be 200
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-7",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-7",
      "resource" : "test-resource-metrics-pbehavior-axe-7",
      "state" : 2
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-7",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-7",
      "resource" : "test-resource-metrics-pbehavior-axe-7",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.downtime" is greater or equal than 2
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "maintenance": 0
        }
      ]
    }
    """
    When I save response downtime={{ (index .lastResponse.data 0).downtime }}
    When I wait 2s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-7",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-7",
      "resource" : "test-resource-metrics-pbehavior-axe-7",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "timestamp": {{ nowDate }},
          "downtime": {{ .downtime }},
          "maintenance": 0
        }
      ]
    }
    """

  Scenario: given minor alarm with SLI critical state should not add downtime metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-pbehavior-axe-8-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-pbehavior-axe-8"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-8",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-8",
      "resource" : "test-resource-metrics-pbehavior-axe-8",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-metrics-pbehavior-axe-8/test-component-metrics-pbehavior-axe-8:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 3
    }
    """
    Then the response code should be 200
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-8",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-8",
      "resource" : "test-resource-metrics-pbehavior-axe-8",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-8",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-8",
      "resource" : "test-resource-metrics-pbehavior-axe-8",
      "state" : 2
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-8",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-8",
      "resource" : "test-resource-metrics-pbehavior-axe-8",
      "state" : 3
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-pbehavior-axe-8",
      "connector_name" : "test-connector-name-metrics-pbehavior-axe-8",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-pbehavior-axe-8",
      "resource" : "test-resource-metrics-pbehavior-axe-8",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I do GET /api/v4/cat/metrics/sli?filter={{ .filterID }}&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": []
    }
    """
