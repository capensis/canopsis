import { AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

export default {
  filterByValue: 'Filtrer par valeur',
  valueFilterMethods: {
    [AVAILABILITY_VALUE_FILTER_METHODS.greater]: 'Plus grand que',
    [AVAILABILITY_VALUE_FILTER_METHODS.less]: 'Moins que',
  },
  popups: {
    exportCSVFailed: 'Échec de l\'exportation des disponibilités au format CSV',
  },
};
