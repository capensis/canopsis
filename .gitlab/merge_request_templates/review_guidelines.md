# Merge Request Canopsis

<message>

## Pré-requis

``` 
Indiquez ici les pré-requis (s'il y en a) pour que la merge request puisse être testée et validée

```

## Pour la validation

- [ ] Le code respecte la PEP8
- [ ] Chaque classe/fonction/méthode contient une docstring complète
- [ ] Le code est revu, les modifications demandées sont appliquées
- [ ] Les messages de commits sont compatibles avec l' [angular spec](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#commit) ?
- [ ] Les modifications apportées contiennent des tests unitaires
- [ ] Les tests unitaires passent
- [ ] (En cas de changement de comportement) : la documentation a été mise à jour

## Tests


1. Installer la branche dans un environnement Canopsis
2. se mettre dans l'environnement Canopsis :`su - canopsis`
3. (re)démarrer Canopsis :`hypcontrol restart`
4. lancer les tests unitaires :  `ut_runner`
5. Retourner dans le dossier où les sources de canopsis sont installées :`cd /vagrant/canopsis`
6. Exécuter les tests fonctionnels : `cd sources/functional_testing/ && python2 runner.py`
