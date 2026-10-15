export default {
  title: 'Interface utilisateur',
  appTitle: 'Titre de l\'application',
  language: 'Langue par défaut',
  footer: 'Page d\'identification : pied de page',
  description: 'Page d\'identification : description',
  loginPage: 'Page d\'identification',
  matchedItems: 'Éléments correspondants',
  animations: 'Animations',
  instructionApprovement: 'Approbation des consignes',
  advanced: 'Avancé',
  logo: 'Logo',
  infoPopupTimeout: 'Délai d\'affichage pour les popups d\'informations',
  errorPopupTimeout: 'Délai d\'affichage pour les popups d\'erreurs',
  allowChangeSeverityToInfo: 'Autorise le changement de criticité en Info',
  showHeaderOnKioskMode: 'Afficher l\'en-tête en mode kiosque',
  showHeaderOnKioskModeTooltip: 'L\'en-tête ne sera pas disponible en mode kiosque plein écran',
  maxMatchedItems: 'Seuil d\'éléments avant avertissement',
  checkCountRequestTimeout: 'Délai d\'expiration de la requête',
  requiredInstructionApprove: 'Approbation des consignes requise',
  disabledTransitions: 'Désactivez certaines animations',
  disabledTransitionsTooltip: 'Cela aidera à améliorer les performances de l\'application',
  autoSuggestPbehaviorName: 'Proposer automatiquement un nom de comportement périodique',
  defaultTheme: 'Thème par défaut',
  dateFormat: 'Format de date et heure',
  dateFormats: {
    dateTime: '31/12/2000 00:00:00',
    dateTimeWithWeekday: '31/12/2000 (lun.) 00:00:00',
    weekdayWithDateTime: 'lun. 31/12/2000 00:00:00',
  },
  versionDescriptionTooltip: 'Info-bulle contextuelle sur la bannière de version',
  defaultVersionDescription: `<div>
  <div>Édition:<strong> {'{{ edition }}'}</strong></div>
  <div>Nom de série:<strong> {'{{ serialName }}'}</strong></div>
  <div>Version mise à jour:<strong> {'{{ timestamp versionUpdated }}'}</strong></div>
</div>`,
  tooltips: {
    dateFormat:
      'Ce format sera appliqué à toutes les colonnes et tous les champs de date/horodatage dans Canopsis',
    maxMatchedItems:
      'Avertit l\'utilisateur lorsque le nombre d\'éléments correspondant aux modèles est supérieur à cette valeur',
    checkCountRequestTimeout:
      'Définit le délai d\'expiration (en secondes) de la demande d\'éléments correspondants',
  },
};
