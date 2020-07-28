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
    slot(
      name="eventCreatePopover",
      v-if="isStart && isShownPopover",
      v-bind="{ placeholder, calendar, day, add, close }"
    )
</template>

<script>
import { DsCalendarEventTimePlaceholder } from 'dayspan-vuetify/src/components';

import popoverMixin from '../mixins/popover';

export default {
  extends: DsCalendarEventTimePlaceholder,
  mixins: [popoverMixin],
  methods: {
    close() {
      this.$emit('clear-placeholder');
    },

    add(calendarEvent) {
      this.$emit('add-event', calendarEvent);
    },
  },
};
</script>
