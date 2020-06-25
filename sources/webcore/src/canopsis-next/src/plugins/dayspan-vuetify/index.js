import DaySpanVuetify from 'dayspan-vuetify/src/plugin';

import DsDay from './components/day.vue';
import DsCalendar from './components/calendar.vue';
import DsCalendarApp from './components/calendar-app.vue';
import DsCalendarEvent from './components/calendar-event.vue';
import DsCalendarEventPlaceholder from './components/calendar-event-placeholder.vue';
import DsCalendarEventTimePlaceholder from './components/calendar-event-time-placeholder.vue';
import CalendarEventPopover from './components/calendar-event-popover.vue';
import DsCalendarEventTime from './components/calendar-event-time.vue';

export default {
  install(Vue, options = {}) {
    Vue.use(DaySpanVuetify, options);
    Vue.component('dsDay', DsDay);
    Vue.component('dsCalendar', DsCalendar);
    Vue.component('dsCalendarApp', DsCalendarApp);
    Vue.component('dsCalendarEvent', DsCalendarEvent);
    Vue.component('dsCalendarEventPlaceholder', DsCalendarEventPlaceholder);
    Vue.component('dsCalendarEventTimePlaceholder', DsCalendarEventTimePlaceholder);
    Vue.component('dsCalendarEventPopover', CalendarEventPopover);
    Vue.component('dsCalendarEventTime', DsCalendarEventTime);
  },
};
