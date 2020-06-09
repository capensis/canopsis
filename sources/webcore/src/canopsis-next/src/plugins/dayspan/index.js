import DaySpanVuetify from 'dayspan-vuetify';
import { CalendarEvent, Schedule } from 'dayspan';

import DsCalendar from './components/calendar.vue';
import DsCalendarApp from './components/calendar-app.vue';
import DsCalendarEvent from './components/calendar-event.vue';
import DsCalendarEventTime from './components/calendar-event-time.vue';


Schedule.prototype.resize = function resize(span) {
  if (this.start) {
    this.start = span.start.start();
  }
  if (this.end) {
    this.end = span.end.end();
  }
};

CalendarEvent.prototype.resize = function resize(span) {
  this.schedule.resize(span);
};

export default {
  install(Vue, options = {}) {
    Vue.use(DaySpanVuetify, options);
    Vue.component('dsCalendar', DsCalendar);
    Vue.component('dsCalendarApp', DsCalendarApp);
    Vue.component('dsCalendarEvent', DsCalendarEvent);
    Vue.component('dsCalendarEventTime', DsCalendarEventTime);
  },
};
