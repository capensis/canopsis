# Format des expressions régulières Canopsis

Dans Canopsis, l'opérateur `regex_match` est disponible dans différentes interfaces, API et moteurs, afin de réaliser une condition si et seulement si un motif d'expression régulière (ou *regex*) est vrai.

Le document suivant présente les formats d'expressions régulières pris en charge par Canopsis.

!!! attention
    Canopsis ajoute une différence par rapport aux formats d'expressions régulières décrits ci-dessous : les opérateurs commençant par un antislash (tels que `\b`, `\w`…) **doivent être échappés une seconde fois** (ce qui donne `\\b`, `\\w`…).

## Format par défaut et format avancé

### Format d'expression régulière par défaut : regex Golang `re2`

Il s'agit du format utilisé par défaut par Canopsis, en raison de sa meilleure vitesse d'exécution et de sa consommation mémoire réduite. Notez que ce format est différent du format « PCRE » plus couramment utilisé.

Ce format est détaillé en anglais dans [le document de syntaxe `re2`](https://github.com/google/re2/wiki/Syntax) de Google.

Les lignes marquées `NOT SUPPORTED` dans le document précédent ne sont ainsi pas prises en charge par ce format. Il n'est pas non plus possible de réaliser des « expressions régulières négatives », ou *backreferences* avec ce format. Seul le format d'expression régulière « PCRE / .NET » décrit dans la section suivante permet de réaliser ces actions.

### Format d'expression régulière avancé : regex PCRE / .NET

Ce format d'expression régulière est utilisé uniquement lorsque votre expression ne se conforme pas au format Golang `re2` décrit précédemment. Canopsis se charge automatiquement de tenter à nouveau l'évaluation de votre `regex_match` au format PCRE / .NET si l'évaluation est impossible au format Golang `re2`.

Ce format avancé ajoute notamment la possibilité de réaliser :

*  des *backreferences* ;
*  des expressions régulières négatives ;
*  etc.

Voyez la [documentation des expressions régulières .NET](https://docs.microsoft.com/fr-fr/dotnet/standard/base-types/regular-expression-language-quick-reference) de Microsoft pour en savoir plus sur les possibilités de ce format avancé.

!!! important
    Les expressions régulières de type PCRE / .NET peuvent nécessiter beaucoup plus de ressources que les expressions régulières Golang, à leur exécution, ce qui peut fortement réduire votre taux de traitement d'évènements Canopsis.

    La [variable d'environnement `REGEXP2_MATCH_TIMEOUT`](../../guide-administration/administration-avancee/variables-environnement.md) permet d'accorder un temps d'exécution maximal à une expression régulière. La durée maximale par défaut est d'une seconde.

    Pour cette raison, il est recommandé de n'avoir recours à ces expressions régulières avancées que lorsque votre cas d'utilisation le justifie vraiment.

## Tester vos expressions régulières en ligne

Le service en ligne [regex101.com](https://regex101.com) peut vous permettre de simuler l'exécution d'une expression régulière sur un ensemble de données de test. Vos tests doivent être exécutés sur les options « Golang » ou « PCRE », en prenant en compte les informations et différences décrites dans ce document.
