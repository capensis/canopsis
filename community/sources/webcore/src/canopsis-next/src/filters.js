import { get } from 'lodash';

import { convertDurationFormToMaxUnitDurationString, convertDurationToString } from '@/helpers/date/duration';
import { convertDateToTimezoneDateString, convertDateToString, convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { convertNumberToFixedString, convertNumberToRoundedPercentString } from '@/helpers/string';
import { stringifyJsonFilter } from '@/helpers/json';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', convertDateToString);
    Vue.filter('dateWithToday', convertDateToStringWithFormatForToday);
    Vue.filter('duration', convertDurationToString);
    Vue.filter('json', stringifyJsonFilter);
    Vue.filter('percentage', convertNumberToRoundedPercentString);
    Vue.filter('maxDurationByUnit', convertDurationFormToMaxUnitDurationString);
    Vue.filter('timezone', convertDateToTimezoneDateString);
    Vue.filter('fixed', convertNumberToFixedString);
  },
};
