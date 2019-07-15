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
│   └── elementCount.js      # Special global assertion
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
│       └── auth.js          # File with tests. If we want to order files we can put number prefix here.
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
If we need set the order of tests running we can put number prefixes for tests files.
