import { durationToString } from '@/helpers/date/duration';
import { convertDateToString } from '@/helpers/date/date';

import get from './get';
import json from './json';
import percentage from './percentage';
import maxDurationByUnit from './max-duration-by-unit';
import timezone from './timezone';
import fixed from './fixed';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', convertDateToString);
    Vue.filter('duration', durationToString);
    Vue.filter('json', json);
    Vue.filter('percentage', percentage);
    Vue.filter('maxDurationByUnit', maxDurationByUnit);
    Vue.filter('timezone', timezone);
    Vue.filter('fixed', fixed);
  },
};
