<template lang="pug">
  div.ds-calendar
    div.ds-month-view.ds-year-view(v-if="isYear")
      ds-weeks-view.ds-expand(
        v-bind="{ $scopedSlots }",
        v-on="$listeners",
        :calendar="calendar",
        :placeholder="placeholder",
        :placeholder-for-create="placeholderForCreate",
        @edit-event="editHandler",
        @add-event="addHandler",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="startAdd",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-down-event="startMove",
        @mouse-start-resize="startResize",
        @clear-placeholder="clearPlaceholder"
      )

    div.ds-month-view(v-if="isMonth")
      ds-weeks-view.ds-expand(
        v-bind="{ $scopedSlots }",
        v-on="$listeners",
        :calendar="calendar",
        :placeholder="placeholder",
        :placeholder-for-create="placeholderForCreate",
        @edit-event="editHandler",
        @add-event="addHandler",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="startAdd",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-down-event="startMove",
        @mouse-start-resize="startResize",
        @clear-placeholder="clearPlaceholder"
      )

    div.ds-week-view(v-if="isWeek || isDay")
      ds-days-view(
        v-bind="{ $scopedSlots }",
        v-on="$listeners",
        :calendar="calendar",
        :placeholder="placeholder",
        :placeholder-for-create="placeholderForCreate",
        @edit-event="editHandler",
        @add-event="addHandler",
        @mouse-move="mouseMove",
        @mouse-down="mouseDown",
        @mouse-up="mouseUp",
        @mouse-down-event="startMove",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="startAdd",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-start-resize="startResize",
        @clear-placeholder="clearPlaceholder"
      )
</template>


