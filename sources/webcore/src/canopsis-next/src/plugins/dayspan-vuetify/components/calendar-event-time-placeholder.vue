<template lang="pug">
  v-menu.ds-calendar-event-placeholder(
    :content-class="contentClass",
    :disabled="!hasPopover",
    v-model="menu",
    v-bind="popoverProps"
  )
    ds-calendar-event-time(
      slot="activator",
      v-bind="{ $scopedSlots }",
      v-on="$listeners",
      :is-placeholder-with-day="day",
      :calendar-event="placeholder",
      :calendar="calendar"
    )
    slot(v-if="isStart", name="eventCreatePopover", v-bind="{ placeholder, calendar, day, add, close }")
</template>

<script>
import { DsCalendarEventTimePlaceholder } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEventTimePlaceholder,
  methods: {
    close() {
      this.$emit('clear-placeholder');
      this.menu = false;
    },

    add(calendarEvent) {
      this.$emit('add-event', calendarEvent);
      this.menu = false;
    },

    triggerClearPlaceholder() { },
  },
};
</script>
