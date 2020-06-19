<template lang="pug">
  v-card.pa-2.ds-calendar-event-popover-card
    slot(v-if="calendar", v-bind="slotData")
</template>

<script>
import { Calendar, CalendarEvent } from 'dayspan';

export default {
  props: {
    placeholder: {
      required: true,
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
    close: {
      type: Function,
    },
  },
  computed: {
    slotData() {
      return {
        placeholder: this.placeholder,
        calendar: this.calendar,
        edit: this.edit,
        close: this.close,
        details: this.details,
        readOnly: this.readOnly,
      };
    },
  },
};
</script>
