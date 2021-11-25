import { Identifier, Constants } from 'dayspan';
import DaySpanVuetify from 'dayspan-vuetify/src/plugin';

import DsDay from './components/day.vue';
import DsDayTimes from './components/day-times.vue';
import DsCalendar from './components/calendar.vue';
import DsCalendarApp from './components/calendar-app.vue';
import DsCalendarEvent from './components/calendar-event.vue';
import DsCalendarEventPlaceholder from './components/calendar-event-placeholder.vue';
import DsCalendarEventTimePlaceholder from './components/calendar-event-time-placeholder.vue';
import DsCalendarEventPopover from './components/calendar-event-popover.vue';
import DsCalendarEventTime from './components/calendar-event-time.vue';
import DsWeekDayHeader from './components/week-day-header.vue';

/* eslint-disable no-proto, no-param-reassign */
/**
 * We've added this code because dayspan is using moment().day() for dayOfWeek. It's incorrect solution for `fr` locale.
 * And we've replaced it to moment().weekday()
 */
const originalGet = Identifier.Day.__proto__.get;

Identifier.Day.__proto__.get = function get(context) {
  context.dayOfWeek = context.date.weekday();

  return originalGet.call(this, context);
};

/**
 * We've added this code to using `seconds` duration inside a schedule
 */
Constants.DURATION_TO_MILLIS = {
  second: Constants.MILLIS_IN_SECOND,
  seconds: Constants.MILLIS_IN_SECOND,

  ...Constants.DURATION_TO_MILLIS,
};
/* eslint-enable no-proto, no-param-reassign */

export default {
  install(Vue, options = {}) {
    Vue.use(DaySpanVuetify, options);
    Vue.component('dsDay', DsDay);
    Vue.component('dsDayTimes', DsDayTimes);
    Vue.component('dsCalendar', DsCalendar);
    Vue.component('dsCalendarApp', DsCalendarApp);
    Vue.component('dsCalendarEvent', DsCalendarEvent);
    Vue.component('dsCalendarEventPlaceholder', DsCalendarEventPlaceholder);
    Vue.component('dsCalendarEventTimePlaceholder', DsCalendarEventTimePlaceholder);
    Vue.component('dsCalendarEventPopover', DsCalendarEventPopover);
    Vue.component('dsCalendarEventTime', DsCalendarEventTime);
    Vue.component('dsWeekDayHeader', DsWeekDayHeader);
  },
};
