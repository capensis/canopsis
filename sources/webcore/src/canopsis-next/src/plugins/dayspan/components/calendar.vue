<template lang="pug">
  .ds-calendar
    .ds-month-view.ds-year-view(v-if="isYear")
      ds-weeks-view.ds-expand(
        v-bind="{ $scopedSlots }",
        v-on="$listeners",
        :calendar="calendar",
        :placeholder="placeholder",
        :placeholder-for-create="placeholderForCreate",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="mouseDownDay",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-down-event="mouseDownEvent",
        @clear-placeholder="clearPlaceholder"
      )

    .ds-month-view(v-if="isMonth")
      ds-weeks-view.ds-expand(
        v-bind="{ $scopedSlots }",
        v-on="$listeners",
        :calendar="calendar",
        :placeholder="placeholder",
        :placeholder-for-create="placeholderForCreate",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="mouseDownDay",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-down-event="mouseDownEvent",
        @clear-placeholder="clearPlaceholder"
      )

    .ds-week-view(v-if="isWeek || isDay")
      ds-days-view(
        v-bind="{ $scopedSlots }",
        v-on="$listeners",
        :calendar="calendar",
        :placeholder="placeholder",
        :placeholder-for-create="placeholderForCreate",
        @mouse-move="mouseMove",
        @mouse-down="mouseDown",
        @mouse-up="mouseUp",
        @mouse-down-event="mouseDownEvent",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="mouseDownDay",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-start-resize="startResize",
        @mouse-end-resize="endResize",
        @clear-placeholder="clearPlaceholder"
      )
</template>


<script>
import { DsCalendar } from 'dayspan-vuetify/src/components';
import { Schedule } from 'dayspan';

export default {
  extends: DsCalendar,
  methods: {
    startResize(event, calendarEvent) {
      this.resizing = true;
      this.resizingEvent = event;
      this.resizingCalendarEvent = calendarEvent;
      this.placeholder = calendarEvent;
    },

    endResize() {
      this.resizing = false;
      this.resizingEvent = null;
      this.resizingCalendarEvent = null;
      this.placeholder = null;
    },

    mouseMove(mouseEvent) {
      if (this.adding && mouseEvent.left) {
        this.addEnd = mouseEvent.timeDrag;

        const min = this.addStart.min(this.addEnd);
        const max = this.addStart.max(this.addEnd);

        this.placeholder.day = min.start();
        this.placeholder.time.start = min;
        this.placeholder.time.end = max;
        this.placeholder.event.schedule = Schedule.forSpan(this.placeholder.time);
      }

      this.mouseMoveCheckReady();

      if (this.moving) {
        let { time } = mouseEvent;
        time = time.relative(-this.movingEvent.offset);
        time = this.$dayspan.roundTime(time, this.$dayspan.rounding.drag);

        this.placeholder.day = time.start();
        this.placeholder.time.start = time;
        this.placeholder.time.end = time.relative(this.movingDuration);
      }

      if (this.resizing) {
        let { time } = mouseEvent;
        time = time.relative(-this.resizingEvent.offset);
        time = this.$dayspan.roundTime(time, this.$dayspan.rounding.drag);

        this.placeholder.time.end = time.relative(this.movingDuration);
      }

      this.mouseMoveCheckEnd(mouseEvent);
    },
  },
};
</script>
