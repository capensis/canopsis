Feature: update service weather on event
  I need to be able to see new service weather state on event

  Scenario: given service for entity should update service weather on entity event
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-1",
      "name": "test-service-weather-1",
      "output_template": "test-service-weather-1-output",
      "category": "test-category-service-weather",
      "enabled": true,
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-1"
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
      "connector": "test-connector-service-weather-1",
      "connector_name": "test-connector-name-service-weather-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-1",
      "resource": "test-resource-service-weather-1",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-service-weather-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-weather-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-service-weather-1",
          "name": "test-service-weather-1",
          "connector": "service",
          "connector_name": "service",
          "component": "test-service-weather-1",
          "resource": "",
          "state": {"val": 2},
          "status": {"val": 1},
          "ack": null,
          "snooze": null,
          "output": "test-service-weather-1-output",
          "infos": {},
          "icon": "person",
          "secondary_icon": "",
          "category": {
            "_id": "test-category-service-weather"
          },
          "is_grey": false,
          "is_action_required": true,
          "counters": {
            "depends": 1,
            "all": 1,
            "active": 1,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 1,
              "critical": 0
            },
            "acked": 0,
            "unacked": 1,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
          },
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

  Scenario: given service for multiple entities should set service weather state as worst entity state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-service-weather-2",
      "enabled": true,
      "output_template": "Test-service-weather-2",
      "category": "test-category-service-weather",
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-weather-2-1",
                "test-resource-service-weather-2-2",
                "test-resource-service-weather-2-3"
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
      "connector": "test-connector-service-weather-2",
      "connector_name": "test-connector-name-service-weather-2",
      "source_type": "resource",
      "event_type": "check",
      "category": "test-category-service-weather",
      "component":  "test-component-service-weather-2",
      "resource": "test-resource-service-weather-2-1",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-weather-2",
      "connector_name": "test-connector-name-service-weather-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-2",
      "resource": "test-resource-service-weather-2-2",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-weather-2",
      "connector_name": "test-connector-name-service-weather-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-2",
      "resource": "test-resource-service-weather-2-3",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-service-weather-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-weather-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-service-weather-2",
          "connector": "service",
          "connector_name": "service",
          "resource": "",
          "state": {"val": 3},
          "status": {"val": 1},
          "ack": null,
          "snooze": null,
          "infos": {},
          "icon": "wb_cloudy",
          "category": {
            "_id": "test-category-service-weather"
          },
          "secondary_icon": "",
          "is_grey": false,
          "is_action_required": true,
          "counters": {
            "depends": 3,
            "all": 3,
            "active": 3,
            "state": {
              "ok": 0,
              "minor": 1,
              "major": 1,
              "critical": 1
            },
            "acked": 0,
            "unacked": 3,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
          },
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

  Scenario: given service for entity and no open alarm should get ok service weather state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-3",
      "name": "test-service-weather-3",
      "enabled": true,
      "output_template": "Test-service-weather-3",
      "category": "test-category-service-weather",
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-3"
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
      "connector": "test-connector-service-weather-3",
      "connector_name": "test-connector-name-service-weather-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-3",
      "resource": "test-resource-service-weather-3",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-service-weather-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-weather-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-service-weather-3",
          "connector": "",
          "connector_name": "",
          "component": "",
          "resource": "",
          "state": {"val": 0, "t": null},
          "status": {"val": 0, "t": null},
          "ack": null,
          "snooze": null,
          "last_update_date": null,
          "infos": {},
          "category": {
            "_id": "test-category-service-weather"
          },
          "icon": "wb_sunny",
          "secondary_icon": "",
          "is_grey": false,
          "is_action_required": false,
          "counters": {
            "depends": 1,
            "all": 0,
            "active": 0,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 0,
              "critical": 0
            },
            "acked": 0,
            "unacked": 0,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
          },
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

  Scenario: given service for multiple entities and acked alarms
    should return false in is_action_required
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-4",
      "name": "test-service-weather-4",
      "enabled": true,
      "output_template": "Test-service-weather-4",
      "category": "test-category-service-weather",
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-weather-4-1",
                "test-resource-service-weather-4-2"
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
      "connector": "test-connector-service-weather-4",
      "connector_name": "test-connector-name-service-weather-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-4",
      "resource": "test-resource-service-weather-4-1",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-weather-4",
      "connector_name": "test-connector-name-service-weather-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-4",
      "resource": "test-resource-service-weather-4-2",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-weather-4",
      "connector_name": "test-connector-name-service-weather-4",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-service-weather-4",
      "resource": "test-resource-service-weather-4-1",
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-weather-4",
      "connector_name": "test-connector-name-service-weather-4",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-service-weather-4",
      "resource": "test-resource-service-weather-4-2",
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-service-weather-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-weather-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-service-weather-4",
          "connector": "service",
          "connector_name": "service",
          "resource": "",
          "state": {"val": 3},
          "status": {"val": 1},
          "ack": null,
          "snooze": null,
          "infos": {},
          "category": {
            "_id": "test-category-service-weather"
          },
          "icon": "wb_cloudy",
          "secondary_icon": "",
          "is_grey": false,
          "is_action_required": false,
          "counters": {
            "depends": 2,
            "all": 2,
            "active": 2,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 1,
              "critical": 1
            },
            "acked": 2,
            "unacked": 0,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
          },
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

  Scenario: given service for entity should calculate impact_state regarding of impact_level
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-service-weather-5",
      "name": "test-service-weather-5",
      "enabled": true,
      "output_template": "Test-service-weather-5",
      "category": "test-category-service-weather",
      "impact_level": 5,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-5"
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
      "connector": "test-connector-service-weather-5",
      "connector_name": "test-connector-name-service-weather-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-5",
      "resource": "test-resource-service-weather-5",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-service-weather-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-weather-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-service-weather-5",
          "connector": "service",
          "connector_name": "service",
          "resource": "",
          "state": {"val": 2},
          "status": {"val": 1},
          "ack": null,
          "snooze": null,
          "infos": {},
          "icon": "person",
          "impact_level": 5,
          "impact_state": 10,
          "secondary_icon": "",
          "is_grey": false,
          "category": {
            "_id": "test-category-service-weather"
          },
          "is_action_required": true,
          "counters": {
            "depends": 1,
            "all": 1,
            "active": 1,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 1,
              "critical": 0
            },
            "acked": 0,
            "unacked": 1,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
          },
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

  Scenario: given service should get service by filter state=critical
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-service-weather-6-1",
      "enabled": true,
      "output_template": "Test-service-weather-6-1",
      "category": "test-category-service-weather",
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-6-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-service-weather-6-2",
      "enabled": true,
      "output_template": "Test-service-weather-6-2",
      "category": "test-category-service-weather",
      "impact_level": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-weather-6-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-weather-6",
      "connector_name": "test-connector-name-service-weather-6",
      "source_type": "resource",
      "event_type": "check",
      "category": "test-category-service-weather",
      "component":  "test-component-service-weather-6",
      "resource": "test-resource-service-weather-6-1",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-service-weather-6",
      "connector_name": "test-connector-name-service-weather-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-service-weather-6",
      "resource": "test-resource-service-weather-6-2",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-service-weather-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-service-weather-6-1",
                "test-service-weather-6-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "state.val",
            "cond": {
              "type": "eq",
              "value": 1
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-service-weather-6-1"
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-service-weather-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-service-weather-6-1",
                "test-service-weather-6-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "state.val",
            "cond": {
              "type": "neq",
              "value": 1
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-service-weather-6-2"
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
