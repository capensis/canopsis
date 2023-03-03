<template lang="pug">
  div.ds-calendar(:class="{ editing, adding, moving, resizing }")
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
        @mouse-down-day="startAddDay",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-down-event="startMove",
        @mouse-start-resize="startResize",
        @mouse-start-edit="startEditing",
        @mouse-end-edit="endEditing",
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
        @mouse-down-day="startAddDay",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-down-event="startMove",
        @mouse-start-resize="startResize",
        @mouse-start-edit="startEditing",
        @mouse-end-edit="endEditing",
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
        @mouse-down="startAdd",
        @mouse-up="mouseUp",
        @mouse-down-event="startMove",
        @mouse-move-day="mouseMoveDay",
        @mouse-down-day="startAddDay",
        @mouse-up-day="mouseUp",
        @mouse-up-event="mouseUp",
        @mouse-start-resize="startResize",
        @mouse-start-edit="startEditing",
        @mouse-end-edit="endEditing",
        @clear-placeholder="clearPlaceholder"
      )
</template>

<script>
import { get } from 'lodash';

import { CalendarEvent, Day, DaySpan, Op, Schedule } from 'dayspan';
import { DsCalendar } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendar,
  props: {
    currentTimeForToday: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      placeholder: null,
      placeholderForCreate: false,
      addStart: null,
      editing: false,
      resizing: false,
      moving: false,
    };
  },
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
      const details = {
        ...calendarEvent.data,
        resizing: false,
        moving: false,
      };
      const span = new DaySpan(calendarEvent.start, calendarEvent.end);

      if (calendarEvent.fullDay && !span.end.isEnd()) {
        span.end = span.end.prev().end();
      }

      const schedule = calendarEvent.fullDay
        ? Schedule.forDay(span.start, span.days(Op.UP))
        : Schedule.forSpan(span);
      const event = this.$dayspan.createEvent(details, schedule, true);
      event.id = calendarEvent.event.id;

      return new CalendarEvent(calendarEvent.id, event, span, calendarEvent.day);
    },

    editHandler({ calendarEvent, closePopover }) {
      const event = this.getEvent('changed', {
        calendarEvent: this.copyCalendarEvent(calendarEvent),
        closePopover,
      });

      this.$emit('changed', event);
    },

    addHandler({ calendarEvent, closePopover }) {
      const event = this.getEvent('added', {
        calendarEvent: this.copyCalendarEvent(calendarEvent),
        closePopover,
      });

      this.$emit('added', event);
    },

    startAddDay(mouseEvent) {
      if (this.placeholder || this.editing) {
        this.endAdd();
        this.endMove();
        this.endResize();
        return;
      }

      if (this.canAdd && mouseEvent.left) {
        this.endEditing();
        const { day } = mouseEvent;

        this.placeholderForCreate = false;

        if (this.currentTimeForToday && day.currentDay) {
          const time = Day.now();
          this.addStart = time;
          this.addEnd = Day.today().end();
          this.placeholder = this.$dayspan.getPlaceholderEventForAdd(time);
          this.placeholder.event.schedule = Schedule.forTime(time, time.asTime());
          this.placeholder.fullDay = false;
        } else {
          this.addStart = day;
          this.addEnd = day;
          this.placeholder = this.$dayspan.getPlaceholderEventForAdd(day);
          this.placeholder.event.schedule = Schedule.forDay(day);
          this.placeholder.fullDay = true;
        }

        this.updatePlaceholderRow();
      }
    },

    startAdd(mouseEvent) {
      if (this.placeholder || this.editing) {
        this.endAdd();
        this.endMove();
        this.endResize();
        return;
      }

      if (this.canAdd && mouseEvent.left) {
        this.endEditing();
        const { time } = mouseEvent;

        this.addStart = time;
        this.placeholderForCreate = false;
        this.placeholder = this.$dayspan.getPlaceholderEventForAdd(time);
        this.placeholder.event.schedule = Schedule.forTime(time, time.asTime());
        this.placeholder.fullDay = false;
      }
    },

    startMove(mouseEvent) {
      if (this.canMove && mouseEvent.left) {
        this.readyToMove = true;
        this.movingEvent = mouseEvent;
      }
    },

    startEditing(event) {
      if (this.editingEvent && this.editingEvent.closePopover) {
        this.editingEvent.closePopover(event);
      }

      this.editing = true;
      this.editingEvent = event;
    },

    mouseMoveCheck() {
      if (this.readyToMove) {
        this.endEditing();

        const { day } = this.movingEvent;
        const { time, schedule } = this.movingEvent.calendarEvent;

        this.moving = true;
        this.movingDuration = time.millis();
        this.movingStartDay = day;
        this.placeholderForCreate = false;
        this.placeholder = this.copyCalendarEvent(this.movingEvent.calendarEvent);
        this.placeholder.data.moving = true;
        this.placeholder.time.end = this.placeholder.fullDay
          ? time.start.next(schedule.durationInDays).end()
          : time.start.relative(this.movingDuration);
        this.readyToMove = false;

        this.updatePlaceholderRow();
      }
    },

    startResize(event) {
      const { calendarEvent } = event;

      if (this.canResize) {
        const { time, schedule } = calendarEvent;

        this.resizing = true;
        this.resizingEvent = event;
        this.resizingBelow = true;
        this.placeholderForCreate = false;
        this.placeholder = this.copyCalendarEvent(calendarEvent);
        this.placeholder.data.resizing = true;
        this.placeholder.time.end = this.placeholder.fullDay
          ? time.start.next(schedule.durationInDays).end()
          : time.end;

        this.updatePlaceholderRow();
        this.endEditing();
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
      this.placeholder.data.moving = false;

      this.handleMoved(mouseEvent);
      this.endMove();
    },

    finishResize(mouseEvent) {
      this.placeholder.data.resizing = false;

      this.handleResized(mouseEvent);
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
      const isDay = mouseEvent.type === 'mouse-up-day';

      if (
        isDay
          ? !target.start.sameDay(source.start)
          : !target.start.sameMinute(source.start)
      ) {
        const event = this.getEvent('moved', {
          mouseEvent,
          movingEvent: this.movingEvent,
          calendarEvent: this.copyCalendarEvent(this.movingEvent.calendarEvent),
          target: this.placeholder.time,
          openPopover: () => this.placeholderForCreate = true,
          closePopover: () => this.clearPlaceholder(),
        });

        this.$emit('moved', event);

        if (!event.handled) {
          event.clearPlaceholder();
        }
      } else {
        this.$nextTick(() => this.clearPlaceholder());
      }
    },

    handleResized(mouseEvent) {
      const target = this.placeholder.time;
      const source = this.resizingEvent.calendarEvent.time;

      if (!(target.start.sameMinute(source.start) && target.end.sameMinute(source.end))) {
        const event = this.getEvent('resized', {
          mouseEvent,
          resizingEvent: this.resizingEvent,
          calendarEvent: this.copyCalendarEvent(this.resizingEvent.calendarEvent),
          target: this.placeholder.time,
          openPopover: () => this.placeholderForCreate = true,
          closePopover: () => this.clearPlaceholder(),
        });

        this.$emit('resized', event);

        if (!event.handled) {
          event.clearPlaceholder();
        }
      } else {
        this.$nextTick(() => this.clearPlaceholder());
      }
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

    endEditing(event) {
      if (
        this.editingEvent
        && get(event, 'calendarEvent.event.id') === get(this.editingEvent, 'calendarEvent.event.id')
      ) {
        this.editing = false;
        this.editingEvent = null;
      }
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
      const { start, end } = this.placeholder.time;
      let { time } = mouseEvent;

      time = this.$dayspan.roundTime(time, this.$dayspan.rounding.drag, true);

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

      this.placeholder.day = this.placeholder.time.start;
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

      this.placeholder.day = this.placeholder.time.start;

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

        this.moving = true;
        this.movingDuration = calendarEvent.time.millis();
        this.movingStartDay = day;
        this.placeholderForCreate = false;
        this.placeholder = this.copyCalendarEvent(calendarEvent);
        this.placeholder.data.moving = true;

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

      if (this.addEnd || this.adding) {
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

<style lang="scss" scoped>
.ds-calendar {
  &.resizing, &.adding {
    & ::v-deep .ds-week-view .ds-day {
      &, * {
        cursor: ns-resize !important;
      }
    }

    & ::v-deep .ds-month-view .ds-day {
      &, * {
        cursor: ew-resize !important;
      }
    }
  }

  &.moving {
    & ::v-deep .ds-day {
      &, * {
        cursor: grabbing !important;
      }
    }
  }

  & ::v-deep .v-menu__activator {
    &--disabled {
      cursor: inherit !important;

      & * {
        cursor: inherit !important;
      }
    }
  }

  & ::v-deep .ds-month-view .ds-dom {
    display: inline-block;
  }
}
</style>
