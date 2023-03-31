export default {
  advancedSearch: '<span>Aide sur la recherche avancée :</span>\n'
    + '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n'
    + '<p>Le "-" avant la recherche est obligatoire</p>\n'
    + '<p>Opérateurs:\n'
    + '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n'
    + '<p>Pour effectuer une recherche dans les "modèles", utilisez le mot-clé "pattern" comme &lt;NomColonne&gt;</p>\n'
    + '<p>Les types de valeurs : Chaîne de caractères entre guillemets doubles, Booléen ("TRUE", "FALSE"), Entier, Nombre flottant, "NULL"</p>\n'
    + '<dl><dt>Examples :</dt><dt>- description = "testdyninfo"</dt>\n'
    + '    <dd>Règles d\'Informations dynamiques dont la description est "testdyninfo"</dd><dt>- pattern = "SEARCHPATTERN1"</dt>\n'
    + '    <dd>Règles d\'Informations dynamiques dont un des modèles est "SEARCHPATTERN1"</dd><dt>- pattern LIKE "SEARCHPATTERN2"</dt>\n'
    + '    <dd>Règles d\'Informations dynamiques dont un des modèles contient "SEARCHPATTERN2"</dd>'
    + '</dl>',
};
