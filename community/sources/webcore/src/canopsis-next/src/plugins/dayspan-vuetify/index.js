import { Calendar, Identifier, Constants } from 'dayspan';

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

/**
 * The same method with `Day.get` except one thing. We've fixed bug with weekday which was described above.
 *
 * @param context
 * @returns {*}
 */
Identifier.Day.__proto__.get = function get(context) {
  context.dayOfWeek = context.date.weekday();

  return originalGet.call(this, context);
};

/**
 * The same method with `refreshEvents` except one thing. This method is working asynchronously and returns Promise.
 *
 * @returns {Promise<Calendar>}
 */
Calendar.prototype.refreshEventsAsync = async function refreshEventsAsync() {
  const days = this.iterateDays().list();
  const eventsDaysEntries = [];
  const iteration = async (day) => {
    if (day.inCalendar || this.eventsOutside) {
      eventsDaysEntries.push([day, this.eventsForDay(day, this.listTimes, this.repeatCovers)]);
    }
  };

  await new Promise((resolve) => {
    for (let i = 0; i < days.length; i += 1) {
      setTimeout(() => iteration(days[i]), 0);
    }

    setTimeout(resolve, 0);
  });

  eventsDaysEntries.forEach(([day, events]) => day.events = events);

  if (this.updateRows) {
    this.refreshRows();
  }
  if (this.updateColumns) {
    this.refreshColumns();
  }

  return this;
};

/**
 * The same method with `addEvents` except one thing. This method is working asynchronously and returns Promise.
 *
 * @param {Array} events
 * @param {boolean} [allowDuplicates = false]
 * @param {boolean} [delayRefresh = false]
 * @returns {Promise<Calendar>}
 */
Calendar.prototype.addEventsAsync = (
  async function addEventsAsync(events, allowDuplicates = false, delayRefresh = false) {
    for (let i = 0; i < events.length; i += 1) {
      this.addEvent(events[i], allowDuplicates, true);
    }
    if (!delayRefresh) {
      return this.refreshEventsAsync();
    }

    return this;
  }
);

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
