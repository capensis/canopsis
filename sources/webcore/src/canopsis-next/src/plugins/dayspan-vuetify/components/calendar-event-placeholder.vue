<template lang="pug">
  v-menu.ds-calendar-event-menu(
    :content-class="contentClass",
    :disabled="!hasPopover",
    v-model="menu",
    v-bind="popoverProps"
  )
    ds-calendar-event(
      slot="activator",
      v-bind="{ $scopedSlots }",
      v-on="$listeners",
      :is-placeholder-with-day="day",
      :calendar-event="placeholder",
      :calendar="calendar",
      :index="index"
    )
    slot(
      name="eventCreatePopover",
      v-if="isStart && menu",
      v-bind="{ placeholder, calendar, day, index, add, close }"
    )
</template>

<script>
import { DsCalendarEventPlaceholder } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEventPlaceholder,
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
