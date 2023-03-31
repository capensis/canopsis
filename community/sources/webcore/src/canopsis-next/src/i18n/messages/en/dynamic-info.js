export default {
  advancedSearch: '<span>Help on the advanced research :</span>\n'
    + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
    + '<p>The "-" before the research is required</p>\n'
    + '<p>Operators :\n'
    + '    <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
    + '<p>For querying patterns, use "pattern" keyword as the &lt;ColumnName&gt; alias</p>\n'
    + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
    + '<dl><dt>Examples :</dt><dt>- description = "testdyninfo"</dt>\n'
    + '    <dd>Dynamic info rules descriptions are "testdyninfo"</dd><dt>- pattern = "SEARCHPATTERN1"</dt>\n'
    + '    <dd>Dynamic info rules whose one of its patterns is equal "SEARCHPATTERN1"</dd><dt>- pattern LIKE "SEARCHPATTERN2"</dt>\n'
    + '    <dd>Dynamic info rules whose one of its patterns contains "SEARCHPATTERN2"</dd>'
    + '</dl>',
};
