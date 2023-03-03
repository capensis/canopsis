Feature: get subcategories list
  I need to be able to get subcategories
  Only admin should be able to get subcategories

  Scenario: GET unauthorized
    When I do GET /api/v4/cat/ocws-noc/subcategories?ci_provenance=project&category_value=category-to-get-1
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/cat/ocws-noc/subcategories?ci_provenance=project&category_value=category-to-get-1
    Then the response code should be 403

  Scenario: GET subcategories success
    When I am admin
    When I do GET /api/v4/cat/ocws-noc/subcategories?ci_provenance=project&category_value=category-to-get-1
    Then the response code should be 200
    Then the response body should be:
    ```
    [
        {
            "sys_id": "choice-to-get-1-subcategory-1",
            "label": "project - 1 - accounting",
            "value": "PROJECT - 1 - ACCOUNTING"
        },
        {
            "sys_id": "choice-to-get-1-subcategory-2",
            "label": "project - 2 - accounting",
            "value": "PROJECT - 2 - ACCOUNTING"
        }
    ]
    ```

  Scenario: GET subcategories when category not found
    When I am admin
    When I do GET /api/v4/cat/ocws-noc/subcategories?ci_provenance=project&category_value=category-to-get-not-found
    Then the response code should be 200
    Then the response body should be:
    ```
    []
    ```