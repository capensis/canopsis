import { LINK_RULE_TYPES } from '@/constants';

export default {
  simpleMode: 'Mode simple',
  advancedMode: 'Mode avancé',
  addLink: 'Ajouter un lien',
  linksEmpty: 'Aucun lien ajouté pour le moment',
  linksEmptyError: 'Vous devez ajouter au moins 1 lien en mode simple ou modifier le code source en mode avancé',
  sourceCodeAlert: 'Veuillez ne modifier ce script que si vous êtes parfaitement conscient de ce que vous faites',
  type: 'Type de lien',
  single: 'Appliquer ce lien uniquement à une seule alarme ?',
  types: {
    [LINK_RULE_TYPES.alarm]: 'Alarme',
    [LINK_RULE_TYPES.entity]: 'Entité',
  },
};
