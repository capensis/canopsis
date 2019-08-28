import get from './get';
import date from './date';
import duration from './duration';
import json from './json';
import percentage from './percentage';

// Pbehaviors
import pbehaviorToForm from './pbehavior/pbehavior-to-form';
import pbehaviorToComments from './pbehavior/pbehavior-to-comments';
import pbehaviorToExdate from './pbehavior/pbehavior-to-exdate';
import formToPbehavior from './pbehavior/form-to-pbehavior';
import commentsToPbehaviorComments from './pbehavior/comments-to-pbehavior-comments';
import exdateToPbehaviorExdate from './pbehavior/exdate-to-pbehavior-exdate';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', date);
    Vue.filter('duration', duration);
    Vue.filter('json', json);
    Vue.filter('percentage', percentage);
    // Pbehaviors
    Vue.filter('pbehaviorToForm', pbehaviorToForm);
    Vue.filter('pbehaviorToComments', pbehaviorToComments);
    Vue.filter('pbehaviorToExdate', pbehaviorToExdate);
    Vue.filter('formToPbehavior', formToPbehavior);
    Vue.filter('commentsToPbehaviorComments', commentsToPbehaviorComments);
    Vue.filter('exdateToPbehaviorExdate', exdateToPbehaviorExdate);
  },
};
