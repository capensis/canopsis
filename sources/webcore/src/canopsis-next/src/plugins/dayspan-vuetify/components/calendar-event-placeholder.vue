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
    slot(v-if="isStart", name="eventCreatePopover", v-bind="slotData")
</template>

<script>
import { DsCalendarEventPlaceholder } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEventPlaceholder,
  computed: {
    slotData() {
      return {
        placeholder: this.placeholder,
        calendar: this.calendar,
        day: this.day,
        index: this.index,
        close: this.close,
      };
    },
  },
  methods: {
    triggerClearPlaceholder() { },
  },
};
</script>
