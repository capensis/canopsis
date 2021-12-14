Feature: add infos to alarm
  I need to be able to add infos to alarm

  Scenario: given dynamic infos and check event should update alarm
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-1-name",
      "description": "test-dynamicinfos-dynamicinfos-1-description",
      "disable_during_periods": [],
      "entity_patterns":[{"name":"test-resource-dynamicinfos-1"}],
      "infos": [
        {
          "name": "test-info-dynamicinfos-1-name",
          "value": "Alarm field: {{ `{{ .Alarm.Value.Connector }}` }}; Entity field: {{ `{{ .Entity.ID }}` }}; Event field: {{ `{{ .Event.ExtraInfos.customer }}` }}"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-1",
      "connector_name" : "test-connector-name-dynamicinfos-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-1",
      "resource" : "test-resource-dynamicinfos-1",
      "state" : 1,
      "output" : "test-output-dynamicinfos-1",
      "customer": "test-customer-dynamicinfos-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"v.resource":"test-resource-dynamicinfos-1"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "infos": {
            "{{ .ruleID }}": {
              "test-info-dynamicinfos-1-name": "Alarm field: test-connector-dynamicinfos-1; Entity field: test-resource-dynamicinfos-1/test-component-dynamicinfos-1; Event field: test-customer-dynamicinfos-1"
            }
          },
          "v": {
            "component": "test-component-dynamicinfos-1",
            "connector": "test-connector-dynamicinfos-1",
            "connector_name": "test-connector-name-dynamicinfos-1",
            "infos": {
              "{{ .ruleID }}": {
                "test-info-dynamicinfos-1-name": "Alarm field: test-connector-dynamicinfos-1; Entity field: test-resource-dynamicinfos-1/test-component-dynamicinfos-1; Event field: test-customer-dynamicinfos-1"
              }
            },
            "resource": "test-resource-dynamicinfos-1"
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

  Scenario: given updated dynamic infos should remove info from alarm with already added info
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-2-name",
      "description": "test-dynamicinfos-dynamicinfos-2-description",
      "disable_during_periods": [],
      "entity_patterns":[{"name":"test-resource-dynamicinfos-2"}],
      "infos": [
        {"name":"test-info-dynamicinfos-2-name", "value":"test-info-dynamicinfos-2-value"}
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-2",
      "connector_name" : "test-connector-name-dynamicinfos-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-2",
      "resource" : "test-resource-dynamicinfos-2",
      "state" : 1,
      "output" : "test-output-dynamicinfos-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.infos.*.test-info-dynamicinfos-2-name":"test-info-dynamicinfos-2-value"}]}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "infos": {
            "{{ .ruleID }}": {
              "test-info-dynamicinfos-2-name": "test-info-dynamicinfos-2-value"
            }
          },
          "v": {
            "infos": {
              "{{ .ruleID }}": {
                "test-info-dynamicinfos-2-name": "test-info-dynamicinfos-2-value"
              }
            },
            "resource": "test-resource-dynamicinfos-2"
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
    When I do PUT /api/v4/cat/dynamic-infos/{{ .ruleID }}:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-2-name",
      "description": "test-dynamicinfos-dynamicinfos-2-description",
      "disable_during_periods": [],
      "entity_patterns":[{"name":"test-resource-dynamicinfos-2-updated"}],
      "infos": [
        {"name":"test-info-dynamicinfos-2-name", "value":"test-info-dynamicinfos-2-value"}
      ],
      "enabled": true
    }
    """
    Then the response code should be 200
    When I wait the next periodical process
    When I do GET /api/v4/alarms?filter={"$and":[{"v.infos.*.test-info-dynamicinfos-2-name":"test-info-dynamicinfos-2-value"}]}
    Then the response code should be 200
    Then the response body should be:
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

  Scenario: given removed dynamic infos should remove info from alarm with already added info
    Given I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamicinfos-dynamicinfos-3-name",
      "description": "test-dynamicinfos-dynamicinfos-3-description",
      "disable_during_periods": [],
      "entity_patterns":[{"name":"test-resource-dynamicinfos-3"}],
      "infos": [
        {"name":"test-info-dynamicinfos-3-name", "value":"test-info-dynamicinfos-3-value"}
      ],
      "enabled": true
    }
    """
    Then the response code should be 201
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector" : "test-connector-dynamicinfos-3",
      "connector_name" : "test-connector-name-dynamicinfos-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-3",
      "resource" : "test-resource-dynamicinfos-3",
      "state" : 1,
      "output" : "test-output-dynamicinfos-3"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.infos.*.test-info-dynamicinfos-3-name":"test-info-dynamicinfos-3-value"}]}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "infos": {
            "{{ .ruleID }}": {
              "test-info-dynamicinfos-3-name": "test-info-dynamicinfos-3-value"
            }
          },
          "v": {
            "resource": "test-resource-dynamicinfos-3"
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
    When I do DELETE /api/v4/cat/dynamic-infos/{{ .ruleID }}
    Then the response code should be 204
    When I wait the next periodical process
    When I do GET /api/v4/alarms?filter={"$and":[{"v.infos.*.test-info-dynamicinfos-3-name":"test-info-dynamicinfos-3-value"}]}
    Then the response code should be 200
    Then the response body should be:
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