<script>
import { CalendarEvent, DaySpan, Op, Schedule } from 'dayspan';
import { DsCalendar } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendar,
  computed: {
    canResize() {
      return !this.readOnly && !this.$dayspan.readOnly;
    },

    hasCreatePopover() {
      return !!this.$scopedSlots.eventCreatePopover;
    },

    hasEditPopover() {
      return !!this.$scopedSlots.eventPopover;
    },

    openPopover() {
      return (!this.placeholder.event.id && this.hasCreatePopover) || this.hasEditPopover;
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
        ? Schedule.forDay(span.start, span.days(Op.UP))
        : Schedule.forSpan(span);
      const event = this.$dayspan.createEvent(details, schedule, true);
      event.id = calendarEvent.event.id;

      return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
    },

    editHandler(calendarEvent) {
      const event = this.getEvent('changed', {
        calendarEvent: this.createEventFromCalendar(this.copyCalendarEvent(calendarEvent)),
      });

      this.$emit('changed', event);
    },

    addHandler(calendarEvent) {
      const event = this.getEvent('added', {
        calendarEvent: this.createEventFromCalendar(this.copyCalendarEvent(calendarEvent)),
      });

      this.$emit('added', event);
    },

    startAdd(mouseEvent) {
      if (this.placeholder) {
        this.endAdd();
        this.endMove();
        this.endResize();
        return;
      }

      if (this.canAdd && mouseEvent.left) {
        const { day } = mouseEvent;

        this.addStart = day;
        this.addEnd = day;
        this.placeholderForCreate = false;
        this.placeholder = this.$dayspan.getPlaceholderEventForAdd(day);
        this.placeholder.event.schedule = Schedule.forDay(day);
        this.placeholder.fullDay = true;

        this.updatePlaceholderRow();
      }
    },

    startMove(mouseEvent) {
      if (this.canMove && mouseEvent.left) {
        this.readyToMove = true;
        this.movingEvent = mouseEvent;
      }
    },

    mouseMoveCheck() {
      if (this.readyToMove) {
        const { day } = this.movingEvent;
        const { time, schedule } = this.movingEvent.calendarEvent;

        this.moving = true;
        this.movingDuration = time.millis();
        this.movingStartDay = day;
        this.placeholderForCreate = false;
        this.placeholder = this.copyCalendarEvent(this.movingEvent.calendarEvent);
        this.placeholder.time.end = this.placeholder.fullDay
          ? time.start.next(schedule.durationInDays).end()
          : time.start.relative(this.movingDuration);
        this.readyToMove = false;

        this.updatePlaceholderRow();
      }
    },

    startResize(event, calendarEvent) {
      if (this.canResize) {
        this.resizing = true;
        this.resizingEvent = event;
        this.resizingBelow = true;
        this.placeholderForCreate = false;
        this.placeholder = this.copyCalendarEvent(calendarEvent);

        this.updatePlaceholderRow();
      }
    },

    finishAdd(mouseEvent) {
      if (!this.openPopover) {
        this.handleAdded(mouseEvent);
      } else {
        this.placeholderForCreate = true;
      }

      this.endAdd();
    },

    finishMove(mouseEvent) {
      if (!this.openPopover) {
        this.handleMoved(mouseEvent);
      } else {
        this.placeholderForCreate = true;
      }

      if (this.placeholder.start.isBefore(this.calendar.filled.start)) {
        this.calendar.prev();
      }

      this.endMove();
    },

    finishResize(mouseEvent) {
      if (!this.openPopover) {
        this.handleResized(mouseEvent);
      } else {
        this.placeholderForCreate = true;
      }

      this.endResize();
    },

    handleAdded(mouseEvent) {
      const event = this.getEvent('added', {
        mouseEvent,
        calendarEvent: this.createEventFromCalendar(this.placeholder),
        span: this.placeholder.time,
      });

      this.$emit('added', event);
    },

    handleMoved(mouseEvent) {
      const target = this.placeholder.time;
      const source = this.movingEvent.calendarEvent.time;
      const sameTime = target.start.sameMinute(source.start);
      const sameDay = target.start.sameDay(source.start);
      const isDay = mouseEvent.type === 'mouse-up-day';

      if ((isDay && !sameDay) || (!isDay && !sameTime)) {
        const calendarEvent = this.copyCalendarEvent(this.placeholder);

        const event = this.getEvent('changed', {
          mouseEvent,
          calendarEvent: this.createEventFromCalendar(calendarEvent),
        });

        this.$emit('changed', event);
      } else {
        this.clearPlaceholder();
      }
    },

    handleResized(mouseEvent) {
      const calendarEvent = this.copyCalendarEvent(this.placeholder);

      const event = this.getEvent('changed', {
        mouseEvent,
        calendarEvent: this.createEventFromCalendar(calendarEvent),
      });

      this.$emit('changed', event);
    },

    endAdd() {
      this.addStart = null;
      this.addEnd = null;
    },

    endMove() {
      this.readyToMove = false;
      this.moving = false;
      this.movingEvent = null;
    },

    endResize() {
      this.resizing = false;
      this.resizingEvent = null;
      this.resizingBelow = true;
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
      const { time } = mouseEvent;
      const { calendarEvent } = this.movingEvent;
      let { offset = 0 } = this.movingEvent;
      let newTime = time;

      if (this.movingStartDay) {
        if (!calendarEvent.start.sameDay(this.movingStartDay)) {
          offset -= calendarEvent.start.date.diff(this.movingStartDay.date, 'milliseconds');
        }
      }

      newTime = newTime.relative(-offset);
      newTime = this.$dayspan.roundTime(newTime, this.$dayspan.rounding.drag);

      this.placeholder.day = newTime.start();
      this.placeholder.time.start = newTime;
      this.placeholder.time.end = newTime.relative(this.movingDuration);
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
        this.placeholder.time.start,
        this.placeholder.time.days(Op.UP),
      );

      this.updatePlaceholderRow();
    },

    changeMoveDayPlaceholder(mouseEvent) {
      const { calendarEvent } = this.movingEvent;
      const { day } = mouseEvent;
      let newDay = day;

      if (this.movingStartDay) {
        const diff = calendarEvent.start.date.diff(this.movingStartDay.date, 'days');

        newDay = newDay.next(diff);
      }

      this.placeholder.day = newDay;
      this.placeholder.time.start = newDay;
      this.placeholder.time.end = newDay.next(this.placeholder.event.schedule.durationInDays).end();

      this.updatePlaceholderRow();
    },

    changeResizeDayPlaceholder(mouseEvent) {
      const { day } = mouseEvent;
      const { start, end } = this.placeholder.time;

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

      this.placeholder.event.schedule = Schedule.forDay(
        this.placeholder.time.start,
        this.placeholder.time.days(Op.UP),
      );

      this.updatePlaceholderRow();
    },

    updatePlaceholderRow() {
      let row = 0;

      this.calendar.iterateDays().iterate((day) => {
        if (this.placeholder.time.matchesDay(day)) {
          row = day.iterateEvents().reduce(
            row,
            (calendarEvent, maxRow) => Math.max(calendarEvent.row + 1, maxRow),
            calendarEvent => calendarEvent.event.id !== this.placeholder.event.id,
          );
        }
      });

      this.placeholder.row = row;
    },

    mouseMove(mouseEvent) {
      if (!mouseEvent.left) {
        return;
      }
      this.mouseMoveCheck(mouseEvent);

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

    mouseMoveCheckReady() {
      if (this.readyToMove) {
        const { calendarEvent, day } = this.movingEvent;

        const ev = this.getEvent('moving', {
          calendarEvent,
          moveEvent: this.movingEvent,
          placeholder: this.$dayspan.getPlaceholderEventForMove(calendarEvent),
        });

        this.$emit('moving', ev);

        if (!ev.handled && ev.placeholder) {
          this.moving = true;
          this.movingDuration = calendarEvent.time.millis();
          this.movingStartDay = day;
          this.placeholderForCreate = false;
          this.placeholder = ev.placeholder;
        }

        this.readyToMove = false;
      }
    },

    mouseMoveDay(mouseEvent) {
      if (!mouseEvent.left) {
        return;
      }
      this.mouseMoveCheckReady();

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
      if (this.readyToMove) {
        this.endMove();
      }

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
      if (!this.moving && !this.resizing) {
        this.placeholder = null;
        this.placeholderForCreate = false;
      }
    },
  },
};
</script>
