Feature: get service entities
  I need to be able to get service entities

  @concurrent
  Scenario: given service for one entity should get one entity
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
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
    When I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "api",
        "connector_name": "api",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-service-weather-entity-1",
      "connector": "test-connector-service-weather-entity-1",
      "connector_name": "test-connector_name-service-weather-entity-1",
      "component": "test-component-service-weather-entity-1",
      "resource": "test-resource-service-weather-entity-1",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-weather-entity-1",
        "connector_name": "test-connector_name-service-weather-entity-1",
        "component": "test-component-service-weather-entity-1",
        "resource": "test-resource-service-weather-entity-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/weather-services/{{ .serviceId }}
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

  @concurrent
  Scenario: given service for multiple entities should get multiple entities
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
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
    When I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "api",
        "connector_name": "api",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-service-weather-entity-2",
      "connector": "test-connector-service-weather-entity-2",
      "connector_name": "test-connector_name-service-weather-entity-2",
      "component": "test-component-service-weather-entity-2",
      "resource": "test-resource-service-weather-entity-2-1",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-weather-entity-2",
        "connector_name": "test-connector_name-service-weather-entity-2",
        "component": "test-component-service-weather-entity-2",
        "resource": "test-resource-service-weather-entity-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-service-weather-entity-2",
      "connector": "test-connector-service-weather-entity-2",
      "connector_name": "test-connector_name-service-weather-entity-2",
      "component": "test-component-service-weather-entity-2",
      "resource": "test-resource-service-weather-entity-2-2",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-weather-entity-2",
        "connector_name": "test-connector_name-service-weather-entity-2",
        "component": "test-component-service-weather-entity-2",
        "resource": "test-resource-service-weather-entity-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 3,
      "output": "test-output-service-weather-entity-2",
      "connector": "test-connector-service-weather-entity-2",
      "connector_name": "test-connector_name-service-weather-entity-2",
      "component": "test-component-service-weather-entity-2",
      "resource": "test-resource-service-weather-entity-2-3",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-weather-entity-2",
        "connector_name": "test-connector_name-service-weather-entity-2",
        "component": "test-component-service-weather-entity-2",
        "resource": "test-resource-service-weather-entity-2-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/weather-services/{{ .serviceId }}?sort_by=state&sort=desc
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
    When I do GET /api/v4/weather-services/{{ .serviceId }}?sort_by=impact_state&sort=desc
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

  @concurrent
  Scenario: given service for one entity and no open alarms should get one entity with ok state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
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
    When I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "api",
        "connector_name": "api",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "test-output-service-weather-entity-3",
      "connector": "test-connector-service-weather-entity-3",
      "connector_name": "test-connector_name-service-weather-entity-3",
      "component": "test-component-service-weather-entity-3",
      "resource": "test-resource-service-weather-entity-3",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-service-weather-entity-3",
        "connector_name": "test-connector_name-service-weather-entity-3",
        "component": "test-component-service-weather-entity-3",
        "resource": "test-resource-service-weather-entity-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/weather-services/{{ .serviceId }}
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

  @concurrent
  Scenario: given service for one entity, weather service entities should return event stats
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
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
    When I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "api",
        "connector_name": "api",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "test-output-service-weather-entity-4",
      "connector": "test-connector-service-weather-entity-4",
      "connector_name": "test-connector_name-service-weather-entity-4",
      "component": "test-component-service-weather-entity-4",
      "resource": "test-resource-service-weather-entity-4",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-service-weather-entity-4",
        "connector_name": "test-connector_name-service-weather-entity-4",
        "component": "test-component-service-weather-entity-4",
        "resource": "test-resource-service-weather-entity-4",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-service-weather-entity-4",
      "connector": "test-connector-service-weather-entity-4",
      "connector_name": "test-connector_name-service-weather-entity-4",
      "component": "test-component-service-weather-entity-4",
      "resource": "test-resource-service-weather-entity-4",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-service-weather-entity-4",
        "connector_name": "test-connector_name-service-weather-entity-4",
        "component": "test-component-service-weather-entity-4",
        "resource": "test-resource-service-weather-entity-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I wait 3s
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "test-output-service-weather-entity-4",
      "connector": "test-connector-service-weather-entity-4",
      "connector_name": "test-connector_name-service-weather-entity-4",
      "component": "test-component-service-weather-entity-4",
      "resource": "test-resource-service-weather-entity-4",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-service-weather-entity-4",
        "connector_name": "test-connector_name-service-weather-entity-4",
        "component": "test-component-service-weather-entity-4",
        "resource": "test-resource-service-weather-entity-4",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/weather-services/{{ .serviceId }}
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

  @concurrent
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

  @concurrent
  Scenario: given service for one entity, old last_ko shouldn't be updated  
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "test-output-service-weather-entity-6",
      "connector": "test-connector-service-weather-entity-6",
      "connector_name": "test-connector_name-service-weather-entity-6",
      "component": "test-component-service-weather-entity-6",
      "resource": "test-resource-service-weather-entity-6",
      "source_type": "resource"
    }
    """
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

  @concurrent
  Scenario: given service and events, should return assigned declare ticket rules in get weather service entities request
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-service-weather-entity-7",
        "connector": "test-connector-service-weather-entity-7",
        "connector_name": "test-connector_name-service-weather-entity-7",
        "component": "test-component-service-weather-entity-7",
        "resource": "test-resource-service-weather-entity-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-service-weather-entity-7",
        "connector": "test-connector-service-weather-entity-7",
        "connector_name": "test-connector_name-service-weather-entity-7",
        "component": "test-component-service-weather-entity-7",
        "resource": "test-resource-service-weather-entity-7-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-service-weather-entity-7",
        "connector": "test-connector-service-weather-entity-7",
        "connector_name": "test-connector_name-service-weather-entity-7",
        "component": "test-component-service-weather-entity-7",
        "resource": "test-resource-service-weather-entity-7-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-service-weather-entity-7",
        "connector": "test-connector-service-weather-entity-7",
        "connector_name": "test-connector_name-service-weather-entity-7",
        "component": "test-component-service-weather-entity-7",
        "resource": "test-resource-service-weather-entity-7-4",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-service-weather-entity-7",
        "connector": "test-connector-service-weather-entity-7",
        "connector_name": "test-connector_name-service-weather-entity-7",
        "component": "test-component-service-weather-entity-7",
        "resource": "test-resource-service-weather-entity-7-5",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-weather-service-rule-7-1",
      "system_name": "test-alarm-service-rule-7-1-name",
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
              "value": "test-resource-service-weather-entity-7-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-7-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-7-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-weather-service-rule-7-2",
      "system_name": "test-alarm-service-rule-7-2-name",
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
              "value": "test-resource-service-weather-entity-7-3"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-entity-7-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleId2={{ .lastResponse._id }}
    When I do GET /api/v4/weather-services/test-service-weather-entity-7?with_declare_tickets=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-service-weather-entity-7-1",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleId1 }}",
              "name": "test-weather-service-rule-7-1"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-7-2",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleId1 }}",
              "name": "test-weather-service-rule-7-1"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-7-3",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleId1 }}",
              "name": "test-weather-service-rule-7-1"
            },
            {
              "_id": "{{ .ruleId2 }}",
              "name": "test-weather-service-rule-7-2"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-7-4",
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleId2 }}",
              "name": "test-weather-service-rule-7-2"
            }
          ]
        },
        {
          "name": "test-resource-service-weather-entity-7-5"
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
