# Canopsis - Changelog

This document references all changes made to Canopsis since 2017/08/21. Some older lines may appear in their original language.

⚠️ **Note:** this ChangeLog is currently unmaintained since Canopsis 3.35.0, as some workflow/Git/build processes are being changed. In the meantime, please consult Git history and <https://doc.canopsis.net>.

## Canopsis 3.34.0 - Due date : 2019-12-18

- [Documentation] Add Service Weather API documentation
- [Documentation] Add documentation for the new engine-dynamic-infos CAT engine
- [Go] CAT: Add a new engine-dynamic-infos engine
- [Go] CAT: Add a new engine-webhook engine; drop the old webhook.so plugin
- [UI] Alarms list ack and declare ticket action
- [UI] Add a 'Dynamic info' CRUD
- [UI] Create entity form (impact and depends)
- [UI] bad filter on pbehavior form
- [UI - Alarmlist] Assocticket mass action
- [UI - Context Explorer] 'Enabled' columns
- [UI - Filter editor] Add field suggestions

## Canopsis 3.33.0 - Due date : 2019-11-22

- [Documentation] Mention that resolve and unsnooze are incompatible with event\_patterns
- [API] Add a new /api/v2/dynamic-infos API, not completely usuable yet
- [Go] Handle RabbitMQ connection loss
- [Go] Fix cache usage in alarm service
- [Go] Fix corrupted traces with large events
- [Go] Remove dead code
- [Go] Allow to use null value for pattern lists
- [UI] Add helpers to deal with integers and string comparisons
- [UI - AlarmList] Fix ack\_resources usage when acknowleding an alarm
- [UI - Context] Fix impact/depends search in entity creation modal
- [UI - pbehavior] Refactor pbehavior form
- [UI - pbehavior] Fix "required" rule on pbehavior filter
- [UX] Various UX improvements in widget containers

## Canopsis 3.32.0 - Due date : 2019-11-08

## Canopsis 3.31.0 - Due date : 2019-10-28

- [Packages] Introduce new "canoctl deploy-go" command
- [Packages] Add initialisation.toml.example, amqp2engines-python.conf.example and amqp2engines-go.conf.example reference files
- [Packages] Introduce new go-engines-vars.conf file for Go engine variables
- [Packages] Sync default Python/Go engines with current recommendations
- [Packages] Reduce package dependencies on CentOS 7
- [Packages] Upgrade from InfluxDB 1.5.4 to 1.5.5 on new installations
- [Packages] Upgrade internal canoctl Ansible version from 2.4.4 to 2.8.5
- [Documentation] Add a new documentation for zabbix2canopsis
- [Documentation] Add more filter examples in Usage Guide
- [Documentation] Document leavemail/trim options in email2canopsis
- [Go] Remove long\_output\_history and initial\_long\_output fields from events
- [Go] Only update alarm’s output when a check event is received
- [Go] Use zerolog in Go engines
- [Go] Manage dependencies with go mod instead of glide
- [Go] Refactored most of watcher engine to have better scalability
- [Go] Remove engine-stat from new installs; use statsng from CAT instead
- [UI - Translations] Some typo fixes and improvements in the French translation

## Canopsis 3.30.0 - Due date : 2019-10-11

- [Docker] Fix email2canopsis Docker image
- [Documentation] Improve Canopsis installation and upgrade guides
- [API] Fix possible error when a pbehavior has no rrule
- [API] Fix last\_ko/last\_event date format, when using engine-stat
- [Go] Significant performance improvements in JSON encoding/decoding
- [LDAP] Add a `username_attr` attribute
- [UI] Fix CAT widgets, so that they're not displayed in the Open-Core edition
- [UI - AlarmList] Fix snooze action regression in Canopsis 3.29.0
- [UI - Context] Fix Context Explorer info tabs
- [UI - Action] Fix comments and exdate fields in pbehavior form
- [UI - SNMP] Display SNMP view both in Python and Go stacks

## Canopsis 3.29.0 - Due date : 2019-09-27

 - [Documentation] Add alarm-filter documentation
 - [Documentation] Fix some 404 hyperlinks
 - [Documentation] Fix some resource/ressource spelling mistakes in API examples
 - [Docker] "init" and Open-core Go engine images are now public
 - [API] Remove TicketAPI API
 - [API] Prevent PUT /api/v2/pbehavior/<id> from removing the pbehaviors' comments and timezone
 - [Go] Release Go engines as open-source. Go plugins remain in CAT
 - [Go] Remove TicketAPI service
 - [Python] Add `postpone_if_active_pbehavior` option to alarm filter
 - [UI] Add a linter for pug templates
 - [UI] Filter available CRUDs by canopsis stack
 - [UI] Add possibility to clone a tab to another view
 - [UI] Add "module-resolver" Babel's plugin
 - [UI - Event filter] Add possibility to edit an action
 - [UI - Heartbeat] Add default rights for Heartbeat's CRUD
 - [UI - Heartbeat] Fix - Wrong labels
 - [UI - Heartbeat] Fix - Delete action
 - [UI - Heartbeat] Fix - Rights application
 - [UI - List alarm] Fix - Column selector on info popup setting
 - [UI - Tests] Tests e2e - Basic stats histograms widget function
 - [UI - Tests] Tests e2e - Basic stats table widget function
 - [UI - Tests] Tests e2e - Basic stats calendar widget function
 - [UI - Tests] Tests e2e - Basic stats curves widget function
 - [UI - Tests] Tests e2e - Basic text widget function
 - [UI - Tests] Tests e2e - Basic alarmlist widget function
 - [UI - Tests] Tests e2e - Remove user's default view
 - [UI - Tests] Tests e2e - Fix topbar tests
 - [UI - Tests] Tests e2e - Basic context explorer widget function

## Canopsis 3.28.0 - Due date : 2019-09-12

 - [Documentation] Big improvements in most guides of doc.canopsis.net
 - [Documentation] Add CAS authentication documentation
 - [Documentation] Fix API requests in Action engine documentation
 - [Packaging] Delete Debian 8 compatibility
 - [Packaging] Add an explicit dependency on an equal version of canopsis-engines-go, when installing canopsis-engines-go-cat
 - [Packaging] Upgrade from Go 1.12.7 to 1.12.9, in the provided binaries
 - [Packaging] Work around a compatibility problem between CentOS 7's old Python and pybars3
 - [UI] Fix the roles list pagination when the create modal is opened
 - [UI] Add an "Edit Action" button to the event-filter CRUD
 - [UI] Add a new Heartbeat CRUD
 - [UI - Listalarm] Fix issue when two identical filters are created
 - [UI - Listalarm] Fix "doesn't contain" filter
 - [UI - Stats] Exclude "No data" when sorting values in a stats table
 - [UI - Tests] Add tests for Weather widget
 - [UI - Tests] Add tests for the language settings of the Top bar

## Canopsis 3.27.0 - Due date : 2019-09-02

 - [Documentation] Mention that -enrichContext flag is required when doing enrichment in Go engines
 - [Packaging] Upgrading to a newer Canopsis release doesn't reset the authkey anymore
 - [API] Allow to set the pbehavior’s id on creation
 - [API] Add search option in pbehaviors API
 - [UI] Fix – Get app infos after login
 - [UI] Fix – Text editor going on top of top-bar on parameters view
 - [UI – Stats] Data format harmonization (durations, percentage) on stats widgets
 - [UI – Listalarm] Add a default period filter setting
 - [UI – Tests] e2e – Roles management view
 - [UI – Tests] e2e – Parameters view

## Canopsis 3.26.0 - Due date : 2019-08-19

 - [UI – Listalarm] Improve "info popup" UX
 - [UI – CRUD Pbehaviors] Add pagination parameters
 - [UI – Listalarm] Improve ack with ticket action workflow
 - [UI – Service weather] Add a message when there is no data to display and when there was an error fetching data
 - [UI – CRUD Webhooks] Add "empty_response" field
 - [UI] Clarified "Default filter" field label
 - [UI] Add translations on pattern edition fields
 - [UI – Tests] e2e tests – Add basic functions for views management
 - [UI – CRUD Rights] Fix – User creation
 - [UI – CRUD Rights] Fix – Page refresh after new right creation
 - [UI – Tests] e2e tests – Empty group deletion
 - [API] Add pagination in pbehavior route
 - [Webhook] Add empty_response field

