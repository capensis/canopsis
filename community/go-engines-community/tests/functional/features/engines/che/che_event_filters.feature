Feature: modify event on event filter
  I need to be able to modify event on event filter


  Scenario: given check event and drop event filter with time interval shouldn't drop event because of exception
    Given I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "_id": "test-eventfilter-20",
      "name": "test-eventfilter-20-name",
      "description": "test-eventfilter-20-description",
      "exdates":[
        {
          "begin": {{ now }},
          "end": {{ nowAdd "6s" }},
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
            "value": "test-component-che-event-filters-20"
          }
        }
      ]],
      "exceptions": [
        "test-eventfilter-20"
      ],
      "description": "test-event-filter-che-event-filters-20-description",
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-event-filters-20",
      "connector_name": "test-connector-name-che-event-filters-20",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-20",
      "resource": "test-resource-che-event-filters-20-1",
      "state": 2,
      "output": "test-output-che-event-filters-20"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=che-event-filters-20
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-20-1/test-component-che-event-filters-20"
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
    When I send an event:
    """json
    {
      "connector": "test-connector-che-event-filters-20",
      "connector_name": "test-connector-name-che-event-filters-20",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-20",
      "resource": "test-resource-che-event-filters-20-2",
      "state": 2,
      "output": "test-output-che-event-filters-20"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=che-event-filters-20
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-20-1/test-component-che-event-filters-20"
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

  Scenario: given check event and drop event filter shouldn't drop event after update request because of interval
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "_id": "test-eventfilter-21",
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
            "value": "test-component-che-event-filters-21"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-21-description",
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-event-filters-21",
      "connector_name": "test-connector-name-che-event-filters-21",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-21",
      "resource": "test-resource-che-event-filters-21-1",
      "state": 2,
      "output": "test-output-che-event-filters-21"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=che-event-filters-21
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
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-21:
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
            "value": "test-component-che-event-filters-21"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-21-description",
      "priority": 1,
      "enabled": true
    }
    """
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-event-filters-21",
      "connector_name": "test-connector-name-che-event-filters-21",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-21",
      "resource": "test-resource-che-event-filters-21-2",
      "state": 2,
      "output": "test-output-che-event-filters-21"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=che-event-filters-21
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-21-2/test-component-che-event-filters-21"
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

  Scenario: given check event and drop event filter should drop event after update request because of removed exdate
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "_id": "test-eventfilter-22",
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
            "value": "test-component-che-event-filters-22"
          }
        }
      ]],
      "exdates": [
        {
          "begin": {{ now }},
          "end": {{ nowAdd "30s"}}
        }
      ],
      "description": "test-event-filter-che-event-filters-22-description",
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-che-event-filters-22",
      "connector_name": "test-connector-name-che-event-filters-22",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-22",
      "resource": "test-resource-che-event-filters-22-1",
      "state": 2,
      "output": "test-output-che-event-filters-22"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=che-event-filters-22
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-22-1/test-component-che-event-filters-22"
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
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-22:
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
            "value": "test-component-che-event-filters-22"
          }
        }
      ]],
      "description": "test-event-filter-che-event-filters-22-description",
      "priority": 1,
      "enabled": true
    }
    """
    When I wait 5s
    When I send an event:
    """json
    {
      "connector": "test-connector-che-event-filters-22",
      "connector_name": "test-connector-name-che-event-filters-22",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-event-filters-22",
      "resource": "test-resource-che-event-filters-22-2",
      "state": 2,
      "output": "test-output-che-event-filters-22"
    }
    """
    When I wait the next periodical process
    When I do GET /api/v4/alarms?search=che-event-filters-22
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-che-event-filters-22-1/test-component-che-event-filters-22"
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
