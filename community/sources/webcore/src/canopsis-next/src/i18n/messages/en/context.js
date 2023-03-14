export default {
  impacts: 'Impacts',
  dependencies: 'Dependencies',
  noEventsFilter: 'No events filter',
  impactChain: 'Impact chain',
  resolvedAlarms: 'Resolved alarms',
  activeAlarm: 'Active alarm',
  impactDepends: 'Impact/Depends',
  treeOfDependencies: 'Tree of dependencies',
  infosSearchLabel: 'Search infos',
  eventStatisticsMessage: '{ok} OK events\n{ko} KO Events',
  eventStatistics: 'Event statistics',
  addSelection: 'Add selection',
  advancedSearch: '<span>Help on the advanced research :</span>\n'
    + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
    + '<p>The "-" before the research is required</p>\n'
    + '<p>Operators :\n'
    + '    <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
    + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
    + '<dl><dt>Examples :</dt><dt>- Name = "name_1"</dt>\n'
    + '    <dd>Entities whose names are "name_1"</dd><dt>- Name="name_1" AND Type="service"</dt>\n'
    + '    <dd>Entities whose names is "name_1" and the types is "service"</dd><dt>- infos.custom.value="Custom value" OR Type="resource"</dt>\n'
    + '    <dd>Entities whose infos.custom.value is "Custom value" or the type is "resource"</dd><dt>- infos.custom.value LIKE 1 OR infos.custom.value LIKE 2</dt>\n'
    + '    <dd>Entities whose infos.custom.value contains 1 or 2</dd><dt>- NOT Name = "name_1"</dt>\n'
    + '    <dd>Entities whose name isn\'t "name_1"</dd>\n'
    + '</dl>',
  actions: {
    titles: {
      editEntity: 'Edit entity',
      duplicateEntity: 'Duplicate entity',
      deleteEntity: 'Delete entity',
      pbehavior: 'Periodical behavior',
      variablesHelp: 'List of available variables',
      massEnable: 'Enable entities',
      massDisable: 'Disable entities',
    },
  },
  fab: {
    common: 'Add a new entity',
    addService: 'Add a new service entity',
  },
  popups: {
    massDeleteWarning: 'The mass deletion cannot be applied for some of selected elements, so they won\'t be deleted.',
  },
};
