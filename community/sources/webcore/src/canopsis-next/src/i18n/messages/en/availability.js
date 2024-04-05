import { AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

export default {
  filterByValue: 'Filter by value',
  advancedSearch: '<span>Help on the advanced research :</span>\n'
    + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
    + '<p>The "-" before the research is required</p>\n'
    + '<p>Operators :\n'
    + '    <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
    + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
    + '<dl><dt>Examples :</dt><dt>- Name = "name_1"</dt>\n'
    + '    <dd>Entities whose names are "name_1"</dd><dt>- Name="name_1" AND Type="service"</dt>\n'
    + '    <dd>Entities whose names is "name_1" and the types is "service"</dd><dt>- infos.custom="Custom value" OR Type="resource"</dt>\n'
    + '    <dd>Entities whose infos.custom is "Custom value" or the type is "resource"</dd><dt>- infos.custom LIKE 1 OR infos.custom LIKE 2</dt>\n'
    + '    <dd>Entities whose infos.custom contains 1 or 2</dd><dt>- NOT Name = "name_1"</dt>\n'
    + '    <dd>Entities whose name isn\'t "name_1"</dd>\n'
    + '</dl>',
  valueFilterMethods: {
    [AVAILABILITY_VALUE_FILTER_METHODS.greater]: 'Greater than',
    [AVAILABILITY_VALUE_FILTER_METHODS.less]: 'Less than',
  },
  popups: {
    exportCSVFailed: 'Failed to export availabilities in CSV format',
  },
};
