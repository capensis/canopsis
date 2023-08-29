Feature: modify event on fifo event filter
  I need to be able to modify event on fifo event filter

  @concurrent
  Scenario: given check event and change_entity rule with invalid connector template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-1-description",
      "enabled": true,
      "config": {
        "connector": "{{ `{{ .Event.Connector` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-1"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-1",
        "connector": "test-connector-fifo-event-filters-failure-1",
        "connector_name": "test-connector-name-fifo-event-filters-failure-1",
        "component": "test-component-fifo-event-filters-failure-1",
        "resource": "test-resource-fifo-event-filters-failure-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-1",
        "connector": "test-connector-fifo-event-filters-failure-1",
        "connector_name": "test-connector-name-fifo-event-filters-failure-1",
        "component": "test-component-fifo-event-filters-failure-1",
        "resource": "test-resource-fifo-event-filters-failure-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "invalid template \"Connector\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
        },
        {
          "type": 1,
          "message": "invalid template \"Connector\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
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
  Scenario: given check event and change_entity rule with not found connector template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-2-description",
      "enabled": true,
      "config": {
        "connector": "{{ `{{ .Event.ExtraInfos.notfound }}` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-2"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-2",
        "connector": "test-connector-fifo-event-filters-failure-2",
        "connector_name": "test-connector-name-fifo-event-filters-failure-2",
        "component": "test-component-fifo-event-filters-failure-2",
        "resource": "test-resource-fifo-event-filters-failure-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-2",
        "connector": "test-connector-fifo-event-filters-failure-2",
        "connector_name": "test-connector-name-fifo-event-filters-failure-2",
        "component": "test-component-fifo-event-filters-failure-2",
        "resource": "test-resource-fifo-event-filters-failure-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "cannot execute template \"Connector\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-2",
            "connector_name": "test-connector-name-fifo-event-filters-failure-2",
            "component": "test-component-fifo-event-filters-failure-2",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"Connector\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-2",
            "connector_name": "test-connector-name-fifo-event-filters-failure-2",
            "component": "test-component-fifo-event-filters-failure-2",
            "source_type": "resource"
          },
          "unread": true
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
  Scenario: given check event and change_entity rule with invalid connector_name template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-3-description",
      "enabled": true,
      "config": {
        "connector_name": "{{ `{{ .Event.ConnectorName` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-3"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-3",
        "connector": "test-connector-fifo-event-filters-failure-3",
        "connector_name": "test-connector-name-fifo-event-filters-failure-3",
        "component": "test-component-fifo-event-filters-failure-3",
        "resource": "test-resource-fifo-event-filters-failure-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-3",
        "connector": "test-connector-fifo-event-filters-failure-3",
        "connector_name": "test-connector-name-fifo-event-filters-failure-3",
        "component": "test-component-fifo-event-filters-failure-3",
        "resource": "test-resource-fifo-event-filters-failure-3-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "invalid template \"ConnectorName\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
        },
        {
          "type": 1,
          "message": "invalid template \"ConnectorName\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
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
  Scenario: given check event and change_entity rule with not found connector_name template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-4-description",
      "enabled": true,
      "config": {
        "connector_name": "{{ `{{ .Event.ExtraInfos.notfound }}` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-4"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-4",
        "connector": "test-connector-fifo-event-filters-failure-4",
        "connector_name": "test-connector-name-fifo-event-filters-failure-4",
        "component": "test-component-fifo-event-filters-failure-4",
        "resource": "test-resource-fifo-event-filters-failure-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-4",
        "connector": "test-connector-fifo-event-filters-failure-4",
        "connector_name": "test-connector-name-fifo-event-filters-failure-4",
        "component": "test-component-fifo-event-filters-failure-4",
        "resource": "test-resource-fifo-event-filters-failure-4-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "cannot execute template \"ConnectorName\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-4",
            "connector_name": "test-connector-name-fifo-event-filters-failure-4",
            "component": "test-component-fifo-event-filters-failure-4",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"ConnectorName\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-4",
            "connector_name": "test-connector-name-fifo-event-filters-failure-4",
            "component": "test-component-fifo-event-filters-failure-4",
            "source_type": "resource"
          },
          "unread": true
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
  Scenario: given check event and change_entity rule with invalid component template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-5-description",
      "enabled": true,
      "config": {
        "component": "{{ `{{ .Event.Component` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-5"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-5",
        "connector": "test-connector-fifo-event-filters-failure-5",
        "connector_name": "test-connector-name-fifo-event-filters-failure-5",
        "component": "test-component-fifo-event-filters-failure-5",
        "resource": "test-resource-fifo-event-filters-failure-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-5",
        "connector": "test-connector-fifo-event-filters-failure-5",
        "connector_name": "test-connector-name-fifo-event-filters-failure-5",
        "component": "test-component-fifo-event-filters-failure-5",
        "resource": "test-resource-fifo-event-filters-failure-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "invalid template \"Component\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
        },
        {
          "type": 1,
          "message": "invalid template \"Component\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
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
  Scenario: given check event and change_entity rule with not found component template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-6-description",
      "enabled": true,
      "config": {
        "component": "{{ `{{ .Event.ExtraInfos.notfound }}` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-6"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-6",
        "connector": "test-connector-fifo-event-filters-failure-6",
        "connector_name": "test-connector-name-fifo-event-filters-failure-6",
        "component": "test-component-fifo-event-filters-failure-6",
        "resource": "test-resource-fifo-event-filters-failure-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-6",
        "connector": "test-connector-fifo-event-filters-failure-6",
        "connector_name": "test-connector-name-fifo-event-filters-failure-6",
        "component": "test-component-fifo-event-filters-failure-6",
        "resource": "test-resource-fifo-event-filters-failure-6-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "cannot execute template \"Component\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-6",
            "connector_name": "test-connector-name-fifo-event-filters-failure-6",
            "component": "test-component-fifo-event-filters-failure-6",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"Component\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-6",
            "connector_name": "test-connector-name-fifo-event-filters-failure-6",
            "component": "test-component-fifo-event-filters-failure-6",
            "source_type": "resource"
          },
          "unread": true
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
  Scenario: given check event and change_entity rule with invalid resource template should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-7-description",
      "enabled": true,
      "config": {
        "resource": "{{ `{{ .Event.Resource` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-7"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-7",
        "connector": "test-connector-fifo-event-filters-failure-7",
        "connector_name": "test-connector-name-fifo-event-filters-failure-7",
        "component": "test-component-fifo-event-filters-failure-7",
        "resource": "test-resource-fifo-event-filters-failure-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-7",
        "connector": "test-connector-fifo-event-filters-failure-7",
        "connector_name": "test-connector-name-fifo-event-filters-failure-7",
        "component": "test-component-fifo-event-filters-failure-7",
        "resource": "test-resource-fifo-event-filters-failure-7-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "invalid template \"Resource\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
        },
        {
          "type": 1,
          "message": "invalid template \"Resource\": template: tpl:1: unclosed action",
          "event": null,
          "unread": true
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
  Scenario: given check event and change_entity rule with not found resource template var should create failure
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "change_entity",
      "description": "test-event-filter-fifo-event-filters-failure-8-description",
      "enabled": true,
      "config": {
        "resource": "{{ `{{ .Event.ExtraInfos.notfound }}` }}"
      },
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-fifo-event-filters-failure-8"
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
    Then I save response ruleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-8",
        "connector": "test-connector-fifo-event-filters-failure-8",
        "connector_name": "test-connector-name-fifo-event-filters-failure-8",
        "component": "test-component-fifo-event-filters-failure-8",
        "resource": "test-resource-fifo-event-filters-failure-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-fifo-event-filters-failure-8",
        "connector": "test-connector-fifo-event-filters-failure-8",
        "connector_name": "test-connector-name-fifo-event-filters-failure-8",
        "component": "test-component-fifo-event-filters-failure-8",
        "resource": "test-resource-fifo-event-filters-failure-8-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/eventfilter/{{ .ruleId }}/failures until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "type": 1,
          "message": "cannot execute template \"Resource\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-8",
            "connector_name": "test-connector-name-fifo-event-filters-failure-8",
            "component": "test-component-fifo-event-filters-failure-8",
            "source_type": "resource"
          },
          "unread": true
        },
        {
          "type": 1,
          "message": "cannot execute template \"Resource\" for event: template: tpl:1:9: executing \"tpl\" at <.Event.ExtraInfos.notfound>: map has no entry for key \"notfound\"",
          "event": {
            "connector": "test-connector-fifo-event-filters-failure-8",
            "connector_name": "test-connector-name-fifo-event-filters-failure-8",
            "component": "test-component-fifo-event-filters-failure-8",
            "source_type": "resource"
          },
          "unread": true
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
