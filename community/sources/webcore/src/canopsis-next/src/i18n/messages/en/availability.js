import { AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

export default {
  filterByValue: 'Filter by value',
  valueFilterMethods: {
    [AVAILABILITY_VALUE_FILTER_METHODS.greater]: 'Greater than',
    [AVAILABILITY_VALUE_FILTER_METHODS.less]: 'Less than',
  },
  popups: {
    exportCSVFailed: 'Failed to export availabilities in CSV format',
  },
};
