import get from './get';
import date from './date';
import duration from './duration';
import json from './json';
import percentage from './percentage';

// Stats filters
import { trendColor as statsTrendColor, trendIcon as statsTrendIcon } from './stats/trend';
import statsFormatValue from './stats/formatted-value';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', date);
    Vue.filter('duration', duration);
    Vue.filter('json', json);
    Vue.filter('percentage', percentage);
    // Stats filters
    Vue.filter('statsTrendColor', statsTrendColor);
    Vue.filter('statsTrendIcon', statsTrendIcon);
    Vue.filter('statsFormatValue', statsFormatValue);
  },
};
