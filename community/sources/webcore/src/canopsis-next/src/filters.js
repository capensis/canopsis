import { get } from 'lodash';

import { convertDurationFormToMaxUnitDurationString, durationToString } from '@/helpers/date/duration';
import { convertDateToTimezoneDateString, convertDateToString } from '@/helpers/date/date';
import { convertNumberToFixedString, convertNumberToRoundedPercentString } from '@/helpers/string';
import { stringifyJson } from '@/helpers/entities';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', convertDateToString);
    Vue.filter('duration', durationToString);
    Vue.filter('json', stringifyJson);
    Vue.filter('percentage', convertNumberToRoundedPercentString);
    Vue.filter('maxDurationByUnit', convertDurationFormToMaxUnitDurationString);
    Vue.filter('timezone', convertDateToTimezoneDateString);
    Vue.filter('fixed', convertNumberToFixedString);
  },
};
