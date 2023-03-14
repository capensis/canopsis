Feature: Get a link rule
  I need to be able to get a link rule
  Only admin should be able to get a link rule

  @concurrent
  Scenario: given search request should return link rules
    When I am admin
    When I do GET /api/v4/link-rules?search=test-link-rule-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-link-rule-to-get-1",
          "name": "test-link-rule-to-get-1-name",
          "type": "alarm",
          "enabled": true,
          "links": [
            {
              "label": "test-link-rule-to-get-1-link-1-label",
              "category": "test-link-rule-to-get-1-link-1-category",
              "icon_name": "test-link-rule-to-get-1-link-1-icon",
              "url": "http://test-link-rule-to-get-1-link-1-url.com"
            }
          ],
          "external_data": null,
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-link-rule-to-get-1-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-link-rule-to-get-1-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-link-rule-to-get-2",
          "name": "test-link-rule-to-get-2-name",
          "type": "alarm",
          "enabled": true,
          "source_code": "function generate(alarms) { return [{\"label\": \"test-link-rule-to-get-2-link-1-label\",\"category\": \"test-link-rule-to-get-2-link-1-category\",\"icon_name\": \"test-link-rule-to-get-2-link-1-icon\",\"url\": \"http://test-link-rule-to-get-2-link-1-url.com\"}] }",
          "external_data": null,
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-link-rule-to-get-2-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-link-rule-to-get-2-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-link-rule-to-get-3",
          "name": "test-link-rule-to-get-3-name",
          "type": "entity",
          "enabled": true,
          "links": [
            {
              "label": "test-link-rule-to-get-3-link-1-label",
              "category": "test-link-rule-to-get-3-link-1-category",
              "icon_name": "test-link-rule-to-get-3-link-1-icon",
              "url": "http://test-link-rule-to-get-3-link-1-url.com"
            }
          ],
          "external_data": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-link-rule-to-get-3-pattern"
                }
              }
            ]
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
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
  Scenario: given get request should return link rule
    When I am admin
    When I do GET /api/v4/link-rules/test-link-rule-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-link-rule-to-get-1",
      "name": "test-link-rule-to-get-1-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-get-1-link-1-label",
          "category": "test-link-rule-to-get-1-link-1-category",
          "icon_name": "test-link-rule-to-get-1-link-1-icon",
          "url": "http://test-link-rule-to-get-1-link-1-url.com"
        }
      ],
      "external_data": null,
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-get-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-get-1-pattern"
            }
          }
        ]
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992
    }
    """

  @concurrent
  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/link-rules
    Then the response code should be 401

  @concurrent
  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/link-rules
    Then the response code should be 403

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/link-rules/test-link-rule-to-get-1
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/link-rules/test-link-rule-to-get-1
    Then the response code should be 403

  @concurrent
  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/link-rules/test-link-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
