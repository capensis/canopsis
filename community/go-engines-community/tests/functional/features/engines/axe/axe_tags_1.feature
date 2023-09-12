Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  @concurrent
  Scenario: given check event should create alarm with tags
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "env": "prod",
        "cloud": "",
        "": "localhost"
      },
      "state": 2,
      "output": "test-output-axe-tags-1",
      "connector": "test-connector-axe-tags-1",
      "connector_name": "test-connector-name-axe-tags-1",
      "component": "test-component-axe-tags-1",
      "resource": "test-resource-axe-tags-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-1
    Then the response code should be 200
    Then the response array key "data.0.tags" should contain only:
    """json
    [
      "env: prod",
      "cloud"
    ]
    """

  @concurrent
  Scenario: given check event with another tags should add new tags
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "env": "prod",
        "cloud": "",
        "": "localhost"
      },
      "state": 2,
      "output": "test-output-axe-tags-2",
      "connector": "test-connector-axe-tags-2",
      "connector_name": "test-connector-name-axe-tags-2",
      "component": "test-component-axe-tags-2",
      "resource": "test-resource-axe-tags-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "env": "test",
        "website": "",
        "cloud": "",
        "": ""
      },
      "state": 2,
      "output": "test-output-axe-tags-2",
      "connector": "test-connector-axe-tags-2",
      "connector_name": "test-connector-name-axe-tags-2",
      "component": "test-component-axe-tags-2",
      "resource": "test-resource-axe-tags-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-2
    Then the response code should be 200
    Then the response array key "data.0.tags" should contain only:
    """json
    [
      "env: prod",
      "cloud",
      "env: test",
      "website"
    ]
    """

  @concurrent
  Scenario: given check event without tags should keep old tags
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        " env ": " prod ",
        "cloud": "  ",
        " ": "localhost"
      },
      "state": 2,
      "output": "test-output-axe-tags-3",
      "connector": "test-connector-axe-tags-3",
      "connector_name": "test-connector-name-axe-tags-3",
      "component": "test-component-axe-tags-3",
      "resource": "test-resource-axe-tags-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-tags-3",
      "connector": "test-connector-axe-tags-3",
      "connector_name": "test-connector-name-axe-tags-3",
      "component": "test-component-axe-tags-3",
      "resource": "test-resource-axe-tags-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-3
    Then the response code should be 200
    Then the response array key "data.0.tags" should contain only:
    """json
    [
      "env: prod",
      "cloud"
    ]
    """

  @concurrent
  Scenario: given check event should create alarm tags
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "test-tag-axe-tags-4-1": " prod ",
        " test-tag-axe-tags-4-2 ": " "
      },
      "state": 2,
      "output": "test-output-axe-tags-4",
      "connector": "test-connector-axe-tags-4",
      "connector_name": "test-connector-name-axe-tags-4",
      "component": "test-component-axe-tags-4",
      "resource": "test-resource-axe-tags-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarm-tags?search=test-tag-axe-tags-4 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 0,
          "value": "test-tag-axe-tags-4-1: prod"
        },
        {
          "type": 0,
          "value": "test-tag-axe-tags-4-2"
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
    When I save response sameLabelColor={{ (index .lastResponse.data 0).color }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "test-tag-axe-tags-4-1": " test ",
        " test-tag-axe-tags-4-3 ": " "
      },
      "state": 2,
      "output": "test-output-axe-tags-4",
      "connector": "test-connector-axe-tags-4",
      "connector_name": "test-connector-name-axe-tags-4",
      "component": "test-component-axe-tags-4",
      "resource": "test-resource-axe-tags-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarm-tags?search=test-tag-axe-tags-4 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 0,
          "value": "test-tag-axe-tags-4-1: prod"
        },
        {
          "type": 0,
          "value": "test-tag-axe-tags-4-1: test",
          "color": "{{ .sameLabelColor }}"
        },
        {
          "type": 0,
          "value": "test-tag-axe-tags-4-2"
        },
        {
          "type": 0,
          "value": "test-tag-axe-tags-4-3"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  @concurrent
  Scenario: given check event should create alarm with internal tags
    Given I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-5-1: prod",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-tags-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-5-2",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-tags-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-tags-5",
      "connector": "test-connector-axe-tags-5",
      "connector_name": "test-connector-name-axe-tags-5",
      "component": "test-component-axe-tags-5",
      "resource": "test-resource-axe-tags-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-5 until response code is 200 and response array key "data.0.tags" contains only:
    """json
    [
      "test-tag-axe-tags-5-1: prod",
      "test-tag-axe-tags-5-2"
    ]
    """

  @concurrent
  Scenario: given new internal tags should update alarms
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-6",
        "connector": "test-connector-axe-tags-6",
        "connector_name": "test-connector-name-axe-tags-6",
        "component": "test-component-axe-tags-6",
        "resource": "test-resource-axe-tags-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-6",
        "connector": "test-connector-axe-tags-6",
        "connector_name": "test-connector-name-axe-tags-6",
        "component": "test-component-axe-tags-6",
        "resource": "test-resource-axe-tags-6-2",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-6: prod",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-axe-tags-6-1",
                "test-resource-axe-tags-6-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-6 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "tags": ["test-tag-axe-tags-6: prod"]
        },
        {
          "tags": ["test-tag-axe-tags-6: prod"]
        }
      ]
    }
    """

  @concurrent
  Scenario: given updated internal tags should update alarms
    Given I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-7",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-axe-tags-7-1",
                "test-resource-axe-tags-7-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response tagId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-7",
        "connector": "test-connector-axe-tags-7",
        "connector_name": "test-connector-name-axe-tags-7",
        "component": "test-component-axe-tags-7",
        "resource": "test-resource-axe-tags-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-7",
        "connector": "test-connector-axe-tags-7",
        "connector_name": "test-connector-name-axe-tags-7",
        "component": "test-component-axe-tags-7",
        "resource": "test-resource-axe-tags-7-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-7",
        "connector": "test-connector-axe-tags-7",
        "connector_name": "test-connector-name-axe-tags-7",
        "component": "test-component-axe-tags-7",
        "resource": "test-resource-axe-tags-7-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-7&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "tags": ["test-tag-axe-tags-7"],
          "v": {
            "resource": "test-resource-axe-tags-7-1"
          }
        },
        {
          "tags": ["test-tag-axe-tags-7"],
          "v": {
            "resource": "test-resource-axe-tags-7-2"
          }
        },
        {
          "tags": [],
          "v": {
            "resource": "test-resource-axe-tags-7-3"
          }
        }
      ]
    }
    """
    When I do PUT /api/v4/alarm-tags/{{ .tagId }}:
    """json
    {
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-axe-tags-7-1",
                "test-resource-axe-tags-7-3",
                "test-resource-axe-tags-7-4"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-tags-7",
      "connector": "test-connector-axe-tags-7",
      "connector_name": "test-connector-name-axe-tags-7",
      "component": "test-component-axe-tags-7",
      "resource": "test-resource-axe-tags-7-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-7&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "tags": ["test-tag-axe-tags-7"],
          "v": {
            "resource": "test-resource-axe-tags-7-1"
          }
        },
        {
          "tags": [],
          "v": {
            "resource": "test-resource-axe-tags-7-2"
          }
        },
        {
          "tags": ["test-tag-axe-tags-7"],
          "v": {
            "resource": "test-resource-axe-tags-7-3"
          }
        },
        {
          "tags": ["test-tag-axe-tags-7"],
          "v": {
            "resource": "test-resource-axe-tags-7-4"
          }
        }
      ]
    }
    """

  @concurrent
  Scenario: given removed internal tags should update alarms
    Given I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-8",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-axe-tags-8-1",
                "test-resource-axe-tags-8-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response tagId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-8",
        "connector": "test-connector-axe-tags-8",
        "connector_name": "test-connector-name-axe-tags-8",
        "component": "test-component-axe-tags-8",
        "resource": "test-resource-axe-tags-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-tags-8",
        "connector": "test-connector-axe-tags-8",
        "connector_name": "test-connector-name-axe-tags-8",
        "component": "test-component-axe-tags-8",
        "resource": "test-resource-axe-tags-8-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-8 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "tags": ["test-tag-axe-tags-8"]
        },
        {
          "tags": ["test-tag-axe-tags-8"]
        }
      ]
    }
    """
    When I do DELETE /api/v4/alarm-tags/{{ .tagId }}
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-8 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "tags": []
        },
        {
          "tags": []
        }
      ]
    }
    """

  @concurrent
  Scenario: given check event should create alarm with external and internal tags
    Given I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-9-1: prod",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-tags-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-tag-axe-tags-9-2",
      "color": "#AABBCC",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-tags-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "test-tag-axe-tags-9-1": "test",
        "test-tag-axe-tags-9-3": ""
      },
      "state": 2,
      "output": "test-output-axe-tags-9",
      "connector": "test-connector-axe-tags-9",
      "connector_name": "test-connector-name-axe-tags-9",
      "component": "test-component-axe-tags-9",
      "resource": "test-resource-axe-tags-9",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-9 until response code is 200 and response array key "data.0.tags" contains only:
    """json
    [
      "test-tag-axe-tags-9-1: prod",
      "test-tag-axe-tags-9-1: test",
      "test-tag-axe-tags-9-2",
      "test-tag-axe-tags-9-3"
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "test-tag-axe-tags-9-1": "test",
        "test-tag-axe-tags-9-3": ""
      },
      "state": 2,
      "output": "test-output-axe-tags-9",
      "connector": "test-connector-axe-tags-9",
      "connector_name": "test-connector-name-axe-tags-9",
      "component": "test-component-axe-tags-9",
      "resource": "test-resource-axe-tags-9",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "tags": {
        "test-tag-axe-tags-9-4": ""
      },
      "state": 2,
      "output": "test-output-axe-tags-9",
      "connector": "test-connector-axe-tags-9",
      "connector_name": "test-connector-name-axe-tags-9",
      "component": "test-component-axe-tags-9",
      "resource": "test-resource-axe-tags-9",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-tags-9 until response code is 200 and response array key "data.0.tags" contains only:
    """json
    [
      "test-tag-axe-tags-9-1: prod",
      "test-tag-axe-tags-9-1: test",
      "test-tag-axe-tags-9-2",
      "test-tag-axe-tags-9-3",
      "test-tag-axe-tags-9-4"
    ]
    """
