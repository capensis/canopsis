Feature: modify event on fifo event filter
  I need to be able to modify event on fifo event filter

  @concurrent
  Scenario: given check event and change_entity rule should change entity defining fields by value
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "description": "test",
      "type": "change_entity",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-fifo-event-filters-1"
            }
          },
          {
            "field": "extra.customer_tags",
            "field_type": "string",
            "cond": {
              "type": "regexp",
              "value": "CMDB:(?P<SI_CMDB>.*?)($|,)"
            }
          }
        ]
      ],
      "config": {
        "component": "test-new-component-fifo-event-filters-1"
      },
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "customer_tags": "CMDB:TEST_PROD",
      "state": 2,
      "output": "test-output-fifo-event-filters-1",
      "connector": "test-connector-fifo-event-filters-1",
      "connector_name": "test-connector-name-fifo-event-filters-1",
      "component": "test-component-fifo-event-filters-1",
      "resource": "test-resource-fifo-event-filters-1",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-fifo-event-filters-1",
      "connector_name": "test-connector-name-fifo-event-filters-1",
      "component": "test-new-component-fifo-event-filters-1",
      "resource": "test-resource-fifo-event-filters-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-fifo-event-filters-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-fifo-event-filters-1",
          "component": "test-new-component-fifo-event-filters-1"
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
  Scenario: given check event and change_entity rule should change entity defining fields by template
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "description": "test",
      "type": "change_entity",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-fifo-event-filters-2"
            }
          },
          {
            "field": "extra.customer_tags",
            "field_type": "string",
            "cond": {
              "type": "regexp",
              "value": "CMDB:(?P<SI_CMDB>.*?)($|,)"
            }
          }
        ]
      ],
      "config": {
        "component": "{{ `{{.RegexMatch.ExtraInfos.customer_tags.SI_CMDB}}` }}"
      },
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "event_type": "check",
      "customer_tags": "CMDB:TEST_PROD",
      "state": 2,
      "output": "test-output-fifo-event-filters-2",
      "connector": "test-connector-fifo-event-filters-2",
      "connector_name": "test-connector-name-fifo-event-filters-2",
      "component": "test-component-fifo-event-filters-2",
      "resource": "test-resource-fifo-event-filters-2",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-fifo-event-filters-2",
      "connector_name": "test-connector-name-fifo-event-filters-2",
      "component": "TEST_PROD",
      "resource": "test-resource-fifo-event-filters-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities?search=test-resource-fifo-event-filters-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-fifo-event-filters-2",
          "component": "TEST_PROD"
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
  Scenario: given check event and change_entity rule should update events count
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-3",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-connector-fifo-event-filters-3"
            }
          }
        ]
      ],
      "config": {
        "component": "test-new-component-fifo-event-filters-3"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "output": "test-output-fifo-event-filters-3",
        "state": 2,
        "connector": "test-connector-fifo-event-filters-3",
        "connector_name": "test-connector-name-fifo-event-filters-3",
        "component": "test-component-fifo-event-filters-3",
        "resource": "test-resource-fifo-event-filters-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-3",
        "connector": "test-connector-fifo-event-filters-3",
        "connector_name": "test-connector-name-fifo-event-filters-3",
        "component": "test-component-fifo-event-filters-3",
        "resource": "test-resource-fifo-event-filters-3",
        "source_type": "resource"
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-fifo-event-filters-3",
        "connector_name": "test-connector-name-fifo-event-filters-3",
        "component": "test-new-component-fifo-event-filters-3",
        "resource": "test-resource-fifo-event-filters-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "test-connector-fifo-event-filters-3",
        "connector_name": "test-connector-name-fifo-event-filters-3",
        "component": "test-new-component-fifo-event-filters-3",
        "resource": "test-resource-fifo-event-filters-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules?search=fifo-event-filters-3 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "description": "test-event-filter-fifo-event-filters-3",
          "events_count": 2
        }
      ]
    }
    """
