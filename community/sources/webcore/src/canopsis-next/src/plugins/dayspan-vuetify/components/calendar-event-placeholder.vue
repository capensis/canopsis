<template lang="pug">
  v-menu(
    :content-class="contentClass",
    :disabled="menuDisabled",
    v-model="menu",
    v-bind="popoverProps"
  )
    template(#activator="{ on }")
      div.ds-calendar-event-menu(v-on="on")
        ds-calendar-event(
          v-bind="{ $scopedSlots }",
          v-on="$listeners",
          :is-placeholder-with-day="day",
          :calendar-event="placeholder",
          :calendar="calendar",
          :index="index"
        )
    slot(
      name="eventCreatePopover",
      v-if="isStart && isShownPopover",
      v-bind="{ placeholder, calendar, day, index, add, close: closePopover }"
    )
</template>

<script>
import { Functions as fn } from 'dayspan';

import { DsCalendarEventPlaceholder } from 'dayspan-vuetify/src/components';

import popoverMixin from '../mixins/popover';

export default {
  extends: DsCalendarEventPlaceholder,
  mixins: [popoverMixin],
  computed: {
    /**
     * We've refactored isStart for resize event. If we will decrease size of event we will not see popover
     *
     * @returns {boolean | string}
     */
    isStart() {
      const { day, start, end } = this.placeholder;
      const placeholderDay = day.isAfter(end) || day.isBefore(start) ? start : day;

      return placeholderDay.sameDay(this.day);
    },
  },
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
