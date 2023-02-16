Feature: get fill-form request
  I need to be able to get fill-form data
  Only admin should be able to get fill-form

  Scenario: GET unauthorized
    When I do GET /api/v4/cat/ocws-noc/fill-form?id=insignificant%2Fentity-to-fill-form-1
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/cat/ocws-noc/fill-form?id=insignificant%2Fentity-to-fill-form-1
    Then the response code should be 403

  Scenario: GET fill-form success
    When I am admin
    When I do GET /api/v4/cat/ocws-noc/fill-form?id=insignificant%2Fentity-to-fill-form-1
    Then the response code should be 200
    Then the response body should contain:
    ```
    {
        "contracts": [
            {
                "sys_id": "ocws_noc_snow_service_contract-to-fill-form-1-sys_id-1",
                "u_contract_owner": "u_contract_owner-1"
            },
            {
                "sys_id": "ocws_noc_snow_service_contract-to-fill-form-1-sys_id-2",
                "u_contract_owner": "u_contract_owner-2"
            },
            {
                "sys_id": "ocws_noc_snow_service_contract-to-fill-form-1-sys_id-3",
                "u_contract_owner": "u_contract_owner-3"
            },
            {
                "sys_id": "ocws_noc_snow_service_contract-to-fill-form-1-sys_id-4",
                "u_contract_owner": "u_contract_owner-4"
            },
            {
                "sys_id": "ocws_noc_snow_service_contract-to-fill-form-1-sys_id-5",
                "u_contract_owner": "u_contract_owner-5"
            }
        ],
        "categories": [
            {
                "sys_id": "choice-fill-form-1-category-1",
                "label": "category - 1 - fill-form-1",
                "value": "CATEGORY - 1 - fill-form-1"
            },
            {
                "sys_id": "choice-fill-form-1-category-2",
                "label": "category - 2 - fill-form-1",
                "value": "CATEGORY - 2 - fill-form-1"
            },
            {
                "sys_id": "choice-fill-form-1-category-3",
                "label": "category - 3 - fill-form-1",
                "value": "CATEGORY - 3 - fill-form-1"
            }
        ],
        "subCategories": [
            {
                "sys_id": "choice-fill-form-1-subcategory-1"
            },
            {
                "sys_id": "choice-fill-form-1-subcategory-2"
            },
            {
                "sys_id": "choice-fill-form-1-subcategory-3"
            },
            {
                "sys_id": "choice-fill-form-1-subcategory-4"
            }
        ],
        "branches": [
            {
                "sys_id": "location-to-fill-form-1",
                "group": "group-to-fill-form-1-name-1"
            },
            {
                "sys_id": "location-to-fill-form-1",
                "group": "group-to-fill-form-1-name-2"
            }
        ],
        "infos": {
            "ci_company_snow": {
                "name": "ci_company_snow",
                "value": "entity-to-fill-form-1-ci_company_snow-1"
            },
            "ci_location_snow": {
                "name": "ci_location_snow",
                "value": "entity-to-fill-form-1-ci_location_snow-1"
            },
            "ci_provenance": {
                "name": "ci_provenance",
                "value": "entity-to-fill-form-1-provenance-1"
            },
            "ci_type_snow": {
                "name": "ci_type_snow",
                "value": "entity-to-fill-form-1-ci_type_snow-1"
            },
            "company_is_communauty_snow": {
                "name": "company_is_communauty_snow",
                "value": "true"
            },
            "flag": {
                "name": "flag",
                "value": 1
            },
            "id_equipement": {
                "name": "id_equipement",
                "description": "entity-to-fill-form-1-description",
                "value": "entity-to-fill-form-1"
            }
        }
    }
    ```

  Scenario: GET fill-form not found
    When I am admin
    When I do GET /api/v4/cat/ocws-noc/fill-form?id=insignificant%2Fentity-to-fill-form-not-found
    Then the response code should be 404
    Then the response body should be:
    ```
    {
        "error": "Not found"
    }
    ```