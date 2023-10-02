Feature: modify event on event filter
  I need to be able to modify event on event filter

  @concurrent
  Scenario: given check event and drop event filter rule should update events count
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "drop",
      "description": "test-event-filter-che-event-filters-third-1-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-third-1"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
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
    [
      {
        "connector": "test-connector-che-event-filters-third-1",
        "connector_name": "test-connector-name-che-event-filters-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-1",
        "resource": "test-resource-che-event-filters-third-1-1",
        "state": 2
      },
      {
        "connector": "test-connector-che-event-filters-third-1",
        "connector_name": "test-connector-name-che-event-filters-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-1",
        "resource": "test-resource-che-event-filters-third-1-2",
        "state": 1
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules?search=che-event-filters-third-1 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "description": "test-event-filter-che-event-filters-third-1-description",
          "events_count": 2
        }
      ]
    }
    """

  @concurrent
  Scenario: given check event and break event filter rule should update events count
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "break",
      "description": "test-event-filter-che-event-filters-third-2-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-third-2"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
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
    [
      {
        "connector": "test-connector-che-event-filters-third-2",
        "connector_name": "test-connector-name-che-event-filters-third-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-2",
        "resource": "test-resource-che-event-filters-third-2-1",
        "state": 2
      },
      {
        "connector": "test-connector-che-event-filters-third-2",
        "connector_name": "test-connector-name-che-event-filters-third-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-2",
        "resource": "test-resource-che-event-filters-third-2-2",
        "state": 1
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules?search=che-event-filters-third-2 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "description": "test-event-filter-che-event-filters-third-2-description",
          "events_count": 2
        }
      ]
    }
    """

  @concurrent
  Scenario: given check event and enrichment event filter rule should update events count
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-third-3-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-third-3"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "customer",
            "description": "Customer",
            "value": "{{ `{{ .Event.ExtraInfos.customer }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-che-event-filters-third-3",
        "connector_name": "test-connector-name-che-event-filters-third-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-3",
        "resource": "test-resource-che-event-filters-third-3-1",
        "customer": "test-customer-che-event-filters-third-3-1",
        "state": 2
      },
      {
        "connector": "test-connector-che-event-filters-third-3",
        "connector_name": "test-connector-name-che-event-filters-third-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-3",
        "resource": "test-resource-che-event-filters-third-3-2",
        "customer": "test-customer-che-event-filters-third-3-2",
        "state": 1
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules?search=che-event-filters-third-3 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "description": "test-event-filter-che-event-filters-third-3-description",
          "events_count": 2
        }
      ]
    }
    """

  @concurrent
  Scenario: given check event and enrichment event filter rule should update zero events count after update
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-third-4-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-third-4"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "customer",
            "description": "Customer",
            "value": "{{ `{{ .Event.ExtraInfos.customer }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-che-event-filters-third-4",
        "connector_name": "test-connector-name-che-event-filters-third-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-4",
        "resource": "test-resource-che-event-filters-third-4-1",
        "customer": "test-customer-che-event-filters-third-4-1",
        "state": 2
      },
      {
        "connector": "test-connector-che-event-filters-third-4",
        "connector_name": "test-connector-name-che-event-filters-third-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-event-filters-third-4",
        "resource": "test-resource-che-event-filters-third-4-2",
        "customer": "test-customer-che-event-filters-third-4-2",
        "state": 1
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules?search=che-event-filters-third-4 until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "description": "test-event-filter-che-event-filters-third-4-description",
          "events_count": 2
        }
      ]
    }
    """
    When I do PUT /api/v4/eventfilter/rules/{{ .ruleId }}:
    """json
    {
      "type": "enrichment",
      "description": "test-event-filter-che-event-filters-third-4-description",
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-che-event-filters-third-4"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_field_from_template",
            "name": "customer",
            "description": "Customer",
            "value": "{{ `{{ .Event.ExtraInfos.customer }}` }}"
          },
          {
            "type": "set_field_from_template",
            "name": "domain",
            "description": "Domain",
            "value": "{{ `{{ .Event.ExtraInfos.domain }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/eventfilter/rules?search=che-event-filters-third-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "description": "test-event-filter-che-event-filters-third-4-description",
          "events_count": 0
        }
      ]
    }
    """
