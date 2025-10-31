# Mode TV (ou Kiosque) 

Le mode TV (ou mode Kiosque) est un mode disponible dans les vues de Canopsis, il permet de changer la manière dont les informations sont affichés à l'écran pour être adapté à des écrans de supervision (visible par plusieurs opérateurs sur le même site. Par exemple: dans un service informatique ou un service de supervision).

Le mode TV possède 3 type différents de vue :  
- [Plein écran](#kiosque--plein-écran) 
- [Kiosque uniquement](#kiosque-uniquement) 
- [Kiosque + Plein écran](#kiosque--plein-écran) 

Il existe deux possibilités d'accéder au mode kiosque.
- Depuis la bac à alarme en passant en mode édition ou en utilisant les raccourcis
- En allant sur l'URL: `https://[canopsis]/kiosk-views/<ID de la vue>/<ID de l'onglet>` *(Ces deux informations peuvent être trouver dans l'URL lorsque vous êtes sur le bac à alarme)

### Configurer le mode kiosque du widget

Pour configurer le mode kiosque, il faut se rendre dans le menu d'édition du bac à alarmes (CTRL + E) puis sur en appuyant sur les 3 points `...` dans l'onglet `Vues` puis `Mode kiosque`.

3 options sont disponibles:
| Option                  | Utilisation                                                                            |
|-------------------------|----------------------------------------------------------------------------------------|
| Masquer les actions | Permet de faire disparaître la colonne des actions rapides.                                 |
| Masquer la sélection de masse     | Permet de faire disparaître la colonne de sélections de masses.                      |
| Masquer la barre des tâches        | Permet de faire disparaître la barre des tâches |


### Affichage par défaut
*Raccourcis: Alt/CMD + Shift + 1*

![Affichage par défaut](./img/kiosque.png)

### Plein écran
*Raccourcis: Alt/CMD + Shift + 2*

Le mode plein écran permet d'afficher la vue en cours en plein écran en désaffichant la barre de navigation du haut ainsi que la barre latéral. Ce mode ne prend pas en compte les éléments configuré dans les paramètres du mode kiosque.

![Mode plein écran](./img/kiosque1.png)

### Kiosque uniquement
*Raccourcis: Alt/CMD + Shift + 3*

Le mode kiosque permet d'afficher la vue en cours en appliquant la configuation choisie par bac à alarme dans le menu `Edition > Vues > Mode kiosque`.

![Mode plein écran](./img/kiosque2.png)


### Kiosque + Plein écran
*Raccourcis: Alt/CMD + Shift + 4*

Le mode kiosque en plein écran permet d'afficher la vue en cours en en plein écran en appliquant la configuation choisie par bac à alarme dans le menu `Edition > Vues > Mode kioskque`.

![Mode plein écran](./img/kiosque3.png)

### En tête de Canopsis

Par défaut, l'en-tête de Canopsis est désactivé, pour le faire apparaître il faut se rendre dans la partie `Administration > Paramètres` puis activé l'option `Afficher l'en-tête en mode kiosque`.

Il est important de noter que cette option n'est pas disponible pour le mode `Kiosque + Plein écran`.

