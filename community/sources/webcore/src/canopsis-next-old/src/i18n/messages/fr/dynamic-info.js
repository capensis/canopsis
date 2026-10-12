import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

export default {
  tabs: {
    dynamicInfos: 'Informations dynamiques',
    templates: 'Modèles',
  },
  createFromTemplateTooltip: 'Créer une information dynamique à partir du modèle',
  templatesList: {
    massRemove: 'Supprimer les modèles d\'informations dynamiques',
  },
  massRemove: 'Supprimer les informations dynamiques',
  massEnable: 'Activer les informations dynamiques',
  massDisable: 'Désactiver les informations dynamiques',
  informationTypes: {
    [DYNAMIC_INFO_INFORMATION_TYPES.setToInfo]: 'Définir une valeur constante',
    [DYNAMIC_INFO_INFORMATION_TYPES.setToInfoFromTemplate]: 'Définir une valeur avec un modèle',
    [DYNAMIC_INFO_INFORMATION_TYPES.copyToInfo]: 'Copier la valeur depuis une entité / alarme / événement',
  },
};
