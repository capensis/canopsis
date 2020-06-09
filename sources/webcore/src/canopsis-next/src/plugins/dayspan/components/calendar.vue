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
import { CalendarEvent, DaySpan, Schedule } from 'dayspan';

export default {
  extends: DsCalendar,
  methods: {
    getPlaceholderEventForResize(calendarEvent) {
      const details = { ...calendarEvent.data };
      const event = this.$dayspan.createEvent(details, calendarEvent.schedule, true);
      const span = new DaySpan(calendarEvent.start, calendarEvent.end);

      return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
    },

    startResize(event, calendarEvent) {
      this.resizing = true;
      this.resizingEvent = event;
      this.resizingBelow = true;
      this.placeholder = this.getPlaceholderEventForResize(calendarEvent);
    },

    endResize() {
      this.resizing = false;
      this.resizingEvent = null;
      this.resizingBelow = true;
      this.placeholder = null;
    },

    finishResize(mouseEvent) {
      const event = this.getEvent('resized', {
        mouseEvent,
        placeholder: this.placeholder,
        span: this.placeholder.time,
      });

      this.$emit('resized', event);

      if (!event.handled) {
        event.clearPlaceholder();
      }

      this.endResize();
    },

    changeAddPlaceholder(mouseEvent) {
      this.addEnd = mouseEvent.timeDrag;

      const min = this.addStart.min(this.addEnd);
      const max = this.addStart.max(this.addEnd);

      this.placeholder.day = min.start();
      this.placeholder.time.start = min;
      this.placeholder.time.end = max;
      this.placeholder.event.schedule = Schedule.forSpan(this.placeholder.time);
    },

    changeMovePlaceholder(mouseEvent) {
      let { time } = mouseEvent;
      time = time.relative(-this.movingEvent.offset);
      time = this.$dayspan.roundTime(time, this.$dayspan.rounding.drag);

      this.placeholder.day = time.start();
      this.placeholder.time.start = time;
      this.placeholder.time.end = time.relative(this.movingDuration);
    },

    changeResizePlaceholder(mouseEvent) {
      let { time } = mouseEvent;
      time = time.relative(-this.resizingEvent.offset);
      time = this.$dayspan.roundTime(time, this.$dayspan.rounding.drag);
      const { start, end } = this.placeholder.time;

      if (this.resizingBelow && start.time > time.time) {
        this.resizingBelow = false;
        this.placeholder.time.end = start;
      }

      if (!this.resizingBelow && end.time < time.time) {
        this.resizingBelow = true;
        this.placeholder.time.start = end;
      }

      if (this.resizingBelow) {
        this.placeholder.time.end = time;
      } else {
        this.placeholder.time.start = time;
      }
    },

    mouseMove(mouseEvent) {
      if (this.adding && mouseEvent.left) {
        this.changeAddPlaceholder(mouseEvent);
      }

      this.mouseMoveCheckReady();
      if (this.moving) {
        this.changeMovePlaceholder(mouseEvent);
      }
      this.mouseMoveCheckEnd(mouseEvent);

      if (this.resizing) {
        this.changeResizePlaceholder(mouseEvent);
      }
    },

    mouseUp(mouseEvent) {
      if (this.addEnd) {
        this.finishAdd(mouseEvent);
      }

      if (this.moving) {
        this.finishMove(mouseEvent);
      }

      if (this.resizing) {
        this.finishResize(mouseEvent);
      }

      this.readyToMove = false;
    },
  },
};
</script>
