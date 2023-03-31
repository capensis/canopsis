Feature: get service entities
  I need to be able to get service entities

  Scenario: given service for one entity should get one entity
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-entity-1",
      "name": "test-service-weather-entity-1",
      "output_template": "Test-service-weather-entity-1",
      "category": "test-category-service-weather-entities",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-1",
      "connector_name" : "test-connector_name-service-weather-entity-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-1",
      "resource" : "test-resource-service-weather-entity-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-1/test-component-service-weather-entity-1",
          "name": "test-resource-service-weather-entity-1",
          "connector": "test-connector-service-weather-entity-1",
          "connector_name": "test-connector_name-service-weather-entity-1",
          "component": "test-component-service-weather-entity-1",
          "resource": "test-resource-service-weather-entity-1",
          "source_type": "resource",
          "state": {"val": 2},
          "status": {"val": 1},
          "impact_state": 2,
          "impact_level": 1,
          "category": null,
          "is_grey": false,
          "icon": "person",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehavior_info": null,
          "pbehaviors": [],
          "depends_count": 0
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

  Scenario: given service for multiple entities should get multiple entities
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-entity-2",
      "name": "test-service-weather-entity-2",
      "output_template": "Test-service-weather-entity-2",
      "category": "test-category-service-weather-entities",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-weather-entity-2-1",
                "test-resource-service-weather-entity-2-2",
                "test-resource-service-weather-entity-2-3"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-2",
      "connector_name" : "test-connector_name-service-weather-entity-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-2",
      "resource" : "test-resource-service-weather-entity-2-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-2",
      "connector_name" : "test-connector_name-service-weather-entity-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-2",
      "resource" : "test-resource-service-weather-entity-2-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-2",
      "connector_name" : "test-connector_name-service-weather-entity-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-2",
      "resource" : "test-resource-service-weather-entity-2-3",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-2?sort_by=state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-2-3/test-component-service-weather-entity-2",
          "name": "test-resource-service-weather-entity-2-3",
          "connector": "test-connector-service-weather-entity-2",
          "connector_name": "test-connector_name-service-weather-entity-2",
          "component": "test-component-service-weather-entity-2",
          "resource": "test-resource-service-weather-entity-2-3",
          "state": {"val": 3},
          "status": {"val": 1},
          "impact_state": 3,
          "impact_level": 1,
          "is_grey": false,
          "icon": "wb_cloudy",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehavior_info": null,
          "pbehaviors": [],
          "depends_count": 0
        },
        {
          "_id": "test-resource-service-weather-entity-2-2/test-component-service-weather-entity-2",
          "name": "test-resource-service-weather-entity-2-2",
          "connector": "test-connector-service-weather-entity-2",
          "connector_name": "test-connector_name-service-weather-entity-2",
          "component": "test-component-service-weather-entity-2",
          "resource": "test-resource-service-weather-entity-2-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "impact_state": 2,
          "impact_level": 1,
          "is_grey": false,
          "icon": "person",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehavior_info": null,
          "pbehaviors": [],
          "depends_count": 0
        },
        {
          "_id": "test-resource-service-weather-entity-2-1/test-component-service-weather-entity-2",
          "name": "test-resource-service-weather-entity-2-1",
          "connector": "test-connector-service-weather-entity-2",
          "connector_name": "test-connector_name-service-weather-entity-2",
          "component": "test-component-service-weather-entity-2",
          "resource": "test-resource-service-weather-entity-2-1",
          "state": {"val": 1},
          "status": {"val": 1},
          "impact_state": 1,
          "impact_level": 1,
          "is_grey": false,
          "icon": "person",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehavior_info": null,
          "pbehaviors": [],
          "depends_count": 0
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/weather-services/test-service-weather-entity-2?sort_by=impact_state&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-2-3/test-component-service-weather-entity-2",
          "impact_state": 3
        },
        {
          "_id": "test-resource-service-weather-entity-2-2/test-component-service-weather-entity-2",
          "impact_state": 2
        },
        {
          "_id": "test-resource-service-weather-entity-2-1/test-component-service-weather-entity-2",
          "impact_state": 1
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given service for one entity and no open alarms should get one entity with ok state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-entity-3",
      "name": "test-service-weather-entity-3",
      "output_template": "Test-service-weather-entity-3",
      "category": "test-category-service-weather-entities",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-3"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-3",
      "connector_name" : "test-connector_name-service-weather-entity-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-3",
      "resource" : "test-resource-service-weather-entity-3",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-3/test-component-service-weather-entity-3",
          "name": "test-resource-service-weather-entity-3",
          "connector": "",
          "connector_name": "",
          "component": "",
          "resource": "",
          "state": {"val": 0, "t": null},
          "status": {"val": 0, "t": null},
          "last_update_date": null,
          "alarm_creation_date": null,
          "impact_state": 0,
          "impact_level": 1,
          "is_grey": false,
          "icon": "wb_sunny",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehavior_info": null,
          "pbehaviors": []
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

  Scenario: given service for one entity, weather service entities should return event stats
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-entity-4",
      "name": "test-service-weather-entity-4",
      "output_template": "Test-service-weather-entity-4",
      "category": "test-category-service-weather-entities",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-4",
      "connector_name" : "test-connector_name-service-weather-entity-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-4",
      "resource" : "test-resource-service-weather-entity-4",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-4",
      "connector_name" : "test-connector_name-service-weather-entity-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-4",
      "resource" : "test-resource-service-weather-entity-4",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I wait 3s
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-4",
      "connector_name" : "test-connector_name-service-weather-entity-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-4",
      "resource" : "test-resource-service-weather-entity-4",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-4
    Then the response code should be 200
    Then the response key "data.0.stats.last_event" should not be "null"
    Then the response key "data.0.stats.last_ko" should not be "null"
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-4/test-component-service-weather-entity-4",
          "name": "test-resource-service-weather-entity-4",
          "stats": {
            "ok": 2,
            "ko": 1
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

  Scenario: given service for one entity, weather service entities shouldn't return old event stats
    Given I am admin
    When I do GET /api/v4/weather-services/test-service-weather-entity-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-5/test-component-service-weather-entity-5",
          "name": "test-resource-service-weather-entity-5",
          "stats": {
            "ok": 0,
            "ko": 0
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

  Scenario: given service for one entity, old last_ko shouldn't be updated  
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-service-weather-entity-6",
      "connector_name" : "test-connector_name-service-weather-entity-6",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-6",
      "resource" : "test-resource-service-weather-entity-6",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-6
    Then the response code should be 200    
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-6/test-component-service-weather-entity-6",
          "name": "test-resource-service-weather-entity-6",
          "stats": {
            "ok": 1,
            "ko": 0,
            "last_ko": 1000000000
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

  Scenario: given service and events, should return instruction statuses and icons in get weather service entities request
    Given I am admin
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-service-weather-entity-7",
        "connector_name" : "test-connector_name-service-weather-entity-7",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-7",
        "resource" : "test-resource-service-weather-entity-7-1",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-7",
        "connector_name" : "test-connector_name-service-weather-entity-7",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-7",
        "resource" : "test-resource-service-weather-entity-7-2",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-7",
        "connector_name" : "test-connector_name-service-weather-entity-7",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-7",
        "resource" : "test-resource-service-weather-entity-7-3",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-7",
        "connector_name" : "test-connector_name-service-weather-entity-7",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-7",
        "resource" : "test-resource-service-weather-entity-7-4",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-7",
        "connector_name" : "test-connector_name-service-weather-entity-7",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-7",
        "resource" : "test-resource-service-weather-entity-7-5",
        "state" : 1,
        "output" : "noveo alarm"
      }
    ]
    """
    When I wait the end of 9 events processing
    When I do GET /api/v4/alarms?search=test-resource-service-weather-entity-7&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-service-weather-entity-7-1"
          }
        },
        {
          "v": {
            "resource": "test-resource-service-weather-entity-7-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-service-weather-entity-7-3"
          }
        },
        {
          "v": {
            "resource": "test-resource-service-weather-entity-7-4"
          }
        },
        {
          "v": {
            "resource": "test-resource-service-weather-entity-7-5"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmID2={{ (index .lastResponse.data 1)._id }}
    When I save response alarmID3={{ (index .lastResponse.data 2)._id }}
    When I save response alarmID4={{ (index .lastResponse.data 3)._id }}
    When I save response alarmID5={{ (index .lastResponse.data 4)._id }}
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID4 }}",
      "instruction": "test-instruction-service-weather-entity-7-4"
    }
    """
    Then the response code should be 200
    When I do POST /api/v4/cat/executions:
    """json
    {
      "alarm": "{{ .alarmID5 }}",
      "instruction": "test-instruction-service-weather-entity-7-5"
    }
    """
    Then the response code should be 200
    When I wait the end of 4 events processing
    When I wait 3s
    When I do GET /api/v4/weather-services/test-service-weather-entity-7?with_instructions=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-service-weather-entity-7-1",
          "instruction_execution_icon": 9
        },
        {
          "name": "test-resource-service-weather-entity-7-2",
          "instruction_execution_icon": 10,
          "successful_auto_instructions": [
            "test-instruction-service-weather-entity-7-2-name"
          ]
        },
        {
          "name": "test-resource-service-weather-entity-7-3",
          "instruction_execution_icon": 3,
          "failed_auto_instructions": [
            "test-instruction-service-weather-entity-7-3-name"
          ]
        },
        {
          "name": "test-resource-service-weather-entity-7-4",
          "instruction_execution_icon": 11,
          "successful_manual_instructions": [
            "test-instruction-service-weather-entity-7-4-name"
          ]
        },
        {
          "name": "test-resource-service-weather-entity-7-5",
          "instruction_execution_icon": 4,
          "failed_manual_instructions": [
            "test-instruction-service-weather-entity-7-5-name"
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """

  Scenario: given service and events, should return assigned declare ticket rules in get weather service entities request
    Given I am admin
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-service-weather-entity-8",
        "connector_name" : "test-connector_name-service-weather-entity-8",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-8",
        "resource" : "test-resource-service-weather-entity-8-1",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-8",
        "connector_name" : "test-connector_name-service-weather-entity-8",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-8",
        "resource" : "test-resource-service-weather-entity-8-2",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-8",
        "connector_name" : "test-connector_name-service-weather-entity-8",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-8",
        "resource" : "test-resource-service-weather-entity-8-3",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-8",
        "connector_name" : "test-connector_name-service-weather-entity-8",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-8",
        "resource" : "test-resource-service-weather-entity-8-4",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-service-weather-entity-8",
        "connector_name" : "test-connector_name-service-weather-entity-8",
        "source_type" : "resource",
        "event_type" : "check",
        "component" :  "test-component-service-weather-entity-8",
        "resource" : "test-resource-service-weather-entity-8-5",
        "state" : 1,
        "output" : "noveo alarm"
      }
    ]
    """
    When I wait the end of 5 events processing
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-weather-service-rule-8-1",
      "system_name": "test-alarm-service-rule-8-1-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-8-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-8-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-8-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-weather-service-rule-8-2",
      "system_name": "test-alarm-service-rule-8-2-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-8-3"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-8-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID2={{ .lastResponse._id }}
    When I do GET /api/v4/weather-services/test-service-weather-entity-8?with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-service-weather-entity-8-1",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID1 }}",
              "name": "test-weather-service-rule-8-1"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-8-2",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID1 }}",
              "name": "test-weather-service-rule-8-1"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-8-3",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID1 }}",
              "name": "test-weather-service-rule-8-1"
            },
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-weather-service-rule-8-2"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-8-4",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleID2 }}",
              "name": "test-weather-service-rule-8-2"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-8-5"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """
    Then the response key "data.4.assigned_declare_ticket_rules" should not exist
