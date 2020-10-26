<template lang="pug">
  v-menu.ds-calendar-event-placeholder(
    :content-class="contentClass",
    :disabled="menuDisabled",
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
      v-bind="{ placeholder, calendar, day, add, close: closePopover }"
    )
</template>

<script>
import { Functions as fn } from 'dayspan';

import { DsCalendarEventTimePlaceholder } from 'dayspan-vuetify/src/components';

import popoverMixin from '../mixins/popover';

export default {
  extends: DsCalendarEventTimePlaceholder,
  mixins: [popoverMixin],
  methods: {
    add(calendarEvent) {
      this.$emit('add-event', this.getEvent('add-event', { calendarEvent }));
    },

    getEvent(type, extra = {}) {
      return fn.extend({
        type,
        calendarEvent: this.calendarEvent,
        closePopover: this.closePopover,
        openPopover: this.openPopover,
        calendar: this.calendar,
        $vm: this,
        $element: this.$el,
      }, extra);
    },
  },
};
</script>
