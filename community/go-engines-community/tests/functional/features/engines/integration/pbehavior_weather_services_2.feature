Feature: get service weather
  I need to be able to get service weather

  @concurrent
  Scenario: given service without pbehavior should not get service by filter icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-1",
      "output_template": "Test service weather 13",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-1-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector":  "test-connector-pbehavior-weather-service-second-1",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-1",
      "resource": "test-resource-pbehavior-weather-service-second-1-1",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "component":  "test-component-pbehavior-weather-service-second-1",
        "resource": "test-resource-pbehavior-weather-service-second-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-1"
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "icon",
            "cond": {
              "type": "eq",
              "value": "build"
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
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given dependency with maintenance pbehavior and another dependency without pbehavior should get service by filter secondary_icon=maintenance
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-second-2",
        "connector_name": "test-connector-name-pbehavior-weather-service-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-2",
        "resource": "test-resource-pbehavior-weather-service-second-2-1",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-second-2",
        "connector_name": "test-connector-name-pbehavior-weather-service-second-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-2",
        "resource": "test-resource-pbehavior-weather-service-second-2-2",
        "state": 3,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-2-1",
      "output_template": "Test service weather 14-1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-weather-service-second-2-1",
                "test-resource-pbehavior-weather-service-second-2-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId1={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-2-2",
      "output_template": "Test service weather 14-2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-weather-service-second-2-1"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId1 }}"
      },
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId2 }}"
      }
    ]
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-second-2-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-2-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "component":  "test-component-pbehavior-weather-service-second-2",
        "resource": "test-resource-pbehavior-weather-service-second-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-2",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-second-2-1",
                "test-pbehavior-weather-service-second-2-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "secondary_icon",
            "cond": {
              "type": "eq",
              "value": "build"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}&sort_by=name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-2-1"
        },
        {
          "name": "test-pbehavior-weather-service-second-2-2"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-2",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-second-2-1",
                "test-pbehavior-weather-service-second-2-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "secondary_icon",
            "cond": {
              "type": "neq",
              "value": "build"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given dependency without pbehavior should not get service by filter secondary_icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-3",
      "output_template": "Test service weather 15",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-3-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector":  "test-connector-pbehavior-weather-service-second-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-3",
      "resource": "test-resource-pbehavior-weather-service-second-3-1",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "component":  "test-component-pbehavior-weather-service-second-3",
        "resource": "test-resource-pbehavior-weather-service-second-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-3",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-3"
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "secondary_icon",
            "cond": {
              "type": "eq",
              "value": "build"
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
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given service with maintenance pbehavior without alarm should get maintenance icon
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-4",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-4",
      "resource": "test-resource-pbehavior-weather-service-second-4",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-4",
      "output_template": "test-pbehavior-weather-service-second-4",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-second-4",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhenter",
      "component": "{{ .serviceId }}"
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-4",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-4"
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
          "name": "test-pbehavior-weather-service-second-4",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "build",
          "secondary_icon": "",
          "is_grey": true,
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
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-second-4",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-second-4",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "build",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
            }
          ]
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
  Scenario: given dependency with maintenance pbehavior without alarm should get maintenance secondary icon
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-5",
      "resource": "test-resource-pbehavior-weather-service-second-5",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-5",
      "output_template": "test-pbehavior-weather-service-second-5",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-second-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "component": "test-component-pbehavior-weather-service-second-5",
        "resource": "test-resource-pbehavior-weather-service-second-5",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-5",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-5"
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
          "name": "test-pbehavior-weather-service-second-5",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": true,
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
            "under_pbh": 1,
            "pbh_types": [
              {
                "count": 1,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "type": "maintenance"
                }
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

  @concurrent
  Scenario: given dependencies with maintenance pbehavior should keep maintenance secondary icon on alarm create and resolve
    Given I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "name": "test-resolve-rule-pbehavior-weather-service-second-6-name",
      "entity_pattern":[
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-6-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-6-2"
            }
          }
        ]
      ],
      "duration": {
        "value": 1,
        "unit": "s"
      },
      "priority": 1
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-second-6",
        "connector_name": "test-connector-name-pbehavior-weather-service-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-1",
        "state": 1,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-second-6",
        "connector_name": "test-connector-name-pbehavior-weather-service-second-6",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-2",
        "state": 0,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-6",
      "output_template": "test-pbehavior-weather-service-second-6",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-weather-service-second-6-1",
                "test-resource-pbehavior-weather-service-second-6-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-second-6",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-weather-service-second-6-1",
                "test-resource-pbehavior-weather-service-second-6-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "pbhenter",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-6",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterId={{ .lastResponse._id }}
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": true,
          "counters": {
            "depends": 2,
            "all": 1,
            "active": 0,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 0,
              "critical": 0
            },
            "acked": 0,
            "unacked": 0,
            "acked_under_pbh": 0,
            "under_pbh": 2,
            "pbh_types": [
              {
                "count": 2,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "type": "maintenance"
                }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-6",
      "resource": "test-resource-pbehavior-weather-service-second-6-2",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": true,
          "counters": {
            "depends": 2,
            "all": 2,
            "active": 0,
            "state": {
              "ok": 2,
              "minor": 0,
              "major": 0,
              "critical": 0
            },
            "acked": 0,
            "unacked": 0,
            "acked_under_pbh": 0,
            "under_pbh": 2,
            "pbh_types": [
              {
                "count": 2,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "type": "maintenance"
                }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-6",
      "resource": "test-resource-pbehavior-weather-service-second-6-1",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "resolve_close",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": true,
          "counters": {
            "depends": 2,
            "all": 1,
            "active": 0,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 0,
              "critical": 0
            },
            "acked": 0,
            "unacked": 0,
            "acked_under_pbh": 0,
            "under_pbh": 2,
            "pbh_types": [
              {
                "count": 2,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "type": "maintenance"
                }
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
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-6",
      "resource": "test-resource-pbehavior-weather-service-second-6-2",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "resolve_close",
        "component":  "test-component-pbehavior-weather-service-second-6",
        "resource": "test-resource-pbehavior-weather-service-second-6-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": true,
          "counters": {
            "depends": 2,
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
            "under_pbh": 2,
            "pbh_types": [
              {
                "count": 2,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "type": "maintenance"
                }
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

  @concurrent
  Scenario: given dependency with maintenance pbehavior should update service state correctly on alarm state change
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-7",
      "resource": "test-resource-pbehavior-weather-service-second-7",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-7",
      "output_template": "test-pbehavior-weather-service-second-7",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-7"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-second-7",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response pbhId={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "component":  "test-component-pbehavior-weather-service-second-7",
        "resource": "test-resource-pbehavior-weather-service-second-7",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-7",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterId={{ .lastResponse._id }}
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-7",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": true,
          "counters": {
            "depends": 1,
            "all": 1,
            "active": 0,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 0,
              "critical": 0
            },
            "acked": 0,
            "unacked": 0,
            "acked_under_pbh": 0,
            "under_pbh": 1,
            "pbh_types": [
              {
                "count": 1,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "type": "maintenance"
                }
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-7",
      "resource": "test-resource-pbehavior-weather-service-second-7",
      "state": 0,
      "output": "noveo alarm"
    }
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-7",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": true,
          "counters": {
            "depends": 1,
            "all": 1,
            "active": 0,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 0,
              "critical": 0
            },
            "acked": 0,
            "unacked": 0,
            "acked_under_pbh": 0,
            "under_pbh": 1,
            "pbh_types": [
              {
                "count": 1,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "type": "maintenance"
                }
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
    When I do DELETE /api/v4/pbehaviors/{{ .pbhId }}
    Then the response code should be 204
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhleave",
        "component":  "test-component-pbehavior-weather-service-second-7",
        "resource": "test-resource-pbehavior-weather-service-second-7",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-7",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "",
          "is_grey": false,
          "counters": {
            "depends": 1,
            "all": 1,
            "active": 1,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 0,
              "critical": 0
            },
            "acked": 0,
            "unacked": 1,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
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
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-second-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-second-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-second-7",
      "resource": "test-resource-pbehavior-weather-service-second-7",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-7",
        "resource": "test-resource-pbehavior-weather-service-second-7",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-second-7",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "",
          "is_grey": false,
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
  Scenario: given service with pbehavior should get service by filter is_grey=true and should not get service by filter is_grey=false
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-8-1",
      "output_template": "Test service weather 20-1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-8"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId1={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-8-2",
      "output_template": "Test service weather 20-2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-8"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-second-8",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-second-8-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhenter",
      "component": "{{ .serviceId1 }}"
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-8",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-second-8-1",
                "test-pbehavior-weather-service-second-8-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "is_grey",
            "cond": {
              "type": "eq",
              "value": true
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
          "name": "test-pbehavior-weather-service-second-8-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-second-8",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-second-8-1",
                "test-pbehavior-weather-service-second-8-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "is_grey",
            "cond": {
              "type": "eq",
              "value": false
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
          "name": "test-pbehavior-weather-service-second-8-2"
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
  Scenario: given dependency with pbehavior should get service by filter is_grey=true and should not get service by filter is_grey=false
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-second-9",
        "connector_name": "test-connector-name-pbehavior-weather-service-second-9",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-9",
        "resource": "test-resource-pbehavior-weather-service-second-9-1",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-second-9",
        "connector_name": "test-connector-name-pbehavior-weather-service-second-9",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-second-9",
        "resource": "test-resource-pbehavior-weather-service-second-9-2",
        "state": 2,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-9-1",
      "output_template": "Test service weather 21-1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-9-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId1={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-second-9-2",
      "output_template": "Test service weather 21-2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-weather-service-second-9-1",
                "test-resource-pbehavior-weather-service-second-9-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    Then I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId1 }}"
      },
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      },
      {
        "event_type": "activate",
        "component": "{{ .serviceId2 }}"
      }
    ]
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-second-9",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-second-9-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "component":  "test-component-pbehavior-weather-service-second-9",
        "resource": "test-resource-pbehavior-weather-service-second-9-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-second-9",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-second-9-1",
                "test-pbehavior-weather-service-second-9-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "is_grey",
            "cond": {
              "type": "eq",
              "value": true
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
          "name": "test-pbehavior-weather-service-second-9-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-second-9",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-second-9-1",
                "test-pbehavior-weather-service-second-9-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "is_grey",
            "cond": {
              "type": "eq",
              "value": false
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
          "name": "test-pbehavior-weather-service-second-9-2"
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
