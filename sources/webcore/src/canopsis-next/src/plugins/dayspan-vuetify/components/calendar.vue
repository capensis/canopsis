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
        @mouse-down-event="startMove",
        @mouse-start-resize="startResize",
        @clear-placeholder="clearPlaceholder",
        @handle-added="handleAdded"
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
        @mouse-down-event="startMove",
        @mouse-start-resize="startResize",
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
        @mouse-down-event="startMove",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="mouseDownDay",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-start-resize="startResize",
        @clear-placeholder="clearPlaceholder"
      )
</template>


<script>
import { DsCalendar } from 'dayspan-vuetify/src/components';
import { CalendarEvent, DaySpan, Op, Schedule } from 'dayspan';

export default {
  extends: DsCalendar,
  computed: {
    canResize() {
      return !this.readOnly && !this.$dayspan.readOnly;
    },

    hasCreatePopover() {
      return !!this.$scopedSlots.eventCreatePopover;
    },
  },
  methods: {
    createEventFromCalendar(calendarEvent) {
      return {
        data: calendarEvent.data,
        schedule: calendarEvent.schedule.toInput(),
        id: calendarEvent.event.id,
      };
    },

    copyCalendarEvent(calendarEvent) {
      const details = { ...calendarEvent.data };
      const span = new DaySpan(calendarEvent.start, calendarEvent.end);
      const schedule = calendarEvent.fullDay
        ? Schedule.forDay(
          calendarEvent.start,
          calendarEvent.time.days(Op.UP),
        )
        : Schedule.forSpan(span);
      const event = this.$dayspan.createEvent(details, schedule);
      event.id = calendarEvent.event.id;

      return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
    },

    startMove(mouseEvent) {
      if (this.canMove && mouseEvent.left) {
        this.movingEvent = mouseEvent;
        this.moving = true;
        this.movingDuration = mouseEvent.calendarEvent.time.millis();
        this.placeholderForCreate = false;
        this.placeholder = this.copyCalendarEvent(mouseEvent.calendarEvent);
      }
    },

    startResize(event, calendarEvent) {
      if (this.canResize) {
        this.resizing = true;
        this.resizingEvent = event;
        this.resizingBelow = true;
        this.placeholderForCreate = false;
        this.placeholder = this.copyCalendarEvent(calendarEvent);
      }
    },

    endResize() {
      this.resizing = false;
      this.resizingEvent = null;
      this.resizingBelow = true;
    },

    finishAdd(mouseEvent) {
      if (!this.hasCreatePopover) {
        this.handleAdded(mouseEvent);
      } else {
        this.placeholderForCreate = true;
      }
    },

    handleAdded(mouseEvent) {
      const event = this.getEvent('added', {
        mouseEvent,
        calendarEvent: this.createEventFromCalendar(this.placeholder),
        span: this.placeholder.time,
      });

      this.$emit('added', event);

      this.endAdd();
    },

    finishMove(mouseEvent) {
      const target = this.placeholder.time;
      const source = this.movingEvent.calendarEvent.time;
      const sameTime = target.start.sameMinute(source.start);
      const sameDay = target.start.sameDay(source.start);
      const isDay = mouseEvent.type === 'mouse-up-day';

      if ((isDay && !sameDay) || (!isDay && !sameTime)) {
        const calendarEvent = this.copyCalendarEvent(this.placeholder);

        const event = this.getEvent('moved', {
          mouseEvent,
          calendarEvent: this.createEventFromCalendar(calendarEvent),
        });

        this.$emit('moved', event);
      }

      this.endMove();
    },

    finishResize(mouseEvent) {
      const calendarEvent = this.copyCalendarEvent(this.placeholder);

      const event = this.getEvent('resized', {
        mouseEvent,
        calendarEvent: this.createEventFromCalendar(calendarEvent),
      });

      this.$emit('resized', event);

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

      if (start.time === time.time || end.time === time.time) {
        return;
      }

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

    changeAddDayPlaceholder(mouseEvent) {
      this.addEnd = mouseEvent.day;

      const min = this.addStart.min(this.addEnd);
      const max = this.addStart.max(this.addEnd);

      this.placeholder.day = min.start();
      this.placeholder.time.start = min;
      this.placeholder.time.end = max.end();
      this.placeholder.event.schedule = Schedule.forDay(
        this.placeholder.start,
        this.placeholder.time.days(Op.UP),
      );

      this.updatePlaceholderRow(this.placeholder);
    },

    changeMoveDayPlaceholder(mouseEvent) {
      const { day } = mouseEvent;

      this.placeholder.day = day;
      this.placeholder.time.start = day;
      this.placeholder.time.end = day.next(this.placeholder.schedule.durationInDays).end();

      this.updatePlaceholderRow(this.placeholder);
    },

    changeResizeDayPlaceholder(mouseEvent) {
      const { day } = mouseEvent;
      const { start, end } = this.placeholder.time;

      if (start.time === day.time || end.time === day.time) {
        return;
      }

      if (this.resizingBelow && start.time > day.time) {
        this.resizingBelow = false;
        this.placeholder.time.end = start.end();
      }

      if (!this.resizingBelow && end.time < day.time) {
        this.resizingBelow = true;
        this.placeholder.time.start = end;
      }

      if (this.resizingBelow) {
        this.placeholder.time.end = day.end();
      } else {
        this.placeholder.time.start = day;
      }

      this.updatePlaceholderRow(this.placeholder);
    },

    mouseMove(mouseEvent) {
      if (!mouseEvent.left) {
        return;
      }

      if (this.adding) {
        this.changeAddPlaceholder(mouseEvent);
      }

      if (this.moving) {
        this.changeMovePlaceholder(mouseEvent);
      }

      if (this.resizing) {
        this.changeResizePlaceholder(mouseEvent);
      }
    },

    mouseMoveDay(mouseEvent) {
      if (!mouseEvent.left) {
        return;
      }

      if (this.adding) {
        this.changeAddDayPlaceholder(mouseEvent);
      }

      if (this.moving) {
        this.changeMoveDayPlaceholder(mouseEvent);
      }

      if (this.resizing) {
        this.changeResizeDayPlaceholder(mouseEvent);
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
    },

    clearPlaceholder() {
      this.placeholder = null;
      this.placeholderForCreate = false;
    },
  },
};
</script>
