import { PATTERN_TYPES } from '@/constants';

export default {
  patterns: 'Modèles',
  myPatterns: 'Mes modèles',
  corporatePatterns: 'Modèles partagés',
  addRule: 'Ajouter une règle',
  addGroup: 'Ajouter un groupe',
  removeRule: 'Supprimer la règle',
  advancedEditor: 'Éditeur avancé',
  simpleEditor: 'Éditeur simple',
  noData: 'Aucun modèle. Cliquez sur \'@:pattern.addGroup\' pour ajouter des champs au modèle',
  noDataDisabled: 'Aucun modèle.',
  discard: 'Effacer le motif',
  oldPatternTooltip: 'Les modèles de filtre ne sont pas migrés',
  types: {
    [PATTERN_TYPES.alarm]: 'Modèle d\'alarme',
    [PATTERN_TYPES.entity]: 'Modèle d\'entité',
    [PATTERN_TYPES.pbehavior]: 'Modèle de comportements périodiques',
  },
  errors: {
    ruleRequired: 'Veuillez ajouter au moins une règle',
    groupRequired: 'Veuillez ajouter au moins un groupe',
    invalidPatterns: 'Les modèles ne sont pas valides ou il y a un champ de modèle désactivé',
    countOverLimit: 'Le modèle que vous avez défini cible {count} éléments. Cela peut affecter les performances, en êtes-vous sûr ?',
    oldPattern: 'Le modèle de filtre actuel est défini dans l\'ancien format. Veuillez utiliser l\'éditeur avancé pour l\'afficher. Les filtres dans l\'ancien format seront bientôt obsolètes. Veuillez créer de nouveaux modèles dans notre interface mise à jour.',
    existExcluded: 'Les règles incluent la règle exclue.',
    required: 'Au moins un modèle doit être défini. Veuillez définir des modèles de filtre pour la règle',
  },
};