## Canopsis 3.25.0 - Due date : 2019-08-02

 - [Documentation] Add documentation for importctx API
 - [Documentation] Document CpsTime fields, such as .Unix
 - [Build] Build Docker Go engines and Debian/CentOS Go packages with the same Go version
 - [Build] Fix unneccessary snmp-mibs-downloader dependency in Core, it was meant for CAT
 - [Build] Add canopsis-engines-go-cat packages for Debian and CentOS
 - [Build] Improve Go engines Docker Makefiles portability for macOS
 - [API - importctx] Fix "create" action and add a "set" action
 - [API - pbehavior] Fix a TypeError when dealing with a future date
 - [Go] Bump Go engines to Go 1.12.7, fixes dlv attachment on che and axe
 - [Go] Port long_output feature from Python to Go
 - [Service Weather] Refactor Service Weather route to significantly improve performance
 - [UI] Add search fields on Users and Roles administration views
 - [UI] Add an internal configuration for specific modals width
 - [UI] Enable empty views groups deletion
 - [UI] Fix – Roles mass deletion
 - [UI] Fix – Tooltips position on views menu
 - [UI – Listalarm] Fix Ack with ticket workflow
 - [UI – Tests] Add e2e tests for horizontal views menu
 - [UI – Tests] Add e2e tests commands for views management
 - [UI - Translations] Make the global default language configurable
 - [UI – Webhooks] Fix – Webhooks deletion

## Canopsis 3.24.0 - Due date : 2019-07-18

 - [Build] Build Docker images with "yarn" instead of "npm"
 - [Build] Fix missing initialisation.toml file in Go engine packages
 - [Build] Don't strip Go binaries in CentOS packages, to keep debug information
 - [Build] Fix RabbitMQ download link when running "canoctl deploy"
 - [API] Fix search on entities informations in get-alarms
 - [Go] Add formattedDate and replace functions to webhooks templates
 - [Go] Fix missing last_update_date computation
 - [UI] Add a label to "Actions" columns, on all tables
 - [UI] Add "description", "logo", and "login page footer" parameters on parameters view + Add a right to access this view
 - [UI] Harmonize default view selector modal
 - [UI] Fix – Close all modals on page change
 - [UI – Listalarm] Add dynamic date selector for live reporting functionnality
 - [UI – Listalarm] Add a setting to send (or not) a message with fast-ack events
 - [UI – Listalarm] Fix – Available actions on alarms that are not ack
 - [UI – Listalarm] Fix "More infos" action’s icon
 - [UI – Service Weather] Fix issue with templates not being displayed in the tiles
 - [UI – Stats] Add helpers on stats settings
 - [UI - Stats] Fix – Trend icon
 - [UI – Tests] Add fundations for e2e tests
 - [UI – Tests] Add e2e tests – Topbar
 - [UI – Tests] Add e2e tests – Users administration’s panel
 - [UI – Tests] Add e2e tests – Left sidebar
 - [UI – Weather] Change "More infos" modal header’s color according to watcher’s state

## Canopsis 3.23.0 - Due date : 2019-07-05

 - [Documentation] Add entity filters examples in Pbehavior use case documentation
 - [Documentation] Add documentation on database cleanup queries
 - [Documentation] Document the new watcher engine
 - [Build] Fix canoctl compatibility problems on some Debian configurations
 - [Build] Fix a syntax compatibility problem with Docker-Compose >= 1.14.0 versions
 - [Build] Make it possible to build most Canopsis Docker images from a macOS host system
 - [Build] Freeze amqp version to 2.4.2 (CentOS 7 compatibility)
 - [Build] Remove -s -w ldflags to keep debugging symbols in Go engines
 - [Build] Fix missing package dependencies for SNMP on Debian and CentOS
 - [Go] Add CPS_DEBUG_TRACE option to axe and watcher engines
 - [Go] Move watcher functionality from axe to a dedicated new watcher engine
 - [Go] Add runtime warnings about TicketService being deprecated and now replaced by Webhooks
 - [Go] Added missing declareticket trigger to action and webhooks
 - [Python] Remove tracebacks on AMQP connection errors due to loss of heartbeat
 - [UI] Automatically redirect to Login page when "401 - Forbidden" error is returned by a Canopsis API
 - [UI] Add "Delete" mass action on Webhooks and Event filter CRUD
 - [UI] Pattern editors now handle null and empty string values
 - [UI - Stats] When selecting month period, stats widget now take current month into account
 - [UI - Stats] Improve error hangling on stats widgets

## Canopsis 3.22.0 - Due date : 2019-06-28

 - [Documentation] Add documentation about the new `{{ internal-link }}` handlebar helper
 - [Documentation] Fix some `{{ entities name="" }}` cases in Service weather widget configuration
 - [Authentication] CAS authentication fixed to be compatible with current UIv3 instead of old UIv2
 - [UI - ListAlarm] Responsiveness fixes
 - [UI - ListAlarm] Timeline refactoring and bugfixes
 - [UI - Weather] Add a new `{{ internal-link }}` handlebar helper

## Canopsis 3.21.0 - Due date : 2019-06-21

 - [email2canopsis] Add leavemails option
 - [Python] Fix ack_time_sla and resolve_time_sla stats to make them take into account the alarms that have not been acknowledged or resolved
 - [Go] Fix profiling options
 - [Go] Optimize GetLastAlarms to fix a performance issue in axe
 - [Go] Ignore changestate events with info state
 - [API] Add recursive option for the MTBF statistic
 - [API] Count ok and ko statistics for the current service period (since the end of the last pbehavior, or since midnight)
 - [API] Fix advanced search in get-alarms
 - [UI] Harmonize settings order for stats widgets
 - [UI - Stats] Fix annotation line's label position
 - [UI - Pareto] Fix - Stat select's modal's title
 - [UI - Stats curves] Fix - Line tension
 - [UI - Stats table] Reorder columns
 - [UI - Stats curves] Point style setting
 - [UI] Fix - Error in console when multiple modals open
 - [UI - Stats] Fix - Confusing end date displayed on histogram and curves
 - [UI - Stats curves] Format durations on tooltip
 - [UI] Clone widget
 - [UI] Improve validations feedback (settings panel + forms on multiple tabs)
 - [UI] Refactor pattern simple editor
 - [UI] Fix list filter rights
 - [Documentation] Add informations about KPIs and their use
 - [Documentation] Add a process describing Canopsis database backup and restoration
 - [Build] Fix Go engines files in Debian and CentOS packages

## Canopsis 3.20.0 - Due date : 2019-06-07

 - [Python] Fix a race condition in the alarm filter
 - [UI] Durations are now formatted with number of days
 - [UI - List alarm] Fix - Advanced search column value computation
 - [UI - Stats table] Add possibility to display the trend next to the stats values
 - [UI - Stats table] Fix - Always display unit,  even with values equal to zero
 - [UI - Text widget] Clarify the fact that selecting a stat is optionnal on this widget

## Canopsis 3.19.0 - Due date : 2019-06-04

 - [Docker] Docker-Compose files now support CPS_EDITION and CPS_STACK variables
 - [API] Fix resolved alarm duration in get-alarms
 - [Go] Add role in alarm steps
 - [CAT] Fix error in statsng with entities containing a newline character
 - [CAT] Fix "Exceeded memory limit" error in statsng
 - [CAT] Fix missing events in statsng
 - [CAT] The MTBF statistic now takes the dependencies' alarms into account
 - [CAT] The "recursive" statistics now only take into account the entity's dependencies
 - [CAT] Fix invalid statistics when the beat has not been executed
 - [CAT] Fix "input field [...] is type float, already exists as type integer" errors in statsng
 - [UI] Fix - Filter validation with nested groups
 - [UI] Fix - Copy to clipboard functionnality
 - [UI] Fix - Fix sidebars/topbar z-indexes
 - [UI - Stats] Change default value for "Duration" setting
 - [UI - Stats] Fix - Automatically select last full hour when changing to custom date interval on "Duration" setting
 - [UI - Stats] Add new "Pareto diagram" widget
 - [UI - Stats] Fix - Loading overlay container
 - [UI - Stats curves] Add format on tooltips for rate stats
 - [UI - CRUD Pbehaviors] Fix - Pbehavior's invalid date formatting when editting a pbehavior containing timestamps with milliseconds
 - [UI - CRUD Pbehaviors] Fix - Complex filter edition
 - [UI - CRUD SNMP] Fix - Rule's "oid" computation
 - [UI - CRUD Event-filter] Fix - Rule's id autocompletion when duplicating a rule
 - [UI - CRUD Webhooks] Fix - Patterns edition on fields containing a dot

## Canopsis 3.18.1 - Due date : 2019-05-22

 - [Packaging] Various build automation fixes
 - [Packaging] Fix missing SNMPRULE rights for admin user
 - [Packaging] Fix missing default rights when calling canopsinit after a fresh install
 - [Documentation] Delete an older version of linkbuilder documentation
 - [UI - Rights] Fix JavaScript error while trying to apply linklist rights
 - [UI - Webhooks] Fix validation of URLs using templates

