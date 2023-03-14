Feature: get service weather
  I need to be able to get service weather

  @concurrent
  Scenario: given one dependency with maintenance pbehavior should get maintenance secondary icon and grey flag
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-1",
      "connector_name": "test-connector-name-pbehavior-weather-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-1",
      "resource": "test-resource-pbehavior-weather-service-1",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-1",
      "output_template": "Test service weather 1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-1"
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
      "name": "test-pbehavior-weather-service-1",
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
              "value": "test-resource-pbehavior-weather-service-1"
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
        "component": "test-component-pbehavior-weather-service-1",
        "resource": "test-resource-pbehavior-weather-service-1",
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
      "title": "test-widgetfilter-pbehavior-weather-service-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-1"
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
          "name": "test-pbehavior-weather-service-1",
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

  @concurrent
  Scenario: given one dependency with active pbehavior should get state icon
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-2",
      "connector_name": "test-connector-name-pbehavior-weather-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-2",
      "resource": "test-resource-pbehavior-weather-service-2",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-2",
      "output_template": "Test service weather 2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-2"
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
      "name": "test-pbehavior-weather-service-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-2"
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
      "component":  "test-component-pbehavior-weather-service-2",
      "resource": "test-resource-pbehavior-weather-service-2",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-2",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-2"
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
          "name": "test-pbehavior-weather-service-2",
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
  Scenario: given dependency with maintenance pbehavior and dependency without pbehavior should get secondary maintenance icon
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector":  "test-connector-pbehavior-weather-service-3",
        "connector_name": "test-connector-name-pbehavior-weather-service-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-3",
        "resource": "test-resource-pbehavior-weather-service-3-1",
        "state": 3,
        "output": "noveo alarm"
      },
      {
        "connector":  "test-connector-pbehavior-weather-service-3",
        "connector_name": "test-connector-name-pbehavior-weather-service-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-3",
        "resource": "test-resource-pbehavior-weather-service-3-2",
        "state": 2,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-3",
      "output_template": "Test service weather 3",
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
                "test-resource-pbehavior-weather-service-3-1",
                "test-resource-pbehavior-weather-service-3-2"
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
      "name": "test-pbehavior-weather-service-3-1",
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
              "value": "test-resource-pbehavior-weather-service-3-1"
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
        "component":  "test-component-pbehavior-weather-service-3",
        "resource": "test-resource-pbehavior-weather-service-3-1",
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
      "title": "test-widgetfilter-pbehavior-weather-service-3",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-3"
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
          "name": "test-pbehavior-weather-service-3",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 2,
            "all": 2,
            "active": 1,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 1,
              "critical": 0
            },
            "acked": 0,
            "unacked": 1,
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
  Scenario: given service with maintenance pbehavior should get maintenance icon and grey flag
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-4",
      "output_template": "Test service weather 4",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-4"
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
      "name": "test-pbehavior-weather-service-4",
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
              "value": "test-pbehavior-weather-service-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response pbehaviorId={{ .lastResponse._id }}
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "{{ .pbehaviorId }}",
      "message": "First comment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "{{ .pbehaviorId }}",
      "message": "Second comment"
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
    When I send an event:
    """json
    {
      "connector":  "test-connector-pbehavior-weather-service-4",
      "connector_name": "test-connector-name-pbehavior-weather-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-4",
      "resource": "test-resource-pbehavior-weather-service-4",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "component":  "test-component-pbehavior-weather-service-4",
        "resource": "test-resource-pbehavior-weather-service-4",
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
      "title": "test-widgetfilter-pbehavior-weather-service-4",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-4"
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
          "name": "test-pbehavior-weather-service-4",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "build",
          "secondary_icon": "",
          "is_grey": true,
          "counters": {
            "depends": 1,
            "all": 1,
            "active": 1,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 0,
              "critical": 1
            },
            "acked": 0,
            "unacked": 1,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
          },
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-4",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-4",
              "last_comment": {
                "message": "Second comment"
              },
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
  Scenario: given service with active pbehavior should get state icon
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-5",
      "output_template": "Test service weather 5",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-5"
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
        "component":  "{{ .serviceId }}",
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
      "name": "test-pbehavior-weather-service-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-5"
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
    When I send an event:
    """json
    {
      "connector":  "test-connector-pbehavior-weather-service-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-5",
      "resource": "test-resource-pbehavior-weather-service-5",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "component":  "test-component-pbehavior-weather-service-5",
        "resource": "test-resource-pbehavior-weather-service-5",
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
      "title": "test-widgetfilter-pbehavior-weather-service-5",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-5"
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
          "name": "test-pbehavior-weather-service-5",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "wb_cloudy",
          "secondary_icon": "",
          "is_grey": false,
          "counters": {
            "depends": 1,
            "all": 1,
            "active": 1,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 0,
              "critical": 1
            },
            "acked": 0,
            "unacked": 1,
            "acked_under_pbh": 0,
            "under_pbh": 0,
            "pbh_types": []
          },
          "pbehavior_info": {
            "canonical_type": "active",
            "name": "test-pbehavior-weather-service-5",
            "icon_name": "brightness_3"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-5",
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
                "_id": "test-active-type-to-engine",
                "icon_name": "brightness_3",
                "name": "Engine active",
                "type": "active"
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
  Scenario: given service with maintenance pbehavior and one dependency with maintenance pbehavior should get maintenance icon and secondary icon
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector":  "test-connector-pbehavior-weather-service-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-6",
      "resource": "test-resource-pbehavior-weather-service-6",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-6",
      "output_template": "Test service weather 6",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-6"
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
      "name": "test-pbehavior-weather-service-6-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-6"
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-6-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-6"
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
        "component":  "test-component-pbehavior-weather-service-6",
        "resource": "test-resource-pbehavior-weather-service-6",
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
      "title": "test-widgetfilter-pbehavior-weather-service-6",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-6"
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
          "name": "test-pbehavior-weather-service-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "build",
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
          },
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-6-1",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-6-1",
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
  Scenario: given service with maintenance pbehavior and dependency with maintenance pbehavior and another dependency without pbehavior should get maintenance icon and maintenance secondary icon
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector":  "test-connector-pbehavior-weather-service-7",
        "connector_name": "test-connector-name-pbehavior-weather-service-7",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-7",
        "resource": "test-resource-pbehavior-weather-service-7-1",
        "state": 3,
        "output": "noveo alarm"
      },
      {
        "connector":  "test-connector-pbehavior-weather-service-7",
        "connector_name": "test-connector-name-pbehavior-weather-service-7",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-7",
        "resource": "test-resource-pbehavior-weather-service-7-2",
        "state": 2,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-7",
      "output_template": "Test service weather 7",
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
                "test-resource-pbehavior-weather-service-7-1",
                "test-resource-pbehavior-weather-service-7-2"
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
      "name": "test-pbehavior-weather-service-7-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-7"
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-7-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-7-1"
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
        "component":  "test-component-pbehavior-weather-service-7",
        "resource": "test-resource-pbehavior-weather-service-7-1",
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
      "title": "test-widgetfilter-pbehavior-weather-service-7",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-7"
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
          "name": "test-pbehavior-weather-service-7",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "build",
          "secondary_icon": "build",
          "is_grey": true,
          "counters": {
            "depends": 2,
            "all": 2,
            "active": 1,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 1,
              "critical": 0
            },
            "acked": 0,
            "unacked": 1,
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
          },
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-7-1",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-7-1",
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
  Scenario: given service with maintenance pbehavior should get service by filter icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-8-1",
      "output_template": "Test service weather 11-1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-8"
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
      "name": "test-pbehavior-weather-service-8-2",
      "output_template": "Test service weather 11-2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-8-not-exist"
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
        "event_type": "check",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      },
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
      "name": "test-pbehavior-weather-service-8",
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
              "value": "test-pbehavior-weather-service-8-1"
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
    When I send an event:
    """json
    {
      "connector":  "test-connector-pbehavior-weather-service-8",
      "connector_name": "test-connector-name-pbehavior-weather-service-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-8",
      "resource": "test-resource-pbehavior-weather-service-8",
      "state": 3,
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "component":  "test-component-pbehavior-weather-service-8",
        "resource": "test-resource-pbehavior-weather-service-8",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-8",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-8-1",
                "test-pbehavior-weather-service-8-2"
              ]
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
      "data": [
        {
          "name": "test-pbehavior-weather-service-8-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-8",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-8-1",
                "test-pbehavior-weather-service-8-2"
              ]
            }
          }
        ]
      ],
      "weather_service_pattern": [
        [
          {
            "field": "icon",
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
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-8-2"
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
  Scenario: given one dependency with maintenance pbehavior should get service by filter secondary_icon=maintenance
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-9",
      "connector_name": "test-connector-name-pbehavior-weather-service-9",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-weather-service-9",
      "resource": "test-resource-pbehavior-weather-service-9",
      "state": 2,
      "output": "noveo alarm"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-9-1",
      "output_template": "Test service weather 12-1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-9"
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
      "name": "test-pbehavior-weather-service-9-2",
      "output_template": "Test service weather 12-2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-9-not-exist"
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
      "name": "test-pbehavior-weather-service-9-1",
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
              "value": "test-resource-pbehavior-weather-service-9"
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
        "component": "test-component-pbehavior-weather-service-9",
        "resource": "test-resource-pbehavior-weather-service-9",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-9",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-9-1",
                "test-pbehavior-weather-service-9-2"
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
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-9-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-9",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-9-1",
                "test-pbehavior-weather-service-9-2"
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
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-9-2"
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
