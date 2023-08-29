Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given get multi sort request should return sorted alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-multi-sort-get-2"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-1"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-3"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-multi-sort-get&multi_sort[]=t,asc&multi_sort[]=v.last_event_date,asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-multi-sort-get-2"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-3"
        },
        {
          "_id": "test-alarm-to-multi-sort-get-1"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  @concurrent
  Scenario: given tags filter with has_not condition should filter alarms properly
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-test-tag-filter&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-test-tag-filter-1"
        },
        {
          "_id": "test-alarm-to-test-tag-filter-2"
        },
        {
          "_id": "test-alarm-to-test-tag-filter-3"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-test-tag-filter&filters[]=test-widgetfilter-to-alarm-test-tag-filter&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-test-tag-filter-1"
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
  Scenario: given get opened alarms request should return alarms with links
    When I am admin
    When I do GET /api/v4/alarms?with_links=true&search=test-resource-to-alarm-link-get&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-link-get-1",
          "links": {
            "test-category-to-alarm-link-get-1": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-1",
                "label": "test-link-rule-to-alarm-link-get-1-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-1-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-1-link-1-url.com?user=root&resources[]=test-resource-to-alarm-link-get-1/test-component-default|test-link-mongo-data-regexp-1-status&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-1",
                "label": "test-link-rule-to-alarm-link-get-1-link-3-label",
                "icon_name": "test-link-rule-to-alarm-link-get-1-link-3-icon",
                "url": "http://test-link-rule-to-alarm-link-get-1-link-3-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default&",
                "single": true,
                "hide_in_menu": true
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-1-url.com?user=root&info[]=test-resource-to-alarm-link-get-1-info-1-val&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-3",
                "label": "test-link-rule-to-alarm-link-get-3-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-3-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-3-link-1-url.com?user=root&resources[]=test-resource-to-alarm-link-get-1/test-component-default&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-1-url.com?user=root&info[]=test-resource-to-alarm-link-get-1-info-1-val|test-link-mongo-data-1-status&"
              }
            ],
            "test-category-to-alarm-link-get-2": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-1",
                "label": "test-link-rule-to-alarm-link-get-1-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-1-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-1-link-2-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default&"
              }
            ],
            "test-category-to-alarm-link-get-3": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-2-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val&"
              }
            ],
            "test-category-to-alarm-link-get-4": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-3",
                "label": "test-link-rule-to-alarm-link-get-3-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-3-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-3-link-2-url.com?resources[]=test-resource-to-alarm-link-get-1/test-component-default&"
              }
            ],
            "test-category-to-alarm-link-get-5": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-2-url.com?info[]=test-resource-to-alarm-link-get-1-info-1-val|test-link-mongo-data-1-status&"
              }
            ]
          }
        },
        {
          "_id": "test-alarm-to-link-get-2",
          "links": {
            "test-category-to-alarm-link-get-1": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-1-url.com?user=root&info[]=test-resource-to-alarm-link-get-2-info-1-val&"
              },
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-1-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-1-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-1-url.com?user=root&info[]=test-resource-to-alarm-link-get-2-info-1-val|test-link-mongo-data-2-status&"
              }
            ],
            "test-category-to-alarm-link-get-3": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-2",
                "label": "test-link-rule-to-alarm-link-get-2-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-2-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-2-link-2-url.com?info[]=test-resource-to-alarm-link-get-2-info-1-val&"
              }
            ],
            "test-category-to-alarm-link-get-5": [
              {
                "rule_id": "test-link-rule-to-alarm-link-get-4",
                "label": "test-link-rule-to-alarm-link-get-4-link-2-label",
                "icon_name": "test-link-rule-to-alarm-link-get-4-link-2-icon",
                "url": "http://test-link-rule-to-alarm-link-get-4-link-2-url.com?info[]=test-resource-to-alarm-link-get-2-info-1-val|test-link-mongo-data-2-status&"
              }
            ]
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
  Scenario: given get alarms request with search by entity infos should return alarms
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-to-alarm-get-1-info-1-value&active_columns[]=entity.infos.test-resource-to-alarm-get-1-info-1.value&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-to-get-1"
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
