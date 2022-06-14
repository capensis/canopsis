Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  Scenario: given new alarm should add created_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-1-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-1",
      "connector_name" : "test-connector-name-metrics-axe-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-1",
      "resource" : "test-resource-metrics-axe-1",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """

  Scenario: given new alarm with auto instruction should add instruction_alarms and non_displayed_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-2-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-2-1"
        },
        {
          "name": "test-resource-metrics-axe-2-2"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-metrics-axe-2",
        "connector_name" : "test-connector-name-metrics-axe-2",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-2",
        "resource" : "test-resource-metrics-axe-2-1",
        "state" : 1
      },
      {
        "connector" : "test-connector-metrics-axe-2",
        "connector_name" : "test-connector-name-metrics-axe-2",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-2",
        "resource" : "test-resource-metrics-axe-2-2",
        "state" : 1
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-metrics-axe-2-1"}]}&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
    """json
    [
      {
        "_t": "autoinstructioncomplete",
        "m": "Instruction test-instruction-metrics-axe-2-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "m": "Instruction test-instruction-metrics-axe-2-2-name."
      }
    ]
    """
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-metrics-axe-2-2"}]}&with_steps=true until response code is 200 and response array key "data.0.v.steps" contains:
    """json
    [
      {
        "_t": "autoinstructioncomplete",
        "m": "Instruction test-instruction-metrics-axe-2-1-name."
      },
      {
        "_t": "autoinstructioncomplete",
        "m": "Instruction test-instruction-metrics-axe-2-2-name."
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=instruction_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "instruction_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """

  Scenario: given new alarm under pbehavior should add pbehavior_alarms and non_displayed_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-3-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-3"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-3",
      "connector_name" : "test-connector-name-metrics-axe-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-3",
      "resource" : "test-resource-metrics-axe-3",
      "state" : 0
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-axe-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-metrics-axe-3"
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
      "connector" : "test-connector-metrics-axe-3",
      "connector_name" : "test-connector-name-metrics-axe-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-3",
      "resource" : "test-resource-metrics-axe-3",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=pbehavior_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "pbehavior_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """

#  todo race condition
#  Scenario: given new alarm and new meta alarm should add correlation_alarms and non_displayed_alarms metrics
#    Given I am admin
#    When I do POST /api/v4/cat/filters:
#    """json
#    {
#      "name": "test-filter-metrics-axe-4-name",
#      "entity_patterns": [
#        {
#          "name": "test-resource-metrics-axe-4"
#        }
#      ]
#    }
#    """
#    Then the response code should be 201
#    When I save response filterID={{ .lastResponse._id }}
#    When I send an event:
#    """json
#    {
#      "connector" : "test-connector-metrics-axe-4",
#      "connector_name" : "test-connector-name-metrics-axe-4",
#      "source_type" : "resource",
#      "event_type" : "check",
#      "component" : "test-component-metrics-axe-4",
#      "resource" : "test-resource-metrics-axe-4",
#      "state" : 1
#    }
#    """
#    When I wait the end of 3 events processing
#    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=correlation_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
#    """json
#    {
#      "data": [
#        {
#          "title": "correlation_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "non_displayed_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        }
#      ]
#    }
#    """
#
#  Scenario: given new alarm and existed meta alarm should add correlation_alarms and non_displayed_alarms metrics
#    Given I am admin
#    When I do POST /api/v4/cat/filters:
#    """json
#    {
#      "name": "test-filter-metrics-axe-5-name",
#      "entity_patterns": [
#        {
#          "name": "test-resource-metrics-axe-5-1"
#        },
#        {
#          "name": "test-resource-metrics-axe-5-2"
#        }
#      ]
#    }
#    """
#    Then the response code should be 201
#    When I save response filterID={{ .lastResponse._id }}
#    When I send an event:
#    """json
#    {
#      "connector" : "test-connector-metrics-axe-5",
#      "connector_name" : "test-connector-name-metrics-axe-5",
#      "source_type" : "resource",
#      "event_type" : "check",
#      "component" : "test-component-metrics-axe-5",
#      "resource" : "test-resource-metrics-axe-5-1",
#      "state" : 1
#    }
#    """
#    When I wait the end of 2 events processing
#    When I send an event:
#    """json
#    {
#      "connector" : "test-connector-metrics-axe-5",
#      "connector_name" : "test-connector-name-metrics-axe-5",
#      "source_type" : "resource",
#      "event_type" : "check",
#      "component" : "test-component-metrics-axe-5",
#      "resource" : "test-resource-metrics-axe-5-2",
#      "state" : 1
#    }
#    """
#    When I wait the end of 2 events processing
#    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=correlation_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
#    """json
#    {
#      "data": [
#        {
#          "title": "correlation_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "non_displayed_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        }
#      ]
#    }
#    """

  Scenario: given acked alarm should add ack_alarms and average_ack metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-6-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-6"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-6",
      "connector_name" : "test-connector-name-metrics-axe-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-6",
      "resource" : "test-resource-metrics-axe-6",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-6",
      "connector_name" : "test-connector-name-metrics-axe-6",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-6",
      "resource" : "test-resource-metrics-axe-6"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ack_alarms&parameters[]=average_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ nowDate }}
            }
          ]
        }
      ]
    }
    """
    Then the response key "data.1.data.0.value" should be greater or equal than 1
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-6",
      "connector_name" : "test-connector-name-metrics-axe-6",
      "source_type" : "resource",
      "event_type" : "ackremove",
      "component" : "test-component-metrics-axe-6",
      "resource" : "test-resource-metrics-axe-6"
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-6",
      "connector_name" : "test-connector-name-metrics-axe-6",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-6",
      "resource" : "test-resource-metrics-axe-6"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ack_alarms&parameters[]=average_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ nowDate }}
            }
          ]
        }
      ]
    }
    """
    Then the response key "data.1.data.0.value" should be greater or equal than 1

  Scenario: given unacked alarm should add cancel_ack_alarms and ack_active_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-7-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-7"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-7",
      "connector_name" : "test-connector-name-metrics-axe-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-7",
      "resource" : "test-resource-metrics-axe-7",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-7",
      "connector_name" : "test-connector-name-metrics-axe-7",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-7",
      "resource" : "test-resource-metrics-axe-7"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-7",
      "connector_name" : "test-connector-name-metrics-axe-7",
      "source_type" : "resource",
      "event_type" : "ackremove",
      "component" : "test-component-metrics-axe-7",
      "resource" : "test-resource-metrics-axe-7"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=cancel_ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "cancel_ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "ack_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-7",
      "connector_name" : "test-connector-name-metrics-axe-7",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-7",
      "resource" : "test-resource-metrics-axe-7"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=cancel_ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "cancel_ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "ack_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """

  Scenario: given alarm with ticket should add ticket_active_alarms and without_ticket_active_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-8-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-8-1"
        },
        {
          "name": "test-resource-metrics-axe-8-2"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-metrics-axe-8",
        "connector_name" : "test-connector-name-metrics-axe-8",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-8",
        "resource" : "test-resource-metrics-axe-8-1",
        "state" : 1
      },
      {
        "connector" : "test-connector-metrics-axe-8",
        "connector_name" : "test-connector-name-metrics-axe-8",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-8",
        "resource" : "test-resource-metrics-axe-8-2",
        "state" : 1
      }
    ]
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-8",
      "connector_name" : "test-connector-name-metrics-axe-8",
      "source_type" : "resource",
      "event_type" : "assocticket",
      "component" : "test-component-metrics-axe-8",
      "resource" : "test-resource-metrics-axe-8-1",
      "ticket": "testticket"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-8",
      "connector_name" : "test-connector-name-metrics-axe-8",
      "source_type" : "resource",
      "event_type" : "assocticket",
      "component" : "test-component-metrics-axe-8",
      "resource" : "test-resource-metrics-axe-8-1",
      "ticket": "testticket"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ticket_active_alarms&parameters[]=without_ticket_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "without_ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """

  Scenario: given resolved alarm should add average_resolve metrics
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-9-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-9"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-9",
      "connector_name" : "test-connector-name-metrics-axe-9",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-9",
      "resource" : "test-resource-metrics-axe-9",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-9",
      "connector_name" : "test-connector-name-metrics-axe-9",
      "source_type" : "resource",
      "event_type" : "cancel",
      "component" : "test-component-metrics-axe-9",
      "resource" : "test-resource-metrics-axe-9"
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-9",
      "connector_name" : "test-connector-name-metrics-axe-9",
      "source_type" : "resource",
      "event_type" : "resolve_cancel",
      "component" : "test-component-metrics-axe-9",
      "resource" : "test-resource-metrics-axe-9"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=average_resolve&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "average_resolve",
          "data": [
            {
              "timestamp": {{ nowDate }}
            }
          ]
        }
      ]
    }
    """
    Then the response key "data.0.data.0.value" should be greater or equal than 1

#  todo race condition
#  Scenario: given new alarm with auto instruction, meta alarm and pbehavior should add non_displayed_alarms metrics only once
#    Given I am admin
#    When I do POST /api/v4/cat/filters:
#    """json
#    {
#      "name": "test-filter-metrics-axe-10-name",
#      "entity_patterns": [
#        {
#          "name": "test-resource-metrics-axe-10"
#        }
#      ]
#    }
#    """
#    Then the response code should be 201
#    When I save response filterID={{ .lastResponse._id }}
#    When I send an event:
#    """json
#    {
#      "connector" : "test-connector-metrics-axe-10",
#      "connector_name" : "test-connector-name-metrics-axe-10",
#      "source_type" : "resource",
#      "event_type" : "check",
#      "component" : "test-component-metrics-axe-10",
#      "resource" : "test-resource-metrics-axe-10",
#      "state" : 0
#    }
#    """
#    When I wait the end of event processing
#    When I do POST /api/v4/pbehaviors:
#    """json
#    {
#      "enabled": true,
#      "name": "test-pbehavior-metrics-axe-10",
#      "tstart": {{ now }},
#      "tstop": {{ nowAdd "1h" }},
#      "type": "test-maintenance-type-to-engine",
#      "reason": "test-reason-to-engine",
#      "filter":{
#        "$and":[
#          {
#            "name": "test-resource-metrics-axe-10"
#          }
#        ]
#      }
#    }
#    """
#    Then the response code should be 201
#    When I wait the end of event processing
#    When I send an event:
#    """json
#    {
#      "connector" : "test-connector-metrics-axe-10",
#      "connector_name" : "test-connector-name-metrics-axe-10",
#      "source_type" : "resource",
#      "event_type" : "check",
#      "component" : "test-component-metrics-axe-10",
#      "resource" : "test-resource-metrics-axe-10",
#      "state" : 1
#    }
#    """
#    When I wait the end of 2 events processing
#    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=instruction_alarms&parameters[]=correlation_alarms&parameters[]=pbehavior_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
#    """json
#    {
#      "data": [
#        {
#          "title": "instruction_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "correlation_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "pbehavior_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "non_displayed_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        }
#      ]
#    }
#    """
#
#  Scenario: given resolved alarm should decrease active_alarms, ratio_correlation, ratio_instructions, ratio_tickets, ratio_non_displayed, ack_active_alarms metrics
#    Given I am admin
#    When I do POST /api/v4/cat/filters:
#    """json
#    {
#      "name": "test-filter-metrics-axe-11-name",
#      "entity_patterns": [
#        {
#          "name": "test-resource-metrics-axe-11-1"
#        },
#        {
#          "name": "test-resource-metrics-axe-11-2"
#        }
#      ]
#    }
#    """
#    Then the response code should be 201
#    When I save response filterID={{ .lastResponse._id }}
#    When I send an event:
#    """json
#    [
#      {
#        "connector" : "test-connector-metrics-axe-11",
#        "connector_name" : "test-connector-name-metrics-axe-11",
#        "source_type" : "resource",
#        "event_type" : "check",
#        "component" : "test-component-metrics-axe-11",
#        "resource" : "test-resource-metrics-axe-11-1",
#        "state" : 1
#      },
#      {
#        "connector" : "test-connector-metrics-axe-11",
#        "connector_name" : "test-connector-name-metrics-axe-11",
#        "source_type" : "resource",
#        "event_type" : "check",
#        "component" : "test-component-metrics-axe-11",
#        "resource" : "test-resource-metrics-axe-11-2",
#        "state" : 1
#      }
#    ]
#    """
#    When I wait the end of 4 events processing
#    When I send an event:
#    """json
#    [
#      {
#        "connector" : "test-connector-metrics-axe-11",
#        "connector_name" : "test-connector-name-metrics-axe-11",
#        "source_type" : "resource",
#        "event_type" : "assocticket",
#        "component" : "test-component-metrics-axe-11",
#        "resource" : "test-resource-metrics-axe-11-1",
#        "ticket": "testticket"
#      },
#      {
#        "connector" : "test-connector-metrics-axe-11",
#        "connector_name" : "test-connector-name-metrics-axe-11",
#        "source_type" : "resource",
#        "event_type" : "assocticket",
#        "component" : "test-component-metrics-axe-11",
#        "resource" : "test-resource-metrics-axe-11-2",
#        "ticket": "testticket"
#      },
#      {
#        "connector" : "test-connector-metrics-axe-11",
#        "connector_name" : "test-connector-name-metrics-axe-11",
#        "source_type" : "resource",
#        "event_type" : "ack",
#        "component" : "test-component-metrics-axe-11",
#        "resource" : "test-resource-metrics-axe-11-1"
#      },
#      {
#        "connector" : "test-connector-metrics-axe-11",
#        "connector_name" : "test-connector-name-metrics-axe-11",
#        "source_type" : "resource",
#        "event_type" : "ack",
#        "component" : "test-component-metrics-axe-11",
#        "resource" : "test-resource-metrics-axe-11-2"
#      }
#    ]
#    """
#    When I wait the end of 4 events processing
#    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=created_alarms&parameters[]=active_alarms&parameters[]=instruction_alarms&parameters[]=ratio_instructions&parameters[]=correlation_alarms&parameters[]=ratio_correlation&parameters[]=non_displayed_alarms&parameters[]=ratio_non_displayed&parameters[]=ticket_active_alarms&parameters[]=ratio_tickets&parameters[]=ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
#    """json
#    {
#      "data": [
#        {
#          "title": "created_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "active_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "instruction_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "ratio_instructions",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 50
#            }
#          ]
#        },
#        {
#          "title": "correlation_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "ratio_correlation",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 100
#            }
#          ]
#        },
#        {
#          "title": "non_displayed_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "ratio_non_displayed",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 100
#            }
#          ]
#        },
#        {
#          "title": "ticket_active_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "ratio_tickets",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 100
#            }
#          ]
#        },
#        {
#          "title": "ack_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "ack_active_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        }
#      ]
#    }
#    """
#    When I send an event:
#    """json
#    {
#      "connector" : "test-connector-metrics-axe-11",
#      "connector_name" : "test-connector-name-metrics-axe-11",
#      "source_type" : "resource",
#      "event_type" : "cancel",
#      "component" : "test-component-metrics-axe-11",
#      "resource" : "test-resource-metrics-axe-11-1"
#    }
#    """
#    When I wait the end of event processing
#    When I send an event:
#    """json
#    {
#      "connector" : "test-connector-metrics-axe-11",
#      "connector_name" : "test-connector-name-metrics-axe-11",
#      "source_type" : "resource",
#      "event_type" : "resolve_cancel",
#      "component" : "test-component-metrics-axe-11",
#      "resource" : "test-resource-metrics-axe-11-1"
#    }
#    """
#    When I wait the end of event processing
#    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=created_alarms&parameters[]=active_alarms&parameters[]=instruction_alarms&parameters[]=ratio_instructions&parameters[]=correlation_alarms&parameters[]=ratio_correlation&parameters[]=non_displayed_alarms&parameters[]=ratio_non_displayed&parameters[]=ticket_active_alarms&parameters[]=ratio_tickets&parameters[]=ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
#    """json
#    {
#      "data": [
#        {
#          "title": "created_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "active_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "instruction_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "ratio_instructions",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 0
#            }
#          ]
#        },
#        {
#          "title": "correlation_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "ratio_correlation",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 100
#            }
#          ]
#        },
#        {
#          "title": "non_displayed_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "ratio_non_displayed",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 100
#            }
#          ]
#        },
#        {
#          "title": "ticket_active_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        },
#        {
#          "title": "ratio_tickets",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 100
#            }
#          ]
#        },
#        {
#          "title": "ack_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 2
#            }
#          ]
#        },
#        {
#          "title": "ack_active_alarms",
#          "data": [
#            {
#              "timestamp": {{ nowDate }},
#              "value": 1
#            }
#          ]
#        }
#      ]
#    }
#    """

  Scenario: given alarm with ticket should add ticket_active_alarms and without_ticket_active_alarms metrics for user
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-12-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-12-1"
        },
        {
          "name": "test-resource-metrics-axe-12-2"
        },
        {
          "name": "test-resource-metrics-axe-12-3"
        },
        {
          "name": "test-resource-metrics-axe-12-4"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-metrics-axe-12",
        "connector_name" : "test-connector-name-metrics-axe-12",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-12",
        "resource" : "test-resource-metrics-axe-12-1",
        "state" : 1
      },
      {
        "connector" : "test-connector-metrics-axe-12",
        "connector_name" : "test-connector-name-metrics-axe-12",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-12",
        "resource" : "test-resource-metrics-axe-12-2",
        "state" : 1
      },
      {
        "connector" : "test-connector-metrics-axe-12",
        "connector_name" : "test-connector-name-metrics-axe-12",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-12",
        "resource" : "test-resource-metrics-axe-12-3",
        "state" : 1
      },
      {
        "connector" : "test-connector-metrics-axe-12",
        "connector_name" : "test-connector-name-metrics-axe-12",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-metrics-axe-12",
        "resource" : "test-resource-metrics-axe-12-4",
        "state" : 1
      }
    ]
    """
    When I wait the end of 5 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-12",
      "connector_name" : "test-connector-name-metrics-axe-12",
      "source_type" : "resource",
      "event_type" : "assocticket",
      "component" : "test-component-metrics-axe-12",
      "resource" : "test-resource-metrics-axe-12-1",
      "ticket": "testticket",
      "user_id": "test-user-metrics-axe-12",
      "author": "test-user-metrics-axe-12",
      "initiator": "user"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ticket_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 4
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/rating?filter={{ .filterID }}&metric=ticket_active_alarms&criteria=3&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "label": "test-user-metrics-axe-12",
          "value": 1
        }
      ]
    }
    """

  Scenario: given double acked alarm should affect metrics only one time
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-axe-13-name",
      "entity_patterns": [
        {
          "name": "test-resource-metrics-axe-13"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-13",
      "connector_name" : "test-connector-name-metrics-axe-13",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-13",
      "resource" : "test-resource-metrics-axe-13",
      "state" : 1
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-13",
      "connector_name" : "test-connector-name-metrics-axe-13",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-13",
      "resource" : "test-resource-metrics-axe-13"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ack_alarms&parameters[]=average_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-metrics-axe-13",
      "connector_name" : "test-connector-name-metrics-axe-13",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-13",
      "resource" : "test-resource-metrics-axe-13"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ack_alarms&parameters[]=average_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
