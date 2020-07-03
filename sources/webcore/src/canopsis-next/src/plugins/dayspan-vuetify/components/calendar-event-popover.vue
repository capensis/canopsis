<template lang="pug">
  v-card.ds-calendar-event-popover-card
    slot(
      v-if="calendar",
      v-bind="{ calendarEvent: calendarEventForSlot, calendar, edit, add, close, readOnly }"
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
  computed: {
    calendarEventForSlot() {
      return this.placeholder || this.calendarEvent;
    },
  },
};
</script>
