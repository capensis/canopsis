Feature: Get alarms
  I need to be able to get a alarms

  Scenario: given alarms in pbehavior should return alarms
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-pbehavior-alarm-get-1",
        "connector_name": "test-connector-name-pbehavior-alarm-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-alarm-get-1",
        "resource": "test-resource-pbehavior-alarm-get-1-1",
        "state": 2,
        "output": "test-output-pbehavior-alarm-get-1"
      },
      {
        "connector": "test-connector-pbehavior-alarm-get-1",
        "connector_name": "test-connector-name-pbehavior-alarm-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-alarm-get-1",
        "resource": "test-resource-pbehavior-alarm-get-1-2",
        "state": 2,
        "output": "test-output-pbehavior-alarm-get-1"
      },
      {
        "connector": "test-connector-pbehavior-alarm-get-1",
        "connector_name": "test-connector-name-pbehavior-alarm-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-alarm-get-1",
        "resource": "test-resource-pbehavior-alarm-get-1-3",
        "state": 2,
        "output": "test-output-pbehavior-alarm-get-1"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-pbehavior-alarm-get-1-1",
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
              "value": "test-resource-pbehavior-alarm-get-1-2"
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
      "name": "test-pbehavior-pbehavior-alarm-get-1-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine-2",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-get-1-3"
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
      "title": "test-widgetfilter-pbehavior-alarm-get-1-1",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-alarm-get-1"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
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
    When I do GET /api/v4/alarms?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-2"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-alarm-get-1-2",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-alarm-get-1"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
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
    When I do GET /api/v4/alarms?filters[]={{ .lastResponse._id }}&sort=asc&sort_by=v.resource
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-1"
          }
        },
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-3"
          }
        }
      ],
      "meta": {
        "total_count": 2
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-alarm-get-1-3",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-alarm-get-1"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
            "cond": {
              "type": "eq",
              "value": "active"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filters[]={{ .lastResponse._id }}&sort=asc&sort_by=v.resource
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-1"
          }
        },
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-3"
          }
        }
      ],
      "meta": {
        "total_count": 2
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-alarm-get-1-4",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-alarm-get-1"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
            "cond": {
              "type": "neq",
              "value": "active"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-2"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-alarm-get-1-5",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-alarm-get-1"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.reason",
            "cond": {
              "type": "eq",
              "value": "test-reason-to-engine"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-2"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-pbehavior-alarm-get-1-6",
      "widget": "test-widget-to-alarm-get",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-alarm-get-1"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.reason",
            "cond": {
              "type": "eq",
              "value": "test-reason-to-engine-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-pbehavior-alarm-get-1",
            "connector_name": "test-connector-name-pbehavior-alarm-get-1",
            "component":  "test-component-pbehavior-alarm-get-1",
            "resource": "test-resource-pbehavior-alarm-get-1-3"
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
