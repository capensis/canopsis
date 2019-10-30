Environment setup
============================
If you need change environment for nightwatch you can do it in the `tests/e2e/.env.local`

For example: If you are receiving the problem with `chromedriver` version, you can set `CHROME_DRIVER_PATH`.
```
# Example for yarn global package
CHROME_DRIVER_PATH=/usr/local/share/.config/yarn/global/node_modules/chromedriver/lib/chromedriver/chromedriver
```


Folder structure conventions
============================
We have the following project structure:
```
e2e
├── custom-assertions        # Folder with assertions
│   └── elementsCount.js     # Special global assertion
├── custom-commands          # Folder with commands
│   ├── completed            # Folder with completed commands. Completed ui action. For example: login, logout and etc.
│   │   └── login.js         # Completed ui command
│   └── customClick.js       # Global command
├── helpers                  # Folder with helpers
│   └── sel.js               # Special helper for something
├── page-objects             # Folder with page objects (https://nightwatchjs.org/guide/#working-with-page-objects)
│   └── admin                # Folder for unioning of admin pages page-objects
│       └── users.js         # Page-object for /admin/users page
├── records                  # Folder with nightwatch-record unsuccessful tests records
├── reports                  # Folder with nightwatch xml reports
├── specs                    # Folder with tests
│   └── 01-auth              # Folder with several tests file. If we want to order files we can put number prefix.
│       └── auth.js          # File with tests. If we want to order files we can put number prefix here. If you want to exclude a test from parallel running, you need to add a suffix .consistently.
├── globals.js               # Global nightwatch methods, properties and etc. (http://nightwatchjs.org/guide#external-globals)
├── nightwatch.config.js     # File with settings for nightwatch (http://nightwatchjs.org/gettingstarted#settings-file)
```

Page Object API
===============
We are using page object api for every pages (https://nightwatchjs.org/guide/#working-with-page-objects).

in the page object we should implement the following moments:

* `url` address of the page. It can be function or property
* `commands` with every user actions: Click on something, setSomething into input, verify something before/after and etc.
* Every `elements` which we need: submitButton, usernameField, view button and etc. Here we have several moments:
    * If we can't identify UI element by existing attributes (class, id, special attribute and etc.) we can put `data-test` attribute for the element in the our application and use it. This attribute will remove from code on the production environment.
    * If we have dynamic element identify (for example: href attribute with view id or something else) we can put `el` helper and use it in the page object instance. Example in `page-objects/layout.js`
* We can split page to `sections` and put it if we want.

Custom command/assertion convention
===================================
We will consider `custom-commands` folder but this we must apply and on `custom-assertions`

* We should put command into global `custom-commands` folder if we need it more than on one test/page-object.
* If we can move this command into page-object and use it in several tests we must do it.
* If our global command it is the simple user action (like: `click`, `setValue` and etc.) we should move it just into `custom-commands` folder. But if our command completed user action (like: `login`, `logout`) we should put it into `custom-commands/completed` folder.

Tests writing convention
========================
If we need set the order of tests running we can put number prefixes for tests files. But it will be better if we will write tests isolated (So that the tests do not depend on each other).

Planned tests
=============
**We must keep this list up to date!**

1. Auth
    - [x] Correct user credentials login
    - [x] Authorized user logout
2. Admin
    * Users
        - [x] Create new user with some name
        - [x] Login by created user credentials
        - [x] Edit special user with username from constants
        - [x] Login by disabled user credentials
        - [x] Remove user with some username
        - [x] Create mass users with some name
        - [x] Check pagination users table
        - [x] Delete mass users with some name
    * Roles
        - [x] Create new role with some name
        - [x] Edit role with some name
        - [x] Remove role with some username
        - [x] Create mass roles with some name
        - [x] Search roles
        - [x] Check pagination roles table
        - [x] Delete mass roles with some name
        - [ ] Check role default view working
    * Rights
        - [ ] Create new right with some name
        - [ ] Adds right new role
    * Parameters
        - [x] Edit app title
        - [x] Switch language
        - [x] Edit footer text
        - [x] Edit description text
        - [ ] Upload logo
        - [x] Check global language
        - [x] Check app title
        - [x] Check login footer
        - [x] Check login description
        - [ ] Check logo
3. Layout
    * Top Bar
        - [x] Open current user modal
        - [x] Select current user default view
        - [x] Check default view
        - [ ] Switch user language
        - [ ] Check user interface language
    * Group Side Bar
        - [x] Add view with some name from constants
        - [x] Checking view copy with name from constants
        - [x] Editing test view with name from constants
        - [x] Deleting all test items view with name from constants
        - [x] Deleting all test group with name from constants
    * Group Top Bar
        - [x] Add view with some name from constants
        - [x] Checking view copy with name from constants
        - [x] Editing test view with name from constants
        - [x] Deleting all test items view with name from constants
        - [x] Deleting all test group with name from constants
4. View
    * Base functions
        - [ ] View open by `id`
        - [ ] View open by `id` and `tabId`
        - [ ] Add tab
        - [ ] Edit tab
        - [ ] Delete tab
        - [ ] Move tab by drag'n'drop
5. Widget
    * Base functions
        - [ ] Copy widget
        - [ ] Delete widget
        - [ ] Delete row
    * Alarm list
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
        * Header
            - [ ] Search
            - [ ] Pagination
            - [ ] Filters
            - [ ] Reporting
        * Hide menu
            - [ ] Periodical behavior
            - [ ] Ack
            - [ ] Fast ack
            - [ ] Cancel ack
            - [ ] Cancel alarm
        * Body
            - [ ] Ack
            - [ ] Fast ack
            - [ ] Snooze alarm
            - [ ] Declare ticket
            - [ ] Associate ticket
            - [ ] Cancel alarm
            - [ ] Periodical behavior
            - [ ] List periodic behaviors
            - [ ] List of available variables
        * Footer
            - [ ] Pagination
            - [ ] Pages
    * Context explorer
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
        * Header
            - [ ] Search
            - [ ] Pagination
            - [ ] Filters
            - [ ] Entities
            - [ ] Watcher
        * Hide menu
            - [ ] Delete entities
            - [ ] Periodical behavior
        * Body
            - [ ] Edit entities
            - [ ] Duplicate entities
            - [ ] Delete entities
            - [ ] Periodical behavior
        * Footer
            - [ ] Pagination
            - [ ] Pages
    * Service weather
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
        * Widget
            - [ ] Help
            - [ ] Info
    * Stats histogram
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
    * Stats table
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
    * Stats calendar
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
    * Stats curves
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
    * Stats number
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
        * Footer
            - [ ] Pages
    * Pareto diagram
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
    * Text
        * Base functions
            - [x] Create widget
            - [ ] Edit widget
6. Exploitation
    * Event filter
        - [ ] Create event filter rule
        - [ ] Edit event filter rule
        - [ ] Copy event filter rule
        - [ ] Delete event filter rule
        - [ ] Create mass event filter rule
        - [ ] Check pagination roles table
        - [ ] Delete mass event filter rule
    * PBehaviors
        - [ ] Create periodical behavior
        - [ ] Edit periodical behavior
        - [ ] Delete periodical behavior
        - [ ] Create mass periodical behavior
        - [ ] Search periodical behavior
        - [ ] Check pagination periodical behavior
        - [ ] Delete mass periodical behavior
    * Webhooks
        - [ ] Create webhook
        - [ ] Edit webhook
        - [ ] Delete webhook
        - [ ] Create mass webhook
        - [ ] Check pagination webhook
        - [ ] Delete mass webhook
