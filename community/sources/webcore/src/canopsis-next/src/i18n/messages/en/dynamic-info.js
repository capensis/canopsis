import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

export default {
  tabs: {
    dynamicInfos: 'Dynamic infos',
    templates: 'Templates',
  },
  createFromTemplateTooltip: 'Create dynamic information from template',
  templatesList: {
    massRemove: 'Remove dynamic info templates',
  },
  massRemove: 'Remove dynamic infos',
  massEnable: 'Enable dynamic infos',
  massDisable: 'Disable dynamic infos',
  informationTypes: {
    [DYNAMIC_INFO_INFORMATION_TYPES.setToInfo]: 'Set constant value',
    [DYNAMIC_INFO_INFORMATION_TYPES.setToInfoFromTemplate]: 'Set value with template',
    [DYNAMIC_INFO_INFORMATION_TYPES.copyToInfo]: 'Copy value from entity / alarm / event',
  },
};
