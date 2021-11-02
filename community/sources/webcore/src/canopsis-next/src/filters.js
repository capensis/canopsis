import { get } from 'lodash';

import { convertDurationToMaxUnitDurationString, convertDurationToString } from '@/helpers/date/duration';
import { convertDateToTimezoneDateString, convertDateToString, convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { convertNumberToFixedString, convertNumberToRoundedPercentString } from '@/helpers/string';
import { stringifyJson } from '@/helpers/json';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', convertDateToString);
    Vue.filter('dateWithToday', convertDateToStringWithFormatForToday);
    Vue.filter('duration', convertDurationToString);
    Vue.filter('json', stringifyJson);
    Vue.filter('percentage', convertNumberToRoundedPercentString);
    Vue.filter('maxDurationByUnit', convertDurationToMaxUnitDurationString);
    Vue.filter('timezone', convertDateToTimezoneDateString);
    Vue.filter('fixed', convertNumberToFixedString);
  },
};
