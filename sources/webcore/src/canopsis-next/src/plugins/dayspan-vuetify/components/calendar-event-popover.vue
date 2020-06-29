<template lang="pug">
  v-card.ds-calendar-event-popover-card
    slot(
      v-if="calendar",
      v-bind="{ placeholder: calendarEventForm, calendar, edit: editHandler, add: addHandler, close, readOnly }"
    )
</template>

<script>
import moment from 'moment';
import { Calendar, CalendarEvent, DaySpan, Op, Schedule, Day } from 'dayspan';

export default {
  props: {
    placeholder: {
      type: CalendarEvent,
    },
    calendarEvent: {
      type: CalendarEvent,
    },
    calendar: {
      required: true,
      type: Calendar,
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
    edit: {
      type: Function,
    },
    add: {
      type: Function,
    },
    close: {
      type: Function,
    },
  },
  computed: {
    calendarEventForm() {
      const calendarEvent = this.placeholder || this.calendarEvent;

      return {
        start: calendarEvent.start.toDate(),
        end: calendarEvent.end.toDate(),
        ...calendarEvent.data,
      };
    },
  },
  methods: {
    addHandler(form) {
      this.add(this.formToEvent(form));
    },

    editHandler(form) {
      this.edit(this.formToEvent(form));
    },

    formToEvent(form) {
      const { end, start, ...details } = form;
      const calendarEvent = this.placeholder || this.calendarEvent;
      const span = new DaySpan(
        Day.fromMoment(moment(start)),
        Day.fromMoment(moment(end)),
      );
      const schedule = calendarEvent.fullDay
        ? Schedule.forDay(span.start, span.days(Op.UP))
        : Schedule.forSpan(span);
      const event = this.$dayspan.createEvent(details, schedule, true);
      event.id = calendarEvent.event.id;

      return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
    },
  },
};
</script>
