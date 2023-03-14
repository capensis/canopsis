Feature: get service weather
  I need to be able to get service weather

  @concurrent
  Scenario: given entity service entity should decrease pbh counter on disabled and increase on enabled actions
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-third-1",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-1",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-1",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-2",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-1",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-3",
        "state": 2,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-third-1",
      "output_template": "test-pbehavior-weather-service-third-1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-resource-pbehavior-weather-service-third-1"
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
      "name": "test-pbehavior-weather-service-third-1",
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
              "value": "test-resource-pbehavior-weather-service-third-1-2"
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
        "component": "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-2",
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
      "title": "test-widgetfilter-pbehavior-weather-service-third-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-third-1"
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
          "name": "test-pbehavior-weather-service-third-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
            "all": 3,
            "active": 2,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 2,
              "critical": 0
            },
            "acked": 0,
            "unacked": 2,
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
                  "priority": 19,
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-third-1-2/test-component-pbehavior-weather-service-third-1:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-2",
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
          "name": "test-pbehavior-weather-service-third-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "",
          "is_grey": false,
          "counters": {
            "depends": 2,
            "all": 2,
            "active": 2,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 2,
              "critical": 0
            },
            "acked": 0,
            "unacked": 2,
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-third-1-2/test-component-pbehavior-weather-service-third-1:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-2",
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
          "name": "test-pbehavior-weather-service-third-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
            "all": 2,
            "active": 2,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 2,
              "critical": 0
            },
            "acked": 0,
            "unacked": 2,
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
                  "priority": 19,
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-third-1-2/test-component-pbehavior-weather-service-third-1:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-2",
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
          "name": "test-pbehavior-weather-service-third-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "",
          "is_grey": false,
          "counters": {
            "depends": 2,
            "all": 2,
            "active": 2,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 2,
              "critical": 0
            },
            "acked": 0,
            "unacked": 2,
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-third-1-2/test-component-pbehavior-weather-service-third-1:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-1",
        "resource": "test-resource-pbehavior-weather-service-third-1-2",
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
          "name": "test-pbehavior-weather-service-third-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
            "all": 2,
            "active": 2,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 2,
              "critical": 0
            },
            "acked": 0,
            "unacked": 2,
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
                  "priority": 19,
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
  Scenario: given entity service entity should be grey if active entity is disabled, while other entities are in pbh state, if disabled entity is enabled again, the entity service should be returned to the ok state.
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-third-2",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-2",
        "resource": "test-resource-pbehavior-weather-service-third-2-1",
        "state": 0,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-2",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-2",
        "resource": "test-resource-pbehavior-weather-service-third-2-2",
        "state": 0,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-2",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-2",
        "resource": "test-resource-pbehavior-weather-service-third-2-3",
        "state": 0,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-third-2",
      "output_template": "test-pbehavior-weather-service-third-2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-resource-pbehavior-weather-service-third-2"
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
      "name": "test-pbehavior-weather-service-third-2-1",
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
              "value": "test-resource-pbehavior-weather-service-third-2-1"
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
        "component": "test-component-pbehavior-weather-service-third-2",
        "resource": "test-resource-pbehavior-weather-service-third-2-1",
        "source_type": "resource"
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
      "name": "test-pbehavior-weather-service-third-2-3",
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
              "value": "test-resource-pbehavior-weather-service-third-2-3"
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
        "component": "test-component-pbehavior-weather-service-third-2",
        "resource": "test-resource-pbehavior-weather-service-third-2-3",
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
      "title": "test-widgetfilter-pbehavior-weather-service-third-2",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-third-2"
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
          "name": "test-pbehavior-weather-service-third-2",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
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
                  "priority": 19,
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-third-2-2/test-component-pbehavior-weather-service-third-2:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "component": "test-component-pbehavior-weather-service-third-2",
      "resource": "test-resource-pbehavior-weather-service-third-2-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-third-2",
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
                  "priority": 19,
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-third-2-2/test-component-pbehavior-weather-service-third-2:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "component": "test-component-pbehavior-weather-service-third-2",
      "resource": "test-resource-pbehavior-weather-service-third-2-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-third-2",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
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
                  "priority": 19,
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
  Scenario: given entity service entity should decrease pbh counter on mass disabled and increase on mass enabled actions
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-third-3",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-1",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-3",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-2",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-3",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-3",
        "state": 2,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-third-3",
      "output_template": "test-pbehavior-weather-service-third-3",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-resource-pbehavior-weather-service-third-3"
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
      "name": "test-pbehavior-weather-service-third-3-1",
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
              "value": "test-resource-pbehavior-weather-service-third-3-1"
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
        "component": "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-1",
        "source_type": "resource"
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
      "name": "test-pbehavior-weather-service-third-3-2",
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
              "value": "test-resource-pbehavior-weather-service-third-3-2"
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
        "component": "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-2",
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
      "title": "test-widgetfilter-pbehavior-weather-service-third-3",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-third-3"
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
          "name": "test-pbehavior-weather-service-third-3",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
            "all": 3,
            "active": 1,
            "state": {
              "ok": 2,
              "minor": 0,
              "major": 1,
              "critical": 0
            },
            "acked": 0,
            "unacked": 1,
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
                  "priority": 19,
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-third-3-1/test-component-pbehavior-weather-service-third-3"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-third-3-2/test-component-pbehavior-weather-service-third-3"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-2",
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
          "name": "test-pbehavior-weather-service-third-3",
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
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-third-3-1/test-component-pbehavior-weather-service-third-3"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-third-3-2/test-component-pbehavior-weather-service-third-3"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-3",
        "resource": "test-resource-pbehavior-weather-service-third-3-2",
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
          "name": "test-pbehavior-weather-service-third-3",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "person",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
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
            "under_pbh": 2,
            "pbh_types": [
              {
                "count": 2,
                "type": {
                  "_id": "test-maintenance-type-to-engine",
                  "description": "Engine maintenance",
                  "icon_name": "build",
                  "name": "Engine maintenance",
                  "priority": 19,
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
  Scenario: given entity service entity should be grey if active entities are disabled by bulk disable action, while other entities are in pbh state, if disabled entities are enabled again, the entity service should be returned to the ok state.
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-third-4",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-4",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-1",
        "state": 0,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-4",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-4",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-2",
        "state": 0,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-4",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-4",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-3",
        "state": 0,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-third-4",
      "output_template": "test-pbehavior-weather-service-third-4",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-resource-pbehavior-weather-service-third-4"
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
      "name": "test-pbehavior-weather-service-third-4-1",
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
              "value": "test-resource-pbehavior-weather-service-third-4-1"
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
        "component": "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-1",
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
      "title": "test-widgetfilter-pbehavior-weather-service-third-4",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-third-4"
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
          "name": "test-pbehavior-weather-service-third-4",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
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
                  "priority": 19,
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-third-4-2/test-component-pbehavior-weather-service-third-4"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-third-4-3/test-component-pbehavior-weather-service-third-4"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-2",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-3",
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
          "name": "test-pbehavior-weather-service-third-4",
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
                  "priority": 19,
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
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-third-4-2/test-component-pbehavior-weather-service-third-4"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-third-4-3/test-component-pbehavior-weather-service-third-4"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-2",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-4",
        "resource": "test-resource-pbehavior-weather-service-third-4-3",
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
          "name": "test-pbehavior-weather-service-third-4",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "wb_sunny",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
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
                  "priority": 19,
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
  Scenario: given entity service entity should recompute counters properly if entity service is enabled before the dependent entity in mass enable action
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-third-5",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-5",
        "resource": "test-resource-pbehavior-weather-service-third-5-1",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-5",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-5",
        "resource": "test-resource-pbehavior-weather-service-third-5-2",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-5",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-5",
        "resource": "test-resource-pbehavior-weather-service-third-5-3",
        "state": 3,
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-third-5",
      "output_template": "test-pbehavior-weather-service-third-5",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-resource-pbehavior-weather-service-third-5"
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
      "name": "test-pbehavior-weather-service-third-5-1",
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
              "value": "test-resource-pbehavior-weather-service-third-5-1"
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
        "component": "test-component-pbehavior-weather-service-third-5",
        "resource": "test-resource-pbehavior-weather-service-third-5-1",
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
      "title": "test-widgetfilter-pbehavior-weather-service-third-5",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-third-5"
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
          "name": "test-pbehavior-weather-service-third-5",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "wb_cloudy",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
            "all": 3,
            "active": 2,
            "state": {
              "ok": 1,
              "minor": 0,
              "major": 1,
              "critical": 1
            },
            "acked": 0,
            "unacked": 2,
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
                  "priority": 19,
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-third-5-1/test-component-pbehavior-weather-service-third-5"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "component": "test-component-pbehavior-weather-service-third-5",
        "resource": "test-resource-pbehavior-weather-service-third-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "{{ .serviceId }}"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "recomputeentityservice",
      "component": "{{ .serviceId }}",
      "source_type": "service"
    }
    """
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-third-5-1/test-component-pbehavior-weather-service-third-5"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "component": "test-component-pbehavior-weather-service-third-5",
      "resource": "test-resource-pbehavior-weather-service-third-5-1",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "{{ .serviceId }}"
      }
    ]
    """
    Then the response code should be 207
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
    When I do GET /api/v4/weather-services?filters[]={{ .filterId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-third-5",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "wb_cloudy",
          "secondary_icon": "build",
          "is_grey": false,
          "counters": {
            "depends": 3,
            "all": 2,
            "active": 2,
            "state": {
              "ok": 0,
              "minor": 0,
              "major": 1,
              "critical": 1
            },
            "acked": 0,
            "unacked": 2,
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
                  "priority": 19,
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
  Scenario: given one acked dependency with maintenance pbehavior should not update ack counter
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-pbehavior-weather-service-third-6",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-6",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-weather-service-third-6",
        "resource": "test-resource-pbehavior-weather-service-third-6",
        "state": 2,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-weather-service-third-6",
        "connector_name": "test-connector-name-pbehavior-weather-service-third-6",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "test-component-pbehavior-weather-service-third-6",
        "resource": "test-resource-pbehavior-weather-service-third-6",
        "output": "noveo alarm"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-third-6",
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
              "value": "test-resource-pbehavior-weather-service-third-6"
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
      "name": "test-pbehavior-weather-service-third-6",
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
              "value": "test-resource-pbehavior-weather-service-third-6"
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
        "component": "test-component-pbehavior-weather-service-third-6",
        "resource": "test-resource-pbehavior-weather-service-third-6",
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
      "title": "test-widgetfilter-pbehavior-weather-service-third-6",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-third-6"
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
          "name": "test-pbehavior-weather-service-third-6",
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
            "acked_under_pbh": 1,
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
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-third-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-third-6",
      "source_type": "resource",
      "event_type": "ackremove",
      "component":  "test-component-pbehavior-weather-service-third-6",
      "resource": "test-resource-pbehavior-weather-service-third-6",
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ackremove",
        "component": "test-component-pbehavior-weather-service-third-6",
        "resource": "test-resource-pbehavior-weather-service-third-6",
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
          "name": "test-pbehavior-weather-service-third-6",
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
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-weather-service-third-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-third-6",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-pbehavior-weather-service-third-6",
      "resource": "test-resource-pbehavior-weather-service-third-6",
      "output": "noveo alarm"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "component": "test-component-pbehavior-weather-service-third-6",
        "resource": "test-resource-pbehavior-weather-service-third-6",
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
          "name": "test-pbehavior-weather-service-third-6",
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
            "acked_under_pbh": 1,
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
        "component": "test-component-pbehavior-weather-service-third-6",
        "resource": "test-resource-pbehavior-weather-service-third-6",
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
          "name": "test-pbehavior-weather-service-third-6",
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
            "acked": 1,
            "unacked": 0,
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