## Canopsis 3.18.0 - Due date : 2019-05-17

 - [Documentation] Improve the readability of the centreon connector reboot procedure
 - [Documentation] Improve linkbuilder documentation
 - [Documentation] Add a warning about the lack of support for internet explorer
 - [Documentation] Add a section about the use of `amqp2tty` on a docker environment and package environment
 - [Documentation] Fix search on some patterns
 - [Python] Fix issue in statsng when an entity is deleted
 - [Go] Add new values in watchers' output_template
 - [Go] Add alarms_acknowledged statistic
 - [Go] Fix alarms_canceled and ack_time_sla statistics
 - [Go] Update links with watchers when an entity is modified by the event-filter
 - [Go] Update alarm's output at each event
 - [Go] Return an error when an unknown value is used in a template in the event-filter
 - [Tooling] Fix default login page description
 - [UI] Fix - Empty Rrule on pbehavior creation
 - [UI] Refactor Login page style
 - [UI] Adapt default fields for filter editor on diferent forms
 - [UI] Add rights to manage access to linklist (Listalarm and Weather widgets)
 - [UI] Fix - Mix filters
 - [UI - ListAlarm] Fix - Always keep at least 3 visible actions buttons + Dropdown menu
 - [UI - ListAlarm] Add setting for each column to enable HTML interpretation + Setting for HTML interpretation on timeline
 - [UI - ListAlarm] Add format for `duration` and `current_state_duration` columns
 - [UI - Context] Adapt watcher's creation form to canopsis stack (python/go)
 - [UI - Weather] Fix - Watcher with minor state's color
 - [UI - Weather] Fix - Alarm list modal filter computation
 - [UI - Stats] Fix - Bug on duration select
 - [UI - Stats] Add new `current_ongoing_alarms_with_ack` and `current_ongoing_alarms_without_ack` stats on stats widgets
 - [UI - Stats] Add new `alarms_acknowledged` stats on stats widgets
 - [UI - Stats] Fix - Edit/Delete widget while data are loading
 - [UI - Stats table] Add format for `current_state` column
 - [UI - Stats table] Add format for `durations` and `percentage` stats value
 - [UI - Stats table] Add setting for table's default sorting column
 - [UI - Calendar] Fix - Display when there's more than one filter
 - [UI - Webhooks] Add 'Disable if active pbehavior' option to webhook creation/edition form
 - [UI - EventFilter] Fix - Event filter rules are enable by default (even when there's no `enabled` field)

## Canopsis 3.17.0 - Due date : 2019-05-07

 - [API] Prevent from sending a changestate or keepstate event with an info state
 - [API] Add informations to the app_info route
 - [Documentation] Add details about the event-filter versions in ackcentreon's documentation
 - [Documentation] Simplify canopsinit command in update guide
 - [Documentation] canopsinit: document the new --canopsis-edition and --canopsis-stack options
 - [Go] Compute watcher's state in axe's beat
 - [Go] Update dependencies to fix CI
 - [Go] Add split and trim functions to webhooks' templates
 - [Go] Add support for python pbehaviors
 - [Go] Add support for python pbehaviors in webhooks
 - [Go] Add -alwaysFlushEntities option to che
 - [Python] Fix compatibility issue between the statsng and axe engines
 - [Python] Add statistics for ongoing alarms with or without an ACK
 - [Python] Fix current_ongoing_alarms statistic which returned negative values
 - [Python] Add author's role in timeline
 - [Python] Add CPS_RECOMPUTE_PBEHAVIORS_ON_NEW_ENTITY option to pbehaviors
 - [UI] Fix - Pbehavior form validation
 - [UI] Add missing french translations
 - [UI] Refactor user profile + Fix user UI language storage
 - [UI] Refactor Login page style + Add 'login_description' display
 - [UI - List alarm] Display users roles on timeline
 - [UI - List alarm] Add a setting to enable multiple ack
 - [UI - List alarm] Delete 'Info' option on change state action's modal
 - [UI - List alarm] Fix - Pbehaviors comments display on extra details column
 - [UI - Context] Fix - Entity mass deletion
 - [UI - Weather] Add 'Alarm canceled' icon
 - [UI - Weather] Add possibility to show both "Alarm list" and "More infos" modal
 - [UI - Weather] Fix - Alarm list modal filter (compatibility with both watchers and watchersNg)
 - [UI - Webhooks] Fix - Display and edit webhooks with no 'declare_ticket' field
 - [UI - Webhooks] Add 'changestate' trigger
 - [UI - Webhooks] Enable 'event_pattern' with 'cancel' trigger
 - [UI - Webhooks] Add field to specify if webhooks need to be disabled when there's an active pbehavior

## Canopsis 3.16.0 - Due date : 2019-04-19

 - [Documentation] Add documentation for LDAP authentication
 - [Documentation] Add required kernel version for docker installation
 - [Documentation] Improve webhooks documentation
 - [Go] Add json and json_unquote functions to webhooks templates
 - [Go] Fix last_update_date, which was not updated by axe
 - [Go] Fix the computation of watchers' dependencies
 - [Go] Fix the conversion of large numbers to strings in webhooks
 - [Go] Fix CI by preventing the tests from being run in parallel
 - [Go] Add statechange trigger
 - [Python] Fix the webhook API, which returned invalid JSON
 - [Python] Fix the action API, which rejected valid actions
 - [Python] Add route that returns all the actions
 - [Python] Add login_page_description to the login_info route
 - [Python] Add ldap_uri in LDAP configuration

# Canopsis 3.15.0 - Due date : 2019-04-04

 - [Packaging] Deprecate Debian 8 support
 - [Documentation] Added pbehavior helper documentation
 - [Go] Allowed to have fast ack and ack on the same alarm
 - [Python] Added route with login, app and user interface informations
 - [Python] Fixed missing field error in service weather when using go watchers
 - [Python] Changed login info route to allow requests while not logged in
 - [UI] Fix - Display filter with long names
 - [UI] Add possibility to display custom logo/app title
 - [UI] Add variables values into helpers modals
 - [UI] Add users personal filters in Listalarm and Context widgets
 - [UI] Add possibility to insert images in text editor
 - [UI] Fix - Use user's id instead of crecord\_name
 - [UI - Filter editor] Fix - Automatic filter update when modifying it on advanced editor
 - [UI - Rights] Add rights for exploitation views
 - [UI - Service weather] Add cancel action
 - [UI - List alarm] Add a setting to require note field when sending an ack event
 - [UI - List alarm] Fix - Links not clickable
 - [UI - CRUD Pbehaviors] Add pbehaviors edition

## Canopsis 3.14.0 - Due date : 2019-03-28

 - [Python] Fix service weather 500 error when using Go watchers
 - [UI] Display user auth key in top bar
 - [UI - ListAlarm] Add type, author and comments on pbehavior extra details
 - [UI - ListAlarm] Fix immediate pbehavior displaying + Comment display condition
 - [UI - Rights] Fix rights behavior between sidebar and topbar
 - [UI - Service Weather] Make all entities in the "More infos" modal grey when a watcher is on pbehavior
 - [UI - Text widget] Don't force a stat to be selected, so that this widget can be used as a simple text widget too
 - [UIv2 - Service Weather] Fixed links renderer function

## Canopsis 3.13.2 - Due date : 2019-03-26

 - [Go] Adding a description field in event filter rules
 - [Python] Adding a description field in event filter rules
 - [UI - CRUD Webhooks] Added authentification fields

## Canopsis 3.13.1 - Due date : 2019-03-26

 - VERSION.txt and CHANGELOG.md had not been updated for 3.13.0

## Canopsis 3.13.0 - Due date : 2019-03-25

 - [Documentation] Fixed documentation for go watchers
 - [Documentation] Revised documentation for webhooks
 - [Documentation] Fixed curl example in heartbeat documentation
 - [Go] Fixed flapping status in axe engine
 - [Go] Added webhook templates for automatic snoozing
 - [Go] Added event-filter patterns on alarms
 - [Go] Added possibility of using http and https proxy for axe
 - [Go] Fixed nested tickets handling in webhooks
 - [Go] Fixed missing certificates error
 - [Go] Added flapping periods counter for alarms
 - [Python] Fixed missing type error in watcherng route
 - [Python] Added hook field in action crud
 - [Python] Added flapping periods counter for alarms
 - [Python] Fixed pbehavior unwanted re-creation
 - [Python] Fixed pybars3 updating problem for centos7
 - [Python] Added configuration file for healthcheck api
 - [UI] Add tab's id into URL
 - [UI] Add rights on 'create widget' and 'create tab' actions
 - [UI] Fix - View display when double click
 - [UI] Refactor form mixin
 - [UI] Fix - Slide to first tab animation when entering a view
 - [UI] Refactor Vuex modules
 - [UI - Rights] Styles improvments
 - [UI - ListAlarm] Add a 'links' column
 - [UI - Context] Add success/failure popup when creation mass pbehaviors
 - [UI - Service weather] Add links renderer
 - [UI - Stats] Refactor stats widgets
 - [UI - Stats] Add customizable text widget
 - [UI - Exploitation] Add type selector on Exploitation forms
 - [UI - CRUD Pbehaviors] Styles improvments
 - [UI - CRUD Webhooks] Add an optionnal 'id' field
 - [UI - CRUD Event filter] Add 'id' and 'description' fields
 - [Packaging] Fix webserver compatibility with CentOS 7

## Canopsis 3.12.0 - Due date : 2019-03-11

 - [Documentation] Added documentation for common use cases
 - [Documentation] Added documentation for go watchers
 - [Documentation] Added documentation for action engine
 - [Go] Added watcher service
 - [Go] Added state change counters to alarms
 - [Go] Added snooze action to action engine
 - [Python] Added action_required field to the service weather API
 - [Python] Added API route to update pbehaviors
 - [Python] Added API for go watchers
 - [UI] Fixed pbehavior creation with RRULE
 - [UI] Integrated external libraries
 - [UI] Reverted addition of the active tab id into the URL (waiting for a bug fix on this functionnality)
 - [UI - Exploitation] Added Webhooks CRUD view
 - [UI - Exploitation] Improved style on Pbehaviors CRUD view
 - [UI - Alarms list] Added access to 'entity' variable on info popup template
 - [UI - Context] Fixed pbehavior mass deletion
 - [UI - Context] Added a right to manage access to 'Delete pbehavior' action
 - [UI - Service Weather] Changed tile blinking condition
 - [UI - Service Weather] Added rights

## Canopsis 3.11.0 - Due date : 2019-02-22

 - [Documentation] Add webhook documentation
 - [Docmentation] Add libreNMS documentation
 - [Go] Add watchers to the context graph and update their impacts and dependencies
 - [Python] Fix memory leak inside the snmp engine
 - [Python] Fix the bug that prevents the python-ldap package from being installed
 - [Python] Make snow2canopsis wait for the webserver to start
 - [Python] Fix performance issue with the stat API when there are too many entities
 - [Python] Allow the user to strip the string used in the email2canopsis templates
 - [Python] Fix a crash that prevents the correct rendering of the PBehaviors through the API
 - [UI] Show the active view on leftbar
 - [UI] Change Text editor library (from QuillJS to Jodit)
 - [UI] Fix bugs with filter editor + Add a variable type selector on it
 - [UI - Tabs] Fix bug with the actions buttons
 - [UI - Views] Add active tab into URL
 - [UI - Login] Redirect when access to view is forbidden
 - [UI - AlarmsList] Automatic filter when displaying 'resolved' alarms
 - [UI - AlarmsList] Fix Open/Resolved setting
 - [UI - AlarmsList] Optimize widget performance
 - [UI - Context] Bug - "More infos" tabs are invisible
 - [UI - Service Weather] Rights
 - [UI - Service Weather] Keep an entity expanded during a periodic refresh
 - [UI - Service Weather] Add a periodic refresh setting
 - [UI - Service Weather] Tooltips on actions
 - [UI - legacy - Weather] Send a correct event keepstate when the validate button is clicked

## Canopsis 3.10.0 - Due date : 2019-02-08

 - [Documentation] Update Service Weather documentation
 - [Documentation] Update available variables on customizable templates.
 - [Documentation] The Centreon connector is now open-source
 - [Documentation] Fix typos in pbehavior documentation
 - [Go] Add webhooks plugin
 - [Go] Fix issue with missing ack statistics
 - [Go] Add conversion of patterns into MongoDB requests
 - [Go] Add pattern list types
 - [Go] Add entities in calls to AxePostProcessor.ProcessAlarms
 - [Go] Add the possibility to declare a ticket without an event
 - [Python] Add support for exclusion dates in pbehaviors
 - [Python] Add webhook API
 - [Tooling] Fix push_docker_images to push the go engines' images
 - [Tooling] Remove python-ldap from docker images (WARNING: this will temporarily break support for ldap authentication, and should be fixed in the next release)
 - [Tooling] Add engine-axe-cat docker image
 - [Tooling] Fix SNMP MIBs import for CAT
 - [UI] Fix default views
 - [UI] Update to VueJS 2.5.21
 - [UI] Fix - View deletion message
 - [UI] Add optionnal "Comment" field on pbehavior creation
 - [UI] Improve sidebar views links
 - [UI] Fix - Variable type change, on filter editor's advanced mode
 - [UI] Fix - View duplication
 - [UI - UX] Fix - User profile menu's position
 - [UI - UX] Improve user profile menu's styling
 - [UI - Tabs] Add tab duplication functionnality
 - [UI - VueX] Unused entities automatic clean-up
 - [UI - Lodash] Add babel-plugin-lodash package
 - [UI - Service Weather] Add automatic refresh when filter setting has changed
 - [UI - Service Weather] Style improvments
 - [UI - Service Weather] Fix - Blinking conditions and Displayed icons
 - [UI - Service Weather] Add actions for actions done on entities

## Canopsis 3.9.0 - Due date : 2019-01-24

 - [Documentation] Add documentation for action API
 - [Tooling] Fix building scripts that used the wrong version of canopsis-next
 - [Go] Add support for MongoDB replicaset
 - [Go] Add a post-processing plugins system in the axe engine
 - [Go] Define triggers in the axe engine
 - [Python] Fix healthcheck API response when criticals is undefined
 - [Python] Fix pbehavior handling when the rrule generates invalid dates
 - [UI] Add "Help" buttons on ListAlarm/Context/Service weather widgets
 - [UI] Add "Reorder tabs" functionnality
 - [UI] Add "Refresh" button on all administration/exploitation views
 - [UI] Improve left sidebar style
 - [UI] Variables harmonization for all templates parameters (Info popup, More infos, Service weather)
 - [UI / Calendar] Fix display bug when only 1 filter is set
 - [UI / Calendar] Improve style (colors)
 - [UI / ListAlarm] Add "Delete" alarms mass action
 - [UI / ListAlarm] Refactor "Info popup" setting
 - [UI / ListAlarm] Fix "Resolved" column problem
 - [UI / ListAlarm] Improve "Info popup" style
 - [UI / Context] Add rights
 - [UI / Service weather] Fix bug on ListAlarm's modal's filter
 - [UI / Service weather] Add entities names customization, in "More info" modal

## Canopsis 3.8.0 - Due date : 2019-01-10

 - [Documentation] Clean up README.md
 - [Documentation] Add documentation for the tabs system
 - [Tooling] Push the init image in the push_docker_image.sh script
 - [Go] Refactor heartbeat engine and API
 - [Python] Fix timezone handling in pbehaviors
 - [Python] Add ack support in alarm filter
 - [UI] Delete lines name display
 - [UI] Fix pagination length bug
 - [UI] Add Pbehaviors CRUD view
 - [UI/ServiceWeather] Add actions on entities
 - [UI/ServiceWeather] Fix weather items blinking conditions
 - [UI/Calendar] Fix filter deletion bug
 - [UI/Context] Refactor "Manage infos" panel
 - [UI/Context] Improve "More infos" expand panel style

## Canopsis 3.7.0 - Due date : 2018-12-27

 - [Documentation] Add documentation for ackcentreon task
 - [Tooling] Fix debian installation method
 - [UI] Add tabs system inside views
 - [UI] Add default views feature for roles and users
 - [UI/ServiceWeather] Add a parameter that, when clicking on a service weather tile, allows choosing between opening a modal with the observer's entities list and opening an alarm list related to the observer
 - [UI/ServiceWeather] Fix random screen freezing
 - [UI/Alarm List] Fix bug that fetched the alarms twice
 - [UI/Context] Fix timestamp problem at pbehavior creation
 - [UI/Context] Add a search bar used to search through an entity's informations
 - [UI/Filters] Add a CRUD managing the event filter rules

## Canopsis 3.6.0 - Due date : 2018-12-13

 - [Documentation] Add documentation for global status check with the healthcheck route
 - [Documentation] Add documentation for mix-filters and entity duplication
 - [Documentation] Various cosmetic improvements
 - [Go] Fix requests to the Observer ticketing API
 - [Go] Fix statecounter steps handling
 - [Tooling] Fix error handling in canopsinit
 - [Tooling] Fix issue with unused parameters in init
 - [UI] Add default views and rights for UIv3
 - [UI/Context] Add "Manage infos" panel for watchers
 - [UI/Context] Fix a bug with resources expand
 - [UI/Context] Add "Clone" action on entities and watchers
 - [UI/Context] Filter required on watcher's creation (at least one valid rule in the filter)
 - [UI/Rights] Fix a bug with confirmation panel not closing when submitting rights
 - [UI/Events] Add an `origin: 'canopsis'` parameter with all events coming from Canopsis UI
 - [UI/Version] Add Canopsis version number on side-bar
 - [UI/Filters] Add "Mix filters" feature
 - [UI/Alarm list] Simplify default sort column selector in settings
 - [UI/Top bar] Fix a bug with group editing on the top bar

## Canopsis 3.5.0 - Due date : 2018-11-29

- [Documentation]: Add a new documentation
- [Python]: Add a new route to fetch a list of entities with their current alarm
- [Python]: Automatically recover from the loss of the primary member of a MongoDB replicaset
- [Go]: Prevent a crash when a snooze has no duration
- [Go]: Add an option to automatically acknowledge the resources below a component
- [Go]: End the implementation of the eventfilter service
- [Go]: Automatically create a ticket when a new alarm is created if the flag autoDeclareTickets is given to the axe engine
- [ServiceWeather]: Fix the message set in the events sent when an action is triggered
- [Tooling]: Update the configuration of catag to handle the new canopsis project
- [Tooling]: Add a VERSION.txt file inside canopsis, display it on the prompt inside the canopsis env and add an API to retrieve it through HTTP
- [Tooling]: Add the missing engine-action in the push_docker_image.sh script

## Canopsis 3.4.0 - Due date : 2018-11-15

- [Go]: Fix make release command
- [Go]: IncidentKey in Service Ticket API is now optional
- [Go]: Introducing new Observer API driver for ticket creation
- [Go]: New actions for event-filter engine (copy and set_entity_info_from_template)
- [Python]: Fix /get-alarms route which was limited to 50 elements, and returned an overestimation of total
- [Python]: New Healthcheck API to check Canopsis status (service connections, engine status...)
- [Python]: Update Heartbeat engine docstring
- [ServiceWeather]: Fix paused pbehavior icon not always correctly display

## Canopsis 3.3.0 - Due date : 2018-10-31

- [Alarm List]: Fix bug where selected filter is not saved in userpreferences
- [Chore]: Some configuration cleanup in sources, configs and docker files
- [Docker]: Change rabbitmq base image and add some new envvars for image version selection
- [email2canopsis]: Converter now invalidate a translation if there is not match
- [Go]: Change build dir (from /tmp to ./build)
- [Go]: New event-filter engine ! (to filter/translate events based on some rules)
- [Go]: New ticket engine ! (to create tickets through external APIs)
- [Go]: Standardization of heartbeat ids generation
- [Python]: Add new event-filter and ticket API (for new engines)
- [ServiceWeather]: Add default message on Validate and Ack actions
- [snow2canopsis]: New connector which read Service Now API to import informations into Canopsis context
- [UI]: Fix canopsis-next not correctly builded with docker

## Canopsis 3.2.5 - Due date : 2018-10-17

- [Doc]: translating recent Upgrade documentations to english
- [Docker] : new CPS_LOGGING_LEVEL envvar to change loglevel in dockerized engines
- [email2canopsis] : now decode encoded subject line
- [Go]: new Action engine (especially pbehaviors from regex)
- [Go]: send alarm resolution informations to stats engine
- [Python] : fix "Socket Error 104" while engines communicate with rabbitmq
- [Tool] : env2cfg can now handle mongo replicaset option
- [UI] : add rights for items in engine menu
- [UI] : fix `has_active_pb` flag not correctly calculated
- [UI] : fix timeline not correctly showing all informations (especially on ack, ticket and canceled events)
- [UI] : integrate listalarm, timeline and querybuilder bricks in central repo

## Canopsis 3.2.4 - Due date : 2018-10-05

- [Chore]: add some system dependencies to ansible installation
- [Go]: fix bagoting alarms never closed if cropped
- [Go]: send axe statistics to statsng engine
- [Python]: change some amqp publishers to pika to prevent odd reconnections
- [Python]: fix has\_active\_pb flag doesn't correctly show all linked pbehaviors
- [Python]: fix pbehavior doesn't correctly handle timezone change (one day gap)
- [Python]: fix performance concern on alert consultation

## Canopsis 3.2.3 - Due date : 2018-09-17

- [Email Connector]: handle base64 encoded parts
- [Email Connector]: option which use redis to resend last known event on run
- [Service weather]: Fix icons consistency between a watcher and his modal

## Canopsis 3.2.2 - Due date : 2018-09-12

- [Go]: Fix che event enrichment with new entities
- [Go]: Fix unused debug flag
- [Python]: Fix pbehavior based desynchronisation on weather route

## Canopsis 3.2.1 - Due date : 2018-09-12

- [CAT] Simplification of the statistics API (!74, !77, !78)
- [CAT] Add trends, SLA, sorting, aggregations and filtering by author to the statistics API  (!75, !76, !80, !87)
- [CAT] Add support for monthly periods to the statistics API (!79)
- [CAT] Add statistics for ongoing alarms and current state of an entity (!81, 82)
- [CAT] Fix an issue with the state of the alarms for statistics on resolved alarms (!83)
- [CAT] Add API route to get the history of the state of an entity (!85)
- [Python]: Fix error on /trap route
- [Python]: Fix amqp driver to revive stucked Consumers
- [Service weather]: Fix wrong pbehavior and maintenance icons used
- [UI]: Fix role default view that cannot be modified

## Canopsis 3.2 - Due date : 2018-09-01

- [Python]: fix rruled pbehavior computation
- [Python]: two bugfix on /get-entities route
- [Service weather]: fix unfoldable item when it contains a %
- [UI]: add ellipsis on hostgroups field
- [UI]: add rights on filters (create, list, modify)
- [UI]: hide "Restore alarm" button on open alarm list

## Canopsis 3.1 - Due date : 2018-08-17

- [Python]: add a feature that track the change of the longoutput field and alter the behavior of the output field of an alarm.
- [Python]: fix a bug that duplicate an alarm.
- [Build]: fix an error during the compilation of canopsis-next.
- [UI]: fix a right issue in the listalarm brick that allow the done button to be displayed without any right.

## Canopsis 3.0.1 - Due date : 2018-08-03

- [CAT]: New stats in statsng: most_alarms_impacting, worst_mtbf, alarm_list, longest_alarms (!67 et !70)
- [CAT]: Add a "periods" parameter to the statsng API (!69)
- [CAT]: New stats in statsng: most_alarms_created, alarms_impacting, state_list (!71)
- [CAT]: send stats with alarm duration
- [CAT]: default sats value
- [Snooze]: faster snooze application
- [Periodic refresh]: fix infinite loop with periodic refresh
- [Link list]: remove linklist
- [Links]: fix links in service weather
- [UI]: views in canopsis-next
- [UI]: alarms in canopsis-next

## Canopsis 3.0.0 - Due date : 2018-07-20

- [Chore]: Finalizing new building process
- [CAT]: New stats in statsng: time_in_state, availability, maintenance and mtbf
- [CAT]: Provisionning for docker demos
- [Deploy]: Automated standalone deployment
- [Docker]: Multiple cleanup and fixes
- [Go]: Fix an uncatched exception in alarm
- [Go]: Large refacto on engines and services
- [Go]: New engine interface
- [Go]: New Done action
- [Go]: Support watcher
- [Python]: Huge code cleanup, upgrade to pymongo 3.6
- [Python]: New Done action
- [Python]: New metric engine
- [Python]: New stat event: statstateinterval
- [Python]: New webserver handler using Flask
- [Python]: Small refacto and optimization on linkbuilder classes
- [UI]: Beginning Canopsis-next integration (context view)
- [UI]: Fix brickmanager when used with canopsis user
- [UI]: New documentation subprocess (visible on search field in Alarm list)
- [UI]: New 'unknown' state marker
- [Alarm list]: Natural search only on visible columns
- [Alarm list]: Possibility to search on ticket number
- [Service weather]: Add rights for actions
- [Service weather]: Fix blink on pbehaviored watchers
- [Service weather]: Fix refresh (again)

## Canopsis 2.7.0 - Due date : 2018-06-28

- [API]: new routes to manage frontend views
- [CAT]: add routes and make some bugfixes for the new engine statsng
- [CAT]: fix priviledge escalation and introduction of a default right group
- [Go]: large service/adapter refacto
- [Go]: fix timeout on mongo requests
- [Service weather]: fixed an issue where periodical refresh process never refresh anymore
- [Service weather]: add permission on action
- [UI]: a new 'done' action is avalaible, which simply mark an alarm as done
- [UI]: fix /sessionstart route
- [UI]: fix snmp view that can be broken
- [Webserver]: authentification trough WebSSO

## Canopsis 2.6.8 - Due date : 2018-06-20

- [Service weather]: fixed an issue where the popup overlay could stay when the popup was closed, freezing the view
- [CAT]: ported the email2canopsis connector to python3 and fixed an issue with pattern matching caused by the new python version.
- [pbehaviors]: refactored the pbheaviors internal API to unify processing and avoid inconsistencies

## Canopsis 2.6.7 - not released

This release was replaced by the 2.6.8 version to integrate an urgent fix on the pbehaviors.

## Canopsis 2.6.6 - Due date : 2018-06-08

- [Connector]: email2canopsis now with python3
- [Docker]: add some usefull tools in debian-9 docker image
- [EventFilter]: now can filter on subkeys
- [Go]: fix crash when a non-check event arrive on an empty alarm
- [Go]: add printEventOnError flag on engines
- [Go]: update dependency, and notablly migrating from mgo to globalsign/mgo
- [Go]: fix event-entity not always correctly associated
- [Go]: fix PBehavior event message never acked in Che
- [Go]: fix inaccurate messages in Che
- [Python]: API can now produce correct HttpError on catched exception (instead of 200)
- [Python]: limit /get-alarms response size by removing steps
- [Python]: upsert mode on import context route
- [Python]: alarms without a corresponding entity are not shadowed anymore
- [Python]: fix active pbehavior search
- [Python]: two new events (statcounterinc and statduration) for statistics purpose
- [Python]: remove ticket engine from amqp2engines.conf
- [UI]: fix BNF research on entity properties
- [UI]: fix sort on current_state_duration
- [Service Weather]: add new rights for actions in service weather

## Canopsis 2.6.5 - Due date : 2018-05-18

- [Go]: Porting steps cropping from python
- [Go]: Huge refactoring on error handling
- [Go]: Events with a timestamp as a float are now correctly parsed
- [Python]: step cropping compatibility update
- [Python]: bugfix where a source_type component is treated as a resource
- [UI]: fix context info presentation in the context explorer
- [UI]: bugfix on hide resource
- [UI]: fix navbar rights
- [Service Weather]: fix possibility to ack twice, change invalidate to cancel on thumb down button
- [Alarms list]: creation date is now presented in full format, even if its the present day

## Canopsis 2.6.4 - Due date : 2018-05-03

- [Setup] A version of the amqp2engines.conf file is provided for the High performance engines
- [Go]: fix error code when an engine crashes
- [Go]: pbehavior events are now correctly parsed by Che and passed to the pbehavior engine
- [Go]: generation of an unique id on every alarm
- [Go]: fix a crash when impacts and depends are not initialized in an entity object
- [Go]: Removed some useless logs from the Context
- [Event filter]: fixed an issue that blocked snooze events
- [Heartbeat]: fix heartbeat never closed on special conditions
- [Stat]: perfomance update on the engine
- [Stat]: fix handling special character in requests
- [UI]: fix presentation of pbehavior (in service-weather), ellipsis and snooze
- [UI]: fix permissions on the "List pbehaviors" action
- [CAT]: fix a crash with datametrie engine
- [Alarms list]: the  'pbehavior' icon is now displayed only when a pbeheavior is active on the entity, not when it is configured
- [Alarms list]: add the snooze End date on the Snooze tooltip
- [Alarms list]: several fixes on the ACK and report ticket action
- [Alarms list]: Pbehaviors can now be filtered.
- [Alarms list]: fix an issue where a pbehavior id could become null with periodic refresh enabled
- [Service weather]: Resources are now greyed out when a Component has a pbehavior
- [Context Graph]: Fixed an issue that prevented the expansion of an item of the Context Graph explorer

## Canopsis 2.6.3 - Due date : 2018-04-26

- [Engines] : Added UNACK, Uncancel, keep state actions to the High performance engines
- [Engines] : Added a default author in all alarms steps in the High performance engines
- [Go]: gracefull initialization of docker env and fix amqp2engines.conf on docker
- [Go]: Context and alarm creation performance improvements
- [Go]: code refactoring, fix alarm duplication and preloading
- [UI]: update rights on default view
- Multiple fixes on Pbheavior and Event_filter

## Canopsis 2.6.2 - Due date : 2018-04-23

### Fixes

- [Alarms list] Fixed a date formatting issue in the alarms list that made the `last_update_date` column appear with a 1 month delay
- [Service weather] reworked the Ticket action to fix a display issue caused by the new "save on exit" workflow
- [Alarms list] Mass actions now correctly get their rights/permissions (inherited from the rights applied on the actions on the single alarm)
- [setup] Canopsinit now requires a flag `--authorize-reinit` to perform any destructive modification to the database as an extra security
- [APIs] the "enabled" flag on all entities is now active
- [Engines] removed alarms caching in the Che engine to avoid  alarms duplication
- [Rights management] : fixed a rights issue with the massive actions on a limited account
- [Service weather] : fixed the components display to put long names on 2 lines instead of truncating it

## Canopsis 2.6.1 - Due date: 2018-04-20

**Not released due to regression**

## Canopsis 2.6.0 - Due date: 2018-04-18

This release introduced the new High performance engines and allowed the renaming on the Alarms list columns.

### Functional changes

- [Alarms list] : columns can now be renamed by the user
- [Alarms list] : pbehaviors can now be filtered
- [Alarms list] : added rights management to some new actions that were missing it
- [Alarms list] : Added "duration" and "current_state_duration" columns that display the total duration of the alarm and the duration of the current state of the alarm (respectively)
- [Alarms list] : the date is now displayed even if the the alarm was created today
- [pbehaviors] : the form has been simplified to give a behavior closer to a Calendar event
- [setup] : replaced the `schema2db` and `canopsis-filldb` commands with the new Canopsinit command (**Warning: see UPGRADING_2.6.md document**)
- [CAT] : The Datametrie connector can now filter some alarms based on their criticity
- [CAT] : The datametrie connector can now use the local date

### Experimental features

- [Engines] : New High Performance engines for heavily loaded environments (**experimental**)
- [Engines] : reimplemented the last_event_date feature on the High performance engines
- [engines] : The engine "stats" (High performance version) can now log the actions and their autors for audit purposes

### Bug fixes

- [Alarms list] : fixed an issue where rights were not saved properly in the admin
- [Alarms list] : fixed an issue that prevented pbehaviors to be saved properly on the High performance engines
- [Alarms list] : fixed an issue that could record a ticket number with the `0` value with the High performance engines
- [Engines] fixed an issue where the cancel Action did not close the alarm with the High performance engines

## Canopsis 2.5.12 (Sprint 03.16) - Due date : 2018-03-16

### Functional changes

- [Service Weather] Added the Fast ACK action in the service Weather widget
- [Service Weather] Added an ACK icon in the service weather titles, to identify entities that have acknowledged alarms
- [Service Weather] Added OK/NOK events statistics in the info popup
- [APIs] APIs can now be requested using Basic Auth, without having to request the authentication route first
- [Timeline] automatic actions now have an Author: "system"
- [list alarms] added "Cancel alarm" and "Change criticity" actions, with correct translation
- [Event filter] added "state" and "state_type" fields in the usable fields list
- [Installation] New Installation method based on RPM/deb packages

### Bug fixes

- [Service weather] Fixed an issue where the watcher state was incorrectly impacted by paused applications
- [Service weather] Fixed a type issue that could prevent the UI to display correctly
- [HA Tools] Fixed an issue of the High Performance engines when the RabbitMQ connection was lost
- [HA Tools] Reduced the downtime of the Canopsis UI when the MongoDB primary instance changes
- [CAT] Fixed an SAMLv2 installation issue on CentOS 7/RHEL 7
- [CAT] Fixed an race condition where SNMP rules could be missing when a trap was received while the engine's rules update function was running
- [CAT] Fixed an issue where an SNMP trap with a routing key > 255 chars could crash the SNMP engine and block the whole AMQP queue

## Canopsis 2.5.6 (sprint 02.2) - Due date : 2018-02-02

### Bug fixes
 - [#598](https://git.canopsis.net/canopsis/canopsis/issues/598) - Certaines entités ne sont pas importées correctement
 - [#593](https://git.canopsis.net/canopsis/canopsis/issues/593) - Pouvoir interroger les alarmes Resolved avec les dates des alarmes
 - [#589](https://git.canopsis.net/canopsis/canopsis/issues/589) - Probleme Dataset sur widget ServiceWeather
 - [#580](https://git.canopsis.net/canopsis/canopsis/issues/580) - Stealthy calculation
 - [#568](https://git.canopsis.net/canopsis/canopsis/issues/568) - Modifier la couleur d’affichage de la plate-forme de pré-prod
 - [#565](https://git.canopsis.net/canopsis/canopsis/issues/565) - Le bouton "Remove alarm (poubelle/corbeille)" n'apparaît pas sur les alarmes critiques
 - [#558](https://git.canopsis.net/canopsis/canopsis/issues/558) - Recherche avec des int
 - [#529](https://git.canopsis.net/canopsis/canopsis/issues/529) - Pouvoir supprimer des alarmes
 - [#528](https://git.canopsis.net/canopsis/canopsis/issues/528) - Disposer d'une information sur la date de l'alarme

### Functional and other changes

 - [#599](https://git.canopsis.net/canopsis/canopsis/issues/599) - Nettoyage des engines
 - [#566](https://git.canopsis.net/canopsis/canopsis/issues/566) - Remapper l'output "Lost 100%" en "Equipement injoignable"
 - [#594](https://git.canopsis.net/canopsis/canopsis/issues/594) - Validation des actions dans le popup MDS

## Canopsis 2.5.5 (sprint 01.19) - Due date : 2018-01-19

### Bug fixes

 - [#579](https://git.canopsis.net/canopsis/canopsis/issues/579) - Impossible de créer un pbehavior dans l'explorateur de context
 - [#564](https://git.canopsis.net/canopsis/canopsis/issues/564) - [API] get-alarms ne remonte pas tous les résultats
 - [#563](https://git.canopsis.net/canopsis/canopsis/issues/563) - [HardLimit] la hardlimit empêche toutes les actions sur une alarme
 - [#485](https://git.canopsis.net/canopsis/canopsis/issues/485) - bac a alarme création de pbhavior pop up calendrier qui ne s'affiche pas

### Functional and other changes

 - [#573](https://git.canopsis.net/canopsis/canopsis/issues/573) - Bac à alarme - recherche insensible à la casse
 - [#526](https://git.canopsis.net/canopsis/canopsis/issues/526) - Pouvoir trier les tuiles du widget service weather
 - [#525](https://git.canopsis.net/canopsis/canopsis/issues/525) - [Météo] Remonter les statistiques d'un scénario
 - [#524](https://git.canopsis.net/canopsis/canopsis/issues/524) - [Météo de service]Disposer d'une information sur la date de l'alarme

## 2.5.4 (Sprint 01.6) - Due date : 2018-01-06

### Bug fixes

 - [#543](https://git.canopsis.net/canopsis/canopsis/issues/543) - [Engine Alerts] Lorsqu'une alerte a atteint sa hard limit, le beat processing plante dans check_alarm_filters
 - [#540](https://git.canopsis.net/canopsis/canopsis/issues/540) - [doc] lancer un filldb update à chaque mise à jour de Canopsis
 - [#538](https://git.canopsis.net/canopsis/canopsis/issues/538) - [docker/CAT] La brique SNMP n'est pas installée
 - [#537](https://git.canopsis.net/canopsis/canopsis/issues/537) - [CRUD Context] Le schéma d'édition d'une entité n'est pas bon
 - [#556](https://git.canopsis.net/canopsis/canopsis/issues/556) - Declare ticket: missing number
 - [#555](https://git.canopsis.net/canopsis/canopsis/issues/555) - [PE] Dysfonctionnement du workflow des alarmes filter
 - [#553](https://git.canopsis.net/canopsis/canopsis/issues/553) - Missing display_name
 - [#548](https://git.canopsis.net/canopsis/canopsis/issues/548) - Recherche naturelle non fonctionnelle sur les display name
 - [#546](https://git.canopsis.net/canopsis/canopsis/issues/546) - Probleme de couleur sur les tuiles en MAJOR
 - [#545](https://git.canopsis.net/canopsis/canopsis/issues/545) - Pb de désynchronisation de statut entre scénario et Application
 - [#544](https://git.canopsis.net/canopsis/canopsis/issues/544) - Pb de bagot sur les alarmes. Elles ne se cloturent plus.
 - [#542](https://git.canopsis.net/canopsis/canopsis/issues/542) - Pbehavior crash when filter is in a bad format
 - [#539](https://git.canopsis.net/canopsis/canopsis/issues/539) - [Bac à alarmes] Des pbehaviors expirés remontent sur les alarmes en cours
 - [#531](https://git.canopsis.net/canopsis/canopsis/issues/531) - Problème avec la recherche naturelle
 - [#513](https://git.canopsis.net/canopsis/canopsis/issues/513) - le formulaire de login  doit catch le 401

### Functional and other changes

 - [#516](https://git.canopsis.net/canopsis/canopsis/issues/516) - Dummy authentication
 - [#536](https://git.canopsis.net/canopsis/canopsis/issues/536) - Ajouter le résutlat des tests dans le template
 - [#535](https://git.canopsis.net/canopsis/canopsis/issues/535) - Ajouter prérequis dans le template
 - [#533](https://git.canopsis.net/canopsis/canopsis/issues/533) - Retirer le bouton PAUSE sur les alarmes CLOSED
 - [#541](https://git.canopsis.net/canopsis/canopsis/issues/541) - Nettoyage des test

## Canopsis 2.5.3 (Sprint 12.15) Due date : 2017-12-15

**Not released due to blocking issue. This release was tagged on Gitlab but not distrubuted. All issues were reported in the 2.5.4 release**

## Canopsis 2.5.2 (11.30) - Due date : 2017-12-05

### Bug fixes

 - [#499](https://git.canopsis.net/canNombre d'issues pour la milestone  "2.5.1 (11.03)" : 11
opsis/canopsis/issues/499) - [PE] la météo ne s'affiche pas pour les applications "standard"
 - [#518](https://git.canopsis.net/canopsis/canopsis/issues/518) - Saml2 import error
 - [#515](https://git.canopsis.net/canopsis/canopsis/issues/515) - Erreur 500 avec avec un identifiants inconnu
 - [#509](https://git.canopsis.net/canopsis/canopsis/issues/509) - Erreur de login webserver
 - [#497](https://git.canopsis.net/canopsis/canopsis/issues/497) - attribut creation_date  en avance par rapport aux dates des évènements sur certaines alarmes.
 - [#450](https://git.canopsis.net/canopsis/canopsis/issues/450) - Probleme de fonctionnement sur la recherche dans la bac à alarmes

### Functional and other changes

  - [#510](https://git.canopsis.net/canopsis/canopsis/issues/510) - Ajouter une date de dernier événement reçu à chaque alarme.
  - [#466](https://git.canopsis.net/canopsis/canopsis/issues/466) - Pouvoir disposer d'id d'alarmes exploitable dans le bac à alarmes

## Canopsis 2.5.1 (11.03) - Due date : 2017-11-03

### Bug fixes

 - [#446](https://git.canopsis.net/canopsis/canopsis/issues/446) - Le paquet dm.xmlsec.binding-1.3.3.tar.gz ne s'installe pas sur CentOS 7
 - [#427](https://git.canopsis.net/canopsis/canopsis/issues/427) - watcher ne se recalcule pas a la fin d'un pbehavior => desync
 - [#383](https://git.canopsis.net/canopsis/canopsis/issues/383) - soucis utf-8 sur les trap laposte
 - [#372](https://git.canopsis.net/canopsis/canopsis/issues/372) - Pas de données sur la route /perfdata
 - [#470](https://git.canopsis.net/canopsis/canopsis/issues/470) - Problème de performance de la route trap
 - [#451](https://git.canopsis.net/canopsis/canopsis/issues/451) - La configuration du popup du widget alarme ne conserve pas les données
  - [#399](https://git.canopsis.net/canopsis/canopsis/issues/399) - probleme utf8  quand on met un display name avec un accent
 - [#432](https://git.canopsis.net/canopsis/canopsis/issues/432) - UI ne charge pas
 - [#430](https://git.canopsis.net/canopsis/canopsis/issues/430) - bug(CRUD context-graph): Crudcontext adapter en doublon
 - [#413](https://git.canopsis.net/canopsis/canopsis/issues/413) - Probleme d'affichage sur le bac à alarmes
 - [#412](https://git.canopsis.net/canopsis/canopsis/issues/412) - Formatage des heures dans l'historique du bac à alarmes

 - [#401](https://git.canopsis.net/canopsis/canopsis/issues/401) - [Bac à alarmes] Lorsqu'un filtre est présent, le premier chargement de la vue échoue.

### Functional and other changes

 - [#472](https://git.canopsis.net/canopsis/canopsis/issues/472) - MAJ Image Docker MongoDB
 - [#406](https://git.canopsis.net/canopsis/canopsis/issues/406) - Installation mono-package python
 - [#478](https://git.canopsis.net/canopsis/canopsis/issues/478) - [Alarms] Optimisation des perfs du beat processing
 - [#357](https://git.canopsis.net/canopsis/canopsis/issues/357) - feat(context-graph): Basic editing of an entity
 - [#356](https://git.canopsis.net/canopsis/canopsis/issues/356) - feat(CRUD context-graph) : list entities
 - [#336](https://git.canopsis.net/canopsis/canopsis/issues/336) - Remplacer les Configurable par confng
 - [#403](https://git.canopsis.net/canopsis/canopsis/issues/403) - [Bac à alarmes] Recherche naturelle
 - [#402](https://git.canopsis.net/canopsis/canopsis/issues/402) - [Météo de service] Ajouter un bouton FastACK sur la popup "scénario"

## Canopsis 2.4.6 and CAT 2.5.0 (25/09/2017)

Canopsis 2.4.6 is a maintenance release for the 2.4 branch of Canopsis.

### Functional changes - CAT

- [#393](https://git.canopsis.net/canopsis/canopsis/issues/393) Feat(Auth) : Compatibilité SAMLV2
- [#375](https://git.canopsis.net/canopsis/canopsis/issues/375) Feat(SNMP) : Les traps SNMP anomalies ne remontent pas

### Bug fixes - CAT

- [#375](https://git.canopsis.net/canopsis/canopsis/issues/375) Fix(SNMP): Les traps SNMP anomalies ne remontent pas

### Functional  and other changes

- [#394](https://git.canopsis.net/canopsis/canopsis/issues/394) feat(UI) : permettre l'ajout d'onglets dropdown dans la vue Header
- [#392](https://git.canopsis.net/canopsis/canopsis/issues/392) feat(Context-graph) : création d'une route ws  pour update du  context

### Bug fixes and other non-functional changes

- [#391](https://git.canopsis.net/canopsis/canopsis/issues/391) fix(Context-Graph) : La route post retourne parfois des doublons
- [#378](https://git.canopsis.net/canopsis/canopsis/issues/378) fix(web): Blocage appli
- [#377](https://git.canopsis.net/canopsis/canopsis/issues/377) fix(Météo de service) : Impossible de faire les actions sur les alarmes
- [#376](https://git.canopsis.net/canopsis/canopsis/issues/376) fix(Météo de service) : La mise en pause d'un scénario en alarme ne remet pas en vert l'application
- [#374](https://git.canopsis.net/canopsis/canopsis/issues/374) fix(Bac à Alarmes) : Désynchro entre une alarme et son historique & fermeture de l'alarme pas toujours prise en compte
- [#349](https://git.canopsis.net/canopsis/canopsis/issues/349) fix(Météo de service) : Nouvelle desynchro entre entités sur 2.4.5
- [#347](https://git.canopsis.net/canopsis/canopsis/issues/347) fix(Météo de service) : création de pbehavior depuis le service weather
- [#371](https://git.canopsis.net/canopsis/canopsis/issues/371) fix(perfs) : Perf sur queue alerts en 2.4.5
- [#364](https://git.canopsis.net/canopsis/canopsis/issues/364) fix(pbehavior) : Engine pbehavior, déconnection et reconnection en boucle
- [#342](https://git.canopsis.net/canopsis/canopsis/issues/342) fix(pbehavior) : création lente
- [#333](https://git.canopsis.net/canopsis/canopsis/issues/333) fix(Bac à Alarmes) : les boutons d'action de masse ne fonctionnent que si 1 alarme est sélectionnée
- [#326](https://git.canopsis.net/canopsis/canopsis/issues/326) fix(métriques) : Route /api/context/metric et recherches par nom
- [#320](https://git.canopsis.net/canopsis/canopsis/issues/320) fix(pbehavior): Modale pbehavior - pas de création de pbheavior
- [#318](https://git.canopsis.net/canopsis/canopsis/issues/318) fix(pbehavior) : Il est possible de créer un pbehavior avec une rrule invalide
- [#317](https://git.canopsis.net/canopsis/canopsis/issues/317) fix(pbehavior) : check des rrules avant insertion

## Canopsis 2.4.5 (25/08/2017)

### Functional changes

- feat(Météo de service): amélioration serviceweather hauteur des tuiles
- [#345](https://git.canopsis.net/canopsis/canopsis/issues/345) trad(Météo de service) : Traduction française de la météo de Service
- [#337](https://git.canopsis.net/canopsis/canopsis/issues/337) feat(Météo de service) : afficher un compteur de temps avant le prochain changement sur une tuile en alarme
- [#323](https://git.canopsis.net/canopsis/canopsis/issues/323) feat(Météo de service) : templatiser les modales du widget
- [#314](https://git.canopsis.net/canopsis/canopsis/issues/314) feat(Bac à alarmes) : pouvoir afficher les infos de la resource affectée par l'alarme
- [#309](https://git.canopsis.net/canopsis/canopsis/issues/309) feat(Bac à alarmes) : ajouter dans les alarmes un champ de la cause d'alarme
- [#305](https://git.canopsis.net/canopsis/canopsis/issues/305) feat(baselines) : Intégration des baselines avec le service weather
- [#295](https://git.canopsis.net/canopsis/canopsis/issues/295) feat(alarmes) :  ajout de dates de création et de date de dernier changement dans une alarme

### Bug fixes and other non-functional changes:

- [#344](https://git.canopsis.net/canopsis/canopsis/issues/344) fix(engines) : Lors de la création d'un pbehavior avec le widget service weather, le cleaner event crash.
- [#339](https://git.canopsis.net/canopsis/canopsis/issues/339) fix(Mété de service) :  Erreur 404 sur api/v2/weather/watcher
- [#338](https://git.canopsis.net/canopsis/canopsis/issues/338) fix(session) : storage non instancié
- [#330](https://git.canopsis.net/canopsis/canopsis/issues/330) fix(Météo de service) : PBehaviors non fonctionnels
- [#329](https://git.canopsis.net/canopsis/canopsis/issues/329) fix(Bac à alarmes) : Le champs actions ne fonctionne pas si la colonne extra_details n'est pas affichée
- [#327](https://git.canopsis.net/canopsis/canopsis/issues/327) fix(Météo de service) : mixin customsendevent : conflit entre versions service-weather et listalarms
- [#324](https://git.canopsis.net/canopsis/canopsis/issues/324) fix(Bac à alarmes) : probleme de Timestamp
- [#310](https://git.canopsis.net/canopsis/canopsis/issues/310) fix(metrics) :  changer les adapters pour récupérer les métriques
- [#298](https://git.canopsis.net/canopsis/canopsis/issues/298) fix(pbehaviors): crash de l'engine pbehaviors lors de la creation de pbehavior |
- [#296](https://git.canopsis.net/canopsis/canopsis/issues/296) fix(runtime): amqp2engines\* dans hypcontrol si on se trouve dans un dossier avec un amqp2engines.conf qui existe
- [#294](https://git.canopsis.net/canopsis/canopsis/issues/294) fix(global) : Internationalisation non fonctionnelle
- [#293](https://git.canopsis.net/canopsis/canopsis/issues/293) fix(Bac à alarmes) : La recherche ne fonctionne pas sur la 2.4.4
- [#291](https://git.canopsis.net/canopsis/canopsis/issues/291) fix(Bac à alarmes) : Le tri automatique sur les dates ne fonctionne pas
- [#290](https://git.canopsis.net/canopsis/canopsis/issues/290) fix(Bac à alarmes) : Pas de rafaichissement automatique de la vue lorsqu'une action est effectuée sur une alarme
- [#289](https://git.canopsis.net/canopsis/canopsis/issues/299) fix(snooze) : fonctionnalité inopérante
- [#287](https://git.canopsis.net/canopsis/canopsis/issues/287) refact(configuration): [interne] remplacer Configurable

## Canopsis 2.4.0

- feat(context-graph) : new database structure that stores the real-world topology of a supervised system as an object graph, allowing us to identify all entities impacted by an alarm
- feat (backend) : Canopsis now generates alarmes based on the status of events. This new object keep tracks of the full history of a real-world alarm
- feat(UI) : new "Service weather" UI brick that can display the status of up to 120 entities on a single window
- feat (UI): new "Alarms list" UI brick as a replacement of the old events list.
