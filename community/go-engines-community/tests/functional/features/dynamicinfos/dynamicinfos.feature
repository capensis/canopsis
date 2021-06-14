Feature: add infos to alarm
  I need to be able to add infos to alarm

  Scenario: given dynamic infos and check event should update alarm
    Given I am admin
    When I do POST /api/v2/dynamic-infos:
    """
    {
      "name": "test-dynamicinfos-dynamicinfos-1-name",
      "description": "test-dynamicinfos-dynamicinfos-1-description",
      "disable_during_periods": [],
      "entity_patterns":[{"name":"test-resource-dynamicinfos-1"}],
      "infos": [
        {"name":"test-info-dynamicinfos-1-name", "value":"test-info-dynamicinfos-1-value"}
      ]
    }
    """
    Then the response code should be 200
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector" : "test-connector-dynamicinfos-1",
      "connector_name" : "test-connector-name-dynamicinfos-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-dynamicinfos-1",
      "resource" : "test-resource-dynamicinfos-1",
      "state" : 1,
      "output" : "test-output-dynamicinfos-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.infos.*.test-info-dynamicinfos-1-name":"test-info-dynamicinfos-1-value"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "infos": {
            "{{ .ruleID }}": {
              "test-info-dynamicinfos-1-name": "test-info-dynamicinfos-1-value"
            }
          },
          "v": {
            "component": "test-component-dynamicinfos-1",
            "connector": "test-connector-dynamicinfos-1",
            "connector_name": "test-connector-name-dynamicinfos-1",
            "infos": {
              "{{ .ruleID }}": {
                "test-info-dynamicinfos-1-name": "test-info-dynamicinfos-1-value"
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