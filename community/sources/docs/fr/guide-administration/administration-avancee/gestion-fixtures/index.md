# Gestion des fixtures (données d’initialisation)

Cette section décrit la gestion des fixtures dans Canopsis, c’est-à-dire les jeux de données d’initialisation permettant de préconfigurer l’application.
Vous y trouverez la procédure d’export pour sauvegarder la configuration actuelle de votre Canopsis et celle d’import pour restaurer ces données sur une nouvelle instance.

- [Export](export.md)
  
- [Import](import.md)

!!! Important

    Par défaut, la commande reconfigure de Canopsis est fournie avec des fixtures prédéfinies.

    Si vous ajoutez un fichier personnalisé de fixtures dans le répertoire dédié, contenant uniquement certaines collections, il est possible que certains liens entre celles-ci soient rompus en raison d’IDs autogénérés différents d’une instance à l’autre.

    Nous vous recommandons donc d’exporter systématiquement l’ensemble des collections ou de vous assurer qu’aucun lien n’existe entre la collection exportée et d’autres collections associées.

    Pour plus d’informations, n’hésitez pas à contacter les équipes de Capensis.

