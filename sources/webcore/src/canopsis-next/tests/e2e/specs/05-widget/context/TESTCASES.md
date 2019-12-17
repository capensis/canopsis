# Test cases for context widget  

## Searching

* [ ] **The search founds a relevant result.**
    1. To input any expected data to the search field in the upper left corner.
    2. To click on button with magnifier.
        
    **Expected result**: the rows with relevant results are shown in the table.

* [ ] **The search can be canceled.**
    1. To input any data to the search field.
    2. To click on the cross icon.
        
    **Expected result**: a search field became clear, a table has default view.

## Pagination

* [ ] **The right arrow at the top of the widget opens the next page.**
    1. To press on right arrow button at the top of the widget.
        
    **Expected result**: the next page is opened.

* [ ] **The left arrow at the top of the widget opens the previous page.**
    1. To press on right arrow button at the top of the widget.
    2. To press on left arrow button at the top of the widget.
        
    **Expected result**: the previous page is opened.

## Filtering

* [ ] **A filter can be selected.**
    1. To select filter in the dropdown at the top of the page.
        
    **Expected result**: only relevant results are shown in the table.

* [ ] **Filter can be changed to another.**
    1. To select filter in the dropdown at the top of the page.
    2. To select another filter.

    **Expected result**: only results suitable to the last filter are shown in the table.

* [ ] **Filter can be canceled.**
    1. To select filter in the dropdown at the top of the page.
    2. To click on the cross near the selected filter.

    **Expected result**: the filter is canceled.

* [ ] **The "AND" option of "Mix filters" works correctly.**
    1. To select any filter.
    2. To turn on "Mix filters".
    3. To select "AND" option in check-box.
    4. To select the second filter.
        
    **Expected result**: the results in table satisfy the conditions of both filters at the same time.

* [ ] **The "OR" option of "Mix filters" works correctly.**
    1. To select any filter.
    2. To turn on "Mix filters".
    3. To select OR option in check-box.
    4. To select the second filter.
        
    **Expected result**: the results in table satisfy the condition of at least one filter.

## Items

* [ ] **The button with “plus” opens a menu with two options.**
    1. To press on the button with “plus” at the upper right corner of the page.
        
    **Expected result**: a menu with the buttons “Entities” and “Watcher” is shown.

* [ ] **An entity can be created.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Entities”.
    3. To fill in all required fields.
    4. To press the button “Submit”.
        
    **Expected result**: an entity is created, a notification “Entity successfully created !” at the upper right corner is appeared.

* [ ] **An entity with additional information can be created.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Entities”.
    3. To fill in all selected field in the tab “Form”.
    4. To open tab “Manage infos”.
    5. To press the button with “plus”.
    6. To fill in all required fields in the appeared modal window.
    7. To press the button “Add”.
    8. To press the button “Submit”.
        
    **Expected result**: a notification “Entity successfully created !” at the upper right corner is appeared, an entity with additional information is created.

* [ ] **A block of information about entity can be edited.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Entities”.
    3. To fill in all selected field in the tab “Form”.
    4. To open tab “Manage infos”.
    5. To press the button with “plus”.
    6. To fill in all required fields in the appeared modal window.
    7. To press the button “Add”.
    8. To press the button with pencil.
    9. To change any fields in the modal window “Edit an information”.
    10. To press the button “Add”.
        
    **Expected result**: a string with information about entity is updated according to the new values.

* [ ] **A block of information about entity can be deleted.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Entities”.
    3. To fill in all selected field in the tab “Form”.
    4. To open tab “Manage infos”.
    5. To press the button with “plus”.
    6. To fill in all required fields in the appeared modal window.
    7. To press the button “Add”.
    8. To press the button with trash bin.
        
    **Expected result**: a string with information about entity is deleted.

* [ ] **A watcher can be created.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Watcher”.
    3. To input name in the tab “Form”.
    4. To select a rule in the tab “Filter”.
    5. To press the button “Submit”.
        
    **Expected result**: a watcher is created, a notification “Watcher is successfully created !” is shown.

* [ ] **A watcher with additional information can be created.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Watcher”.
    3. To input name in the tab “Form”.
    4. To select a rule in the tab “Filter”.
    5. To press the “plus” button in the tab “Manage infos”.
    6. To fill in all required fields in the appeared modal window.
    7. To press the button “Add”.
    8. To press the button “Submit”.
        
    **Expected result**: a watcher with additional information is created, a notification “Watcher is successfully created !” is shown.

* [ ] **A block of information about watcher can be edited.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Watcher”.
    3. To input name in the tab “Form”.
    4. To select a rule in the tab “Filter”.
    5. To press the “plus” button in the tab “Manage infos”.
    6. To fill in all required fields in the appeared modal window.
    7. To press the button “Add”.
    8. To press the button with pencil.
    9. To change any fields in the modal window “Edit an information”.
    10. To press the button “Add”.
        
    **Expected result**: a string with information about watcher is updated according to the new values.

* [ ] **A block of information about watcher can be deleted.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Watcher”.
    3. To input name in the tab “Form”.
    4. To select a rule in the tab “Filter”.
    5. To press the “plus” button in the tab “Manage infos”.
    6. To fill in all required fields in the appeared modal window.
    7. To press the button “Add”.
    8. To press the button with pencil.
    9. To press the button with trash bin.
        
    **Expected result**: a string with information about watcher is deleted.

* [ ] **A block of information about watcher can be deleted.**
    1. To click on button with “plus” to open a menu.
    2. To select the “Watcher”.
    3. To input name in the tab “Form”.
    4. To select a rule in the tab “Filter”.
    5. To press the “plus” button in the tab “Manage infos”.
    6. To fill in all required fields in the appeared modal window.
    7. To press the button “Add”.
    8. To press the button with pencil.
    9. To press the button with trash bin.
        
    **Expected result**: a string with information about watcher is deleted.

## Mass actions

* [ ] **All entities can be selected.**
    1. To select a check-box at the header of the table.
        
    **Expected result**: all rows are selected.

* [ ] **The only one entity can be selected.**
    1. To select a check-box of any entity.
        
    **Expected result**: the entity is selected.
