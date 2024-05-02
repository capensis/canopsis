import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

export default {
  isEnabled: 'Est actif',
  begins: 'Débute',
  ends: 'Se termine',
  lastAlarmDate: 'Date de la dernière alarme',
  alarmCount: 'Nombre d\'alarmes',
  massRemove: 'Supprimer les comportements périodiques',
  massEnable: 'Activer les comportements périodiques',
  massDisable: 'Désactiver les comportements périodiques',
  periodsCalendar: 'Calendrier avec périodes',
  notEditable: 'Ne peut pas être modifié',
  pbehaviorInfo: 'PBehavior infos',
  pbehaviorType: 'Type PBehavior',
  pbehaviorReason: 'Raison PBehavior',
  pbehaviorName: 'Nom PBehavior',
  pbehaviorCanonicalType: 'Type canonique de PBehavior',
  rruleEnd: 'Fin de récurrence',
  buttons: {
    addFilter: 'Ajouter un filtre',
    editFilter: 'Modifier le filtre',
    addRRule: 'Ajouter une règle de récurrence',
    editRrule: 'Modifier la règle de récurrence',
  },
  tabs: {
    type: 'Type',
    reason: 'Raison',
    exceptions: 'Dates d\'exception',
  },

  exceptions: {
    title: 'Dates d\'exception',
    create: 'Ajouter une date d\'exception',
    choose: 'Sélectionnez la liste d\'exclusion',
    usingException: 'Ne peut pas être supprimé car il est en cours d\'utilisation.',
    emptyExceptions: 'Aucune exception ajoutée pour le moment.',
  },

  types: {
    usingType: 'Le type ne peut être supprimé car il est en cours d\'utilisation.',
    defaultType: 'Le type est par défaut, vous ne pouvez modifier que le champ de couleur.',
    hidden: 'Masquer ce type sur le formulaire de comportement ?',
    types: {
      [PBEHAVIOR_TYPE_TYPES.active]: 'Actif',
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactif',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Maintenance',
    },
  },

  reasons: {
    usingReason: 'La raison ne peut pas être supprimée car elle est en cours d\'utilisation.',
    hidden: 'Masquez cette raison sur le formulaire de comportement ?',
  },
};
