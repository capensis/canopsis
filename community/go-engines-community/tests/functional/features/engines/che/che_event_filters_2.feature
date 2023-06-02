Feature: modify event on event filter
  I need to be able to modify event on event filter

  @concurrent
  Scenario: given check event and drop event filter with several event patterns should drop event
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "entity_pattern": [[
        {
          "field": "name",
          "cond": {
            "type": "regexp",
            "value": "CMDB:(?P<SI_CMDB>.*?)($|,)"
          }
        }
      ]],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "test_template",
            "description": "test template",
            "value": "{{ `{{ .RegexMatch.Entity.Name.SI_CMDB }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "description": "test-event-filter-che-event-filters-second-1-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test_connector",
      "connector_name": "test_connector_name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-1",
      "resource": "CMDB:TEST_PROD",
      "state": 2
    }
    """
    When I do GET /api/v4/entitybasics?_id=CMDB:TEST_PROD/test-component-che-event-filters-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "CMDB:TEST_PROD/test-component-che-event-filters-second-1",
      "infos": {
        "test_template": {
          "name": "test_template",
          "description": "test template",
          "value": "TEST_PROD"
        }
      }
    }
    """

  @concurrent
  Scenario: given rule with old patterns format, backward compatibility test
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-eventfilter-to-backward-compatibility-1",
      "connector_name": "test-connector-name-eventfilter-to-backward-compatibility-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-eventfilter-to-backward-compatibility-1",
      "resource": "CMDD:TEST_PROD",
      "state": 2,
      "output": "test-output-eventfilter-to-backward-compatibility-1",
      "customer": "test-customer-eventfilter-to-backward-compatibility-1",
      "manager": "test-manager-eventfilter-to-backward-compatibility-1"
    }
    """
    When I do GET /api/v4/entitybasics?_id=CMDD:TEST_PROD/test-component-eventfilter-to-backward-compatibility-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "infos": {
        "customer": {
          "description": "customer",
          "name": "customer",
          "value": "TEST_PROD"
        },
        "manager": {
          "description": "manager",
          "name": "manager",
          "value": "TEST_PROD"
        }
      }
    }
    """

  @concurrent
  Scenario: given check event and drop event filter with time interval should drop event in that interval
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "start": {{ nowAdd "4s" }},
      "stop": {{ nowAdd "10s" }},
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-3"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-second-3-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-3",
      "connector_name": "test-connector-name-che-event-filters-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-3",
      "resource": "test-resource-che-event-filters-second-3-1",
      "state": 2,
      "output": "test-output-che-event-filters-second-3"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-3-1/test-component-che-event-filters-second-3"
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
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-3",
      "connector_name": "test-connector-name-che-event-filters-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-3",
      "resource": "test-resource-che-event-filters-second-3-2",
      "state": 2,
      "output": "test-output-che-event-filters-second-3"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-3-1/test-component-che-event-filters-second-3"
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
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-3",
      "connector_name": "test-connector-name-che-event-filters-second-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-3",
      "resource": "test-resource-che-event-filters-second-3-3",
      "state": 2,
      "output": "test-output-che-event-filters-second-3"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-3&sort_by=t&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-3-1/test-component-che-event-filters-second-3"
          }
        },
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-3-3/test-component-che-event-filters-second-3"
          }
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

  @concurrent
  Scenario: given check event and drop event filter with time interval should drop event by rrule
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "start": {{ now }},
      "stop": {{ nowAdd "5s" }},
      "rrule": "FREQ=SECONDLY",
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-4"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-second-4-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-4",
      "connector_name": "test-connector-name-che-event-filters-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-4",
      "resource": "test-resource-che-event-filters-second-4-1",
      "state": 2,
      "output": "test-output-che-event-filters-second-4"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-4
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
  Scenario: given check event and drop event filter with time interval should drop event by rrule with count
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "start": {{ nowAdd "-12s" }},
      "stop": {{ nowAdd "5s" }},
      "rrule": "FREQ=SECONDLY;COUNT=2",
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-5"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-second-5-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-5",
      "connector_name": "test-connector-name-che-event-filters-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-5",
      "resource": "test-resource-che-event-filters-second-5-1",
      "state": 2,
      "output": "test-output-che-event-filters-second-5"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-5
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
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-5",
      "connector_name": "test-connector-name-che-event-filters-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-5",
      "resource": "test-resource-che-event-filters-second-5-2",
      "state": 2,
      "output": "test-output-che-event-filters-second-5"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-5-2/test-component-che-event-filters-second-5"
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
  Scenario: given check event and drop event filter with time interval shouldn't drop event because of exdate
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "start": {{ now }},
      "stop": {{ nowAdd "1m" }},
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-6"
          }
        }
      ]],
      "exdates": [
        {
          "begin": {{ now }},
          "end": {{ nowAdd "8s" }}
        }
      ],
      "description": "test-event-filter-che-event-filters-second-6-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-6",
      "connector_name": "test-connector-name-che-event-filters-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-6",
      "resource": "test-resource-che-event-filters-second-6-1",
      "state": 2,
      "output": "test-output-che-event-filters-second-6"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-6-1/test-component-che-event-filters-second-6"
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
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-6",
      "connector_name": "test-connector-name-che-event-filters-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-6",
      "resource": "test-resource-che-event-filters-second-6-2",
      "state": 2,
      "output": "test-output-che-event-filters-second-6"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-6-1/test-component-che-event-filters-second-6"
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
  Scenario: given check event and drop event filter with time interval shouldn't drop event because of exception
    Given I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """json
    {
      "_id": "test-eventfilter-22",
      "name": "test-eventfilter-22-name",
      "description": "test-eventfilter-22-description",
      "exdates":[
        {
          "begin": {{ now }},
          "end": {{ nowAdd "8s" }},
          "type": "test-type-to-exception-edit-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "start": {{ now }},
      "stop": {{ nowAdd "1m" }},
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-7"
          }
        }
      ]],
      "exceptions": [
        "test-eventfilter-22"
      ],
      "description": "test-event-filter-che-event-filters-second-7-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-7",
      "connector_name": "test-connector-name-che-event-filters-second-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-7",
      "resource": "test-resource-che-event-filters-second-7-1",
      "state": 2,
      "output": "test-output-che-event-filters-second-7"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-7-1/test-component-che-event-filters-second-7"
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
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-7",
      "connector_name": "test-connector-name-che-event-filters-second-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-7",
      "resource": "test-resource-che-event-filters-second-7-2",
      "state": 2,
      "output": "test-output-che-event-filters-second-7"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-7-1/test-component-che-event-filters-second-7"
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
  Scenario: given check event and drop event filter shouldn't drop event after update request because of interval
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-8"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-second-8-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    Then I save response eventFilterId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-8",
      "connector_name": "test-connector-name-che-event-filters-second-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-8",
      "resource": "test-resource-che-event-filters-second-8-1",
      "state": 2,
      "output": "test-output-che-event-filters-second-8"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-8
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
    When I do PUT /api/v4/eventfilter/rules/{{ .eventFilterId }}:
    """json
    {
      "type": "drop",
      "start": {{ nowAdd "-4s" }},
      "stop": {{ now }},
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-8"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-second-8-description",
      "enabled": true
    }
    """
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-8",
      "connector_name": "test-connector-name-che-event-filters-second-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-8",
      "resource": "test-resource-che-event-filters-second-8-2",
      "state": 2,
      "output": "test-output-che-event-filters-second-8"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-8-2/test-component-che-event-filters-second-8"
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
  Scenario: given check event and drop event filter should drop event after update request because of removed exdate
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "start": {{ now }},
      "stop": {{ nowAdd "1m"}},
      "type": "drop",
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-9"
          }
        }
      ]],
      "exdates": [
        {
          "begin": {{ now }},
          "end": {{ nowAdd "30s"}}
        }
      ],
      "description": "test-event-filter-che-event-filters-second-9-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    Then I save response eventFilterId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-9",
      "connector_name": "test-connector-name-che-event-filters-second-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-9",
      "resource": "test-resource-che-event-filters-second-9-1",
      "state": 2,
      "output": "test-output-che-event-filters-second-9"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-9-1/test-component-che-event-filters-second-9"
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
    When I do PUT /api/v4/eventfilter/rules/{{ .eventFilterId }}:
    """json
    {
      "type": "drop",
      "start": {{ now }},
      "stop": {{ nowAdd "1m"}},
      "event_pattern":[[
        {
          "field": "event_type",
          "cond": {
            "type": "eq",
            "value": "check"
          }
        },
        {
          "field": "component",
          "cond": {
            "type": "eq",
            "value": "test-component-che-event-filters-second-9"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-second-9-description",
      "priority": 1,
      "enabled": true
    }
    """
    When I wait 5s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-event-filters-second-9",
      "connector_name": "test-connector-name-che-event-filters-second-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-second-9",
      "resource": "test-resource-che-event-filters-second-9-2",
      "state": 2,
      "output": "test-output-che-event-filters-second-9"
    }
    """
    When I do GET /api/v4/alarms?search=che-event-filters-second-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-second-9-1/test-component-che-event-filters-second-9"
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
