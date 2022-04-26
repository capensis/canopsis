Feature: modify event on fifo event filter
  I need to be able to modify event on fifo event filter

  Scenario: given check event and change_entity rule should change entity defining fields by value
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test",
      "type": "change_entity",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test_connector"
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
        "component": "new component"
      },
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test_connector",
      "connector_name": "test_connector_name",
      "source_type": "resource",
      "event_type": "check",
      "component": "some useless data",
      "resource": "fifo_event_filters_test_resource",
      "customer_tags": "CMDB:TEST_PROD",
      "output": "blabla",
      "state": 2
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=fifo_event_filters_test_resource
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "fifo_event_filters_test_resource",
          "component": "new component"
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

  Scenario: given check event and change_entity rule should change entity defining fields by template
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test",
      "type": "change_entity",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test_connector"
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
    """
    {
      "connector": "test_connector",
      "connector_name": "test_connector_name",
      "source_type": "resource",
      "event_type": "check",
      "component": "some useless data",
      "resource": "fifo_event_filters_test_resource_2",
      "customer_tags": "CMDB:TEST_PROD",
      "output": "blabla",
      "state": 2
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=fifo_event_filters_test_resource_2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "fifo_event_filters_test_resource_2",
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
