import DaySpanVuetify from 'dayspan-vuetify/src/plugin';

import DsCalendar from './components/calendar.vue';
import DsCalendarApp from './components/calendar-app.vue';
import DsCalendarEvent from './components/calendar-event.vue';
import DsCalendarEventTime from './components/calendar-event-time.vue';

export default {
  install(Vue, options = {}) {
    Vue.use(DaySpanVuetify, options);
    Vue.component('dsCalendar', DsCalendar);
    Vue.component('dsCalendarApp', DsCalendarApp);
    Vue.component('dsCalendarEvent', DsCalendarEvent);
    Vue.component('dsCalendarEventTime', DsCalendarEventTime);
  },
};
