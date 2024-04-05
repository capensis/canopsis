import { AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

export default {
  filterByValue: 'Filtrer par valeur',
  advancedSearch: '<span>Aide sur la recherche avancée :</span>\n'
    + '<p>- [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt;</p> [ AND|OR [ NOT ] &lt;NomColonne&gt; &lt;Opérateur&gt; &lt;Valeur&gt; ]\n'
    + '<p>Le "-" avant la recherche est obligatoire</p>\n'
    + '<p>Opérateurs:\n'
    + '    <=, <,=, !=,>=, >, LIKE (Pour les expressions régulières MongoDB)</p>\n'
    + '<p>Les types de valeurs : Chaîne de caractères entre guillemets doubles, Booléen ("TRUE", "FALSE"), Entier, Nombre flottant, "NULL"</p>\n'
    + '<dl><dt>Exemples :</dt><dt>- Name = "name_1"</dt>\n'
    + '    <dd>Entités dont le names est "name_1"</dd><dt>- Name="name_1" AND Type="service"</dt>\n'
    + '    <dd>Entités dont le names est "name_1" et la types est "service"</dd><dt>- infos.custom="Custom value" OR Type="resource"</dt>\n'
    + '    <dd>Entités dont le infos.custom est "Custom value" ou la type est "resource"</dd><dt>- infos.custom LIKE 1 OR infos.custom LIKE 2</dt>\n'
    + '    <dd>Entités dont le infos.custom contient 1 or 2</dd><dt>- NOT Name = "name_1"</dt>\n'
    + '    <dd>Entités dont le name n\'est pas "name_1"</dd>\n'
    + '</dl>',
  valueFilterMethods: {
    [AVAILABILITY_VALUE_FILTER_METHODS.greater]: 'Plus grand que',
    [AVAILABILITY_VALUE_FILTER_METHODS.less]: 'Moins que',
  },
  popups: {
    exportCSVFailed: 'Échec de l\'exportation des disponibilités au format CSV',
  },
};
