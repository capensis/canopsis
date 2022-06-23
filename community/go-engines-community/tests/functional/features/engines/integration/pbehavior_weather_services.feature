Feature: get service weather
  I need to be able to get service weather

  Scenario: given one dependency with maintenance pbehavior should get maintenance icon and grey flag
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-1",
      "connector_name": "test-connector-name-pbehavior-weather-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-1",
      "resource" : "test-resource-pbehavior-weather-service-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
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
    When I wait the end of 2 events processing
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-1",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given one dependency with active pbehavior should get state icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-2",
      "connector_name": "test-connector-name-pbehavior-weather-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-2",
      "resource" : "test-resource-pbehavior-weather-service-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
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
    When I wait the end of event processing
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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

  Scenario: given dependency with maintenance pbehavior
    and dependency without pbehavior should get secondary maintenance icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-3",
      "resource" : "test-resource-pbehavior-weather-service-3-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-3",
      "resource" : "test-resource-pbehavior-weather-service-3-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-3-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
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
    When I wait the end of 2 events processing
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-3",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

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
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-4",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
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
    When I save response pbehaviorID={{ .lastResponse._id }}
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "{{ .pbehaviorID }}",
      "message": "First comment"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "{{ .pbehaviorID }}",
      "message": "Second comment"
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-4",
      "connector_name": "test-connector-name-pbehavior-weather-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-4",
      "resource" : "test-resource-pbehavior-weather-service-4",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-4",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [],
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
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
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-5",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-5",
      "resource" : "test-resource-pbehavior-weather-service-5",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-5",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "critical",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": [],
          "pbehavior_info": {
            "canonical_type": "active",
            "name": "test-pbehavior-weather-service-5",
            "icon_name": "brightness_3"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-5"
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

  Scenario: given service with maintenance pbehavior and one dependency with maintenance pbehavior
    should get maintenance icon and not get secondary icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-6",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-6",
      "resource" : "test-resource-pbehavior-weather-service-6",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-6-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
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
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-6-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
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
    When I wait the end of 2 events processing
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
          ],
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-6-1",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-6-1"
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

  Scenario: given service with maintenance pbehavior and dependency with maintenance pbehavior
    and another dependency without pbehavior should get maintenance icon and maintenance secondary icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-7",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-7",
      "resource" : "test-resource-pbehavior-weather-service-7-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-7",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-7",
      "resource" : "test-resource-pbehavior-weather-service-7-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-7-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
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
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-7-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
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
    When I wait the end of 2 events processing
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-7",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "maintenance",
          "secondary_icon": "maintenance",
          "is_grey": true,
          "alarm_counters": [
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
          ],
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-7-1",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-7-1"
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

  Scenario: given service with maintenance pbehavior should get service by filter icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-11-1",
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
              "value": "test-resource-pbehavior-weather-service-11"
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
      "name": "test-pbehavior-weather-service-11-2",
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
              "value": "test-resource-pbehavior-weather-service-11-not-exist"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-11",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-11-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-11",
      "connector_name": "test-connector-name-pbehavior-weather-service-11",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-11",
      "resource" : "test-resource-pbehavior-weather-service-11",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-11",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-11-1",
                "test-pbehavior-weather-service-11-2"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-11-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-11",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-11-1",
                "test-pbehavior-weather-service-11-2"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-11-2"
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

  Scenario: given one dependency with maintenance pbehavior should get service by filter icon=maintenance
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-12",
      "connector_name": "test-connector-name-pbehavior-weather-service-12",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-12",
      "resource" : "test-resource-pbehavior-weather-service-12",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-12-1",
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
              "value": "test-resource-pbehavior-weather-service-12"
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
      "name": "test-pbehavior-weather-service-12-2",
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
              "value": "test-resource-pbehavior-weather-service-12-not-exist"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-12-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-12"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-12",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-12-1",
                "test-pbehavior-weather-service-12-2"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-12-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-12",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-12-1",
                "test-pbehavior-weather-service-12-2"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-12-2"
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

  Scenario: given service without pbehavior should not get service by filter icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-13",
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
              "value": "test-resource-pbehavior-weather-service-13-1"
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
      "connector" :  "test-connector-pbehavior-weather-service-13",
      "connector_name": "test-connector-name-pbehavior-weather-service-13",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-13",
      "resource" : "test-resource-pbehavior-weather-service-13-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-13",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-13"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
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

  Scenario: given dependency with maintenance pbehavior
    and another dependency without pbehavior should get service by filter secondary_icon=maintenance
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-14",
      "connector_name": "test-connector-name-pbehavior-weather-service-14",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-14",
      "resource" : "test-resource-pbehavior-weather-service-14-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-14",
      "connector_name": "test-connector-name-pbehavior-weather-service-14",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-14",
      "resource" : "test-resource-pbehavior-weather-service-14-2",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-14-1",
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
                "test-resource-pbehavior-weather-service-14-1",
                "test-resource-pbehavior-weather-service-14-2"
              ]
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
      "name": "test-pbehavior-weather-service-14-2",
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
                "test-resource-pbehavior-weather-service-14-1"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-14-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-14-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 3 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-14",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-14-1",
                "test-pbehavior-weather-service-14-2"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-14-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-14",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-14-1",
                "test-pbehavior-weather-service-14-2"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-14-2"
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

  Scenario: given dependency without pbehavior should not get service by filter secondary_icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-15",
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
              "value": "test-resource-pbehavior-weather-service-15-1"
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
      "connector" :  "test-connector-pbehavior-weather-service-15",
      "connector_name": "test-connector-name-pbehavior-weather-service-15",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-15",
      "resource" : "test-resource-pbehavior-weather-service-15-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-15",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-15"
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
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
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

  Scenario: given service with maintenance pbehavior without alarm should get maintenance icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-16",
      "connector_name": "test-connector-name-pbehavior-weather-service-16",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-16",
      "resource" : "test-resource-pbehavior-weather-service-16",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-16",
      "name": "test-pbehavior-weather-service-16",
      "output_template": "test-pbehavior-weather-service-16",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-16"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-16",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-16"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-16",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-16"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-16",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [],
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-16",
            "icon_name": "build"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-16"
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

  Scenario: given dependency with maintenance pbehavior without alarm should get maintenance icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-17",
      "connector_name": "test-connector-name-pbehavior-weather-service-17",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-17",
      "resource" : "test-resource-pbehavior-weather-service-17",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-17",
      "name": "test-pbehavior-weather-service-17",
      "output_template": "test-pbehavior-weather-service-17",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-17"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-17",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-17"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-17",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-17"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-17",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given dependencies with maintenance pbehavior should keep maintenance icon on alarm create and resolve
    Given I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-pbehavior-weather-service-18",
      "name": "test-resolve-rule-pbehavior-weather-service-18-name",
      "entity_pattern":[
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-18-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-18-2"
            }
          }
        ]
      ],
      "duration": {
        "value": 1,
        "unit": "s"
      },
      "priority": 10
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-pbehavior-weather-service-18",
        "connector_name": "test-connector-name-pbehavior-weather-service-18",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-18",
        "resource" : "test-resource-pbehavior-weather-service-18-1",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-pbehavior-weather-service-18",
        "connector_name": "test-connector-name-pbehavior-weather-service-18",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-18",
        "resource" : "test-resource-pbehavior-weather-service-18-2",
        "state" : 0,
        "output" : "noveo alarm"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-18",
      "name": "test-pbehavior-weather-service-18",
      "output_template": "test-pbehavior-weather-service-18",
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
                "test-resource-pbehavior-weather-service-18-1",
                "test-resource-pbehavior-weather-service-18-2"
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-18",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-pbehavior-weather-service-18-1",
                "test-resource-pbehavior-weather-service-18-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-18",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-18"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      "connector" : "test-connector-pbehavior-weather-service-18",
      "connector_name": "test-connector-name-pbehavior-weather-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-18",
      "resource" : "test-resource-pbehavior-weather-service-18-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      "connector" : "test-connector-pbehavior-weather-service-18",
      "connector_name": "test-connector-name-pbehavior-weather-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-18",
      "resource" : "test-resource-pbehavior-weather-service-18-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      "connector" : "test-connector-pbehavior-weather-service-18",
      "connector_name": "test-connector-name-pbehavior-weather-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-18",
      "resource" : "test-resource-pbehavior-weather-service-18-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given dependency with maintenance pbehavior should update service state correctly on alarm state change
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-19",
      "connector_name": "test-connector-name-pbehavior-weather-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-19",
      "resource" : "test-resource-pbehavior-weather-service-19",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-19",
      "name": "test-pbehavior-weather-service-19",
      "output_template": "test-pbehavior-weather-service-19",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-19"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-19",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-19"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response pbhID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-19",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-19"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      "connector" : "test-connector-pbehavior-weather-service-19",
      "connector_name": "test-connector-name-pbehavior-weather-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-19",
      "resource" : "test-resource-pbehavior-weather-service-19",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do DELETE /api/v4/pbehaviors/{{ .pbhID }}
    Then the response code should be 204
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "ok",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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
      "connector" : "test-connector-pbehavior-weather-service-19",
      "connector_name": "test-connector-name-pbehavior-weather-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-19",
      "resource" : "test-resource-pbehavior-weather-service-19",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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

  Scenario: given service with pbehavior should get service by filter is_grey=true and should not get service by filter is_grey=false
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-20-1",
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
              "value": "test-resource-pbehavior-weather-service-20"
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
      "name": "test-pbehavior-weather-service-20-2",
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
              "value": "test-resource-pbehavior-weather-service-20"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-20",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-weather-service-20-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-20",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-20-1",
                "test-pbehavior-weather-service-20-2"
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-20-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-20",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-20-1",
                "test-pbehavior-weather-service-20-2"
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-20-2"
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

  Scenario: given dependency with pbehavior should get service by filter is_grey=true and should not get service by filter is_grey=false
    Given I am admin
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-pbehavior-weather-service-21",
        "connector_name": "test-connector-name-pbehavior-weather-service-21",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-21",
        "resource" : "test-resource-pbehavior-weather-service-21-1",
        "state" : 2,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-pbehavior-weather-service-21",
        "connector_name": "test-connector-name-pbehavior-weather-service-21",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-21",
        "resource" : "test-resource-pbehavior-weather-service-21-2",
        "state" : 2,
        "output" : "noveo alarm"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-21-1",
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
              "value": "test-resource-pbehavior-weather-service-21-1"
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
      "name": "test-pbehavior-weather-service-21-2",
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
                "test-resource-pbehavior-weather-service-21-1",
                "test-resource-pbehavior-weather-service-21-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-21",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-weather-service-21-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 3 events processing
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-weather-service-21",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-21-1",
                "test-pbehavior-weather-service-21-2"
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-21-1"
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
      "title": "test-widgetfilter-pbehavior-weather-service-21",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-pbehavior-weather-service-21-1",
                "test-pbehavior-weather-service-21-2"
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
    When I do GET /api/v4/weather-services?filter={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-21-2"
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
