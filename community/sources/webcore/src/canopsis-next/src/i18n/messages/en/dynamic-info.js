import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

export default {
  informationTypes: {
    [DYNAMIC_INFO_INFORMATION_TYPES.setToInfo]: 'Set constant value',
    [DYNAMIC_INFO_INFORMATION_TYPES.setToInfoFromTemplate]: 'Set value with template',
    [DYNAMIC_INFO_INFORMATION_TYPES.copyToInfo]: 'Copy value from entity / alarm / event',
  },
};
