<template lang="pug">
  v-card.pa-2.ds-calendar-event-popover-card
    slot(
      v-if="calendar",
      v-bind="{ placeholder: calendarEventForm, calendar, edit: editHandler, add: addHandler, close, readOnly }"
    )
</template>

<script>
import { Calendar, CalendarEvent } from 'dayspan';

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
  data() {
    return {
      calendarEventForm: this.eventToForm(this.placeholder || this.calendarEvent),
    };
  },
  methods: {
    addHandler() {
      this.add(this.formToEvent(this.calendarEventForm));
    },

    editHandler() {
      this.edit(this.formToEvent(this.calendarEventForm));
    },

    eventToForm(calendarEvent) {
      return {
        meta: calendarEvent.meta,
        start: calendarEvent.start,
        end: calendarEvent.end,
        ...calendarEvent.data,
      };
    },

    formToEvent() {
      // TODO prepare form to calendar event logic should be added
      return this.placeholder || this.calendarEvent;
    },
  },
};
</script>
