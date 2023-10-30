<template>
  <v-layout
    class="c-calendar"
    fill-height
    column
  >
    <v-layout
      justify-space-around
      class="c-calendar__toolbar py-3"
    >
      <calendar-today-btn @click="setToday" />
      <v-layout justify-center>
        <calendar-pagination
          :focus.sync="focus"
          :type="type"
          @prev="prev"
          @next="next"
        />
      </v-layout>
      <div>
        <calendar-view-mode v-model="type" />
        <slot name="menu-right" />
      </div>
    </v-layout>

    <v-calendar
      ref="calendar"
      v-bind="$attrs"
      v-on="calendarListeners"
      :class="calendarClasses"
      :value="focus"
      :events="availableEvents"
      :type="type"
      :event-color="getEventColor"
      :event-height="calendarEventHeight"
    >
      <template #event="{ event, start, end }">
        <v-layout
          class="pl-1"
          align-center
        >
          <v-icon
            v-if="event.icon"
            :color="event.iconColor"
            size="14"
            class="pr-1"
          >
            {{ event.icon }}
          </v-icon>
          <span v-if="start">{{ event.name }}</span>
          <div
            v-if="end"
            :class="['c-calendar__event-drag-bottom', { 'c-calendar__event-drag-bottom--right': isMonthType }]"
            @mousedown.stop="startResize(event)"
          />
        </v-layout>
      </template>
    </v-calendar>

    <div class="c-calendar__loader">
      <slot name="loader" />
    </div>

    <v-menu
      :value="popoverOpen"
      :close-on-content-click="false"
      :close-on-click="false"
      :position-x="positionX"
      :position-y="positionY"
      content-class="c-calendar__popover-wrapper"
    >
      <v-card v-if="popoverOpen">
        <v-card-text>
          <slot
            name="form-event"
            :close="closeCreateEventPopover"
            :event="newEvent || editEvent || popoverEvent"
          />
        </v-card-text>
      </v-card>
    </v-menu>
  </v-layout>
</template>

<script>
import { CALENDAR_TYPES } from '@/constants';

import { convertDateToEndOfDayDateObject, convertDateToStartOfDayDateObject } from '@/helpers/date/date';
import { colorToRgba } from '@/helpers/color';

import CalendarTodayBtn from './partials/calendar-today-btn.vue';
import CalendarViewMode from './partials/calendar-view-mode.vue';
import CalendarPagination from './partials/calendar-pagination.vue';

export default {
  components: { CalendarPagination, CalendarViewMode, CalendarTodayBtn },
  inheritAttrs: false,
  props: {
    events: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      type: CALENDAR_TYPES.month,
      focus: new Date(),
      filled: {
        start: null,
        stop: null,
      },

      popoverOpen: false,
      dragging: false,
      resizing: false,
      creating: false,

      dragStartTime: 0,

      newEvent: null,
      popoverEvent: null,
      editEvent: null,

      movingEnd: true,

      positionY: null,
      positionX: null,
    };
  },
  computed: {
    isMonthType() {
      return this.type === CALENDAR_TYPES.month;
    },

    availableEvents() {
      if (this.editEvent) {
        return this.events.map(event => (event.id === this.editEvent.id ? this.editEvent : event));
      }

      return this.newEvent
        ? [...this.events, this.newEvent]
        : this.events;
    },

    calendarEventHeight() {
      return {
        [CALENDAR_TYPES.day]: 0,
        [CALENDAR_TYPES.week]: 0,
        [CALENDAR_TYPES.month]: 20,
      }[this.type];
    },

    calendarClasses() {
      return ['c-calendar__calendar', {
        'c-calendar__calendar--dragging': this.dragging,
        'c-calendar__calendar--creating': this.creating,
        'c-calendar__calendar--resizing-month': this.resizing && this.isMonthType,
        'c-calendar__calendar--resizing-day': this.resizing && !this.isMonthType,
      }];
    },

    calendarListeners() {
      const listeners = {
        input: this.setCalendarFocus,
        change: this.updateRange,

        'click:event': this.showEventDetails,

        'mousedown:time': this.startCreateEvent,
        'mousedown:day': this.startCreateEvent,
        'mousedown:event': this.startDragEvent,
      };

      if (this.creating) {
        listeners['mouseup:time'] = this.finishCreateEvent;
        listeners['mouseup:day'] = this.finishCreateEvent;

        listeners['mousemove:day'] = this.handleMouseMoveEvent;
        listeners['mousemove:time'] = this.handleMouseMoveEvent;
      }

      if (this.dragging) {
        listeners['mouseup:time'] = this.finishDragEvent;
        listeners['mouseup:day'] = this.finishDragEvent;

        listeners['mousemove:day'] = this.handleMouseMoveEventDrag;
        listeners['mousemove:time'] = this.handleMouseMoveEventDrag;
      }

      if (this.resizing) {
        listeners['mouseup:time'] = this.finishResizeEvent;
        listeners['mouseup:day'] = this.finishResizeEvent;

        listeners['mousemove:day'] = this.handleMouseMoveEvent;
        listeners['mousemove:time'] = this.handleMouseMoveEvent;
      }

      return listeners;
    },
  },
  mounted() {
    this.$refs.calendar.checkChange();
  },
  methods: {
    getEventColor(event) {
      if (event.id === this.editEvent?.id) {
        return colorToRgba(event.color, 0.5);
      }

      return event.color;
    },

    updateRange({ start, end }) {
      this.filled = {
        start: convertDateToStartOfDayDateObject(start.date),
        end: convertDateToEndOfDayDateObject(end.date),
      };

      this.$emit('change');
    },

    setFocusDate(date) {
      this.focus = date;
    },

    setToday() {
      this.setFocusDate(new Date());
    },

    setCalendarFocus(date) {
      this.setFocusDate(new Date(date));
    },

    prev() {
      this.$refs.calendar.prev();
    },

    next() {
      this.$refs.calendar.next();
    },

    showCreateEventPopover(event, target) {
      const { top, left, width } = target.getBoundingClientRect();
      this.popoverEvent = event;

      this.positionX = left + width / 2;
      this.positionY = top;

      this.popoverOpen = true;
    },

    closeCreateEventPopover() {
      this.popoverOpen = false;
      this.popoverEvent = null;
      this.newEvent = null;
      this.editEvent = null;
    },

    showEventDetails({ nativeEvent, event }) {
      if (this.$listeners['click:event']) {
        this.$emit('click:event');
        return;
      }

      this.showCreateEventPopover(event, nativeEvent.target);

      nativeEvent.stopPropagation();
    },

    getDateByEvent(event) {
      const date = new Date(event.year, event.month - 1, event.day);

      if (event.time) {
        date.setHours(event.hour);
        date.setMinutes(event.minute);
      }

      return date;
    },

    getStartDate(date, timed) {
      return timed ? new Date(date) : convertDateToStartOfDayDateObject(date);
    },

    getEndDate(date, timed) {
      return timed ? new Date(date) : convertDateToEndOfDayDateObject(date);
    },

    handleMouseMoveEvent(event) {
      const timed = !!event.time;

      const eventDate = this.getDateByEvent(event);
      const newEnd = this.getEndDate(eventDate, timed);
      const newStart = this.getStartDate(eventDate, timed);

      const currentEvent = this.creating ? this.newEvent : this.editEvent;

      if (this.movingEnd && newEnd.getTime() < currentEvent.start.getTime()) {
        this.movingEnd = false;
        currentEvent.end = this.getEndDate(currentEvent.start, timed);
      } else if (!this.movingEnd && newStart.getTime() > currentEvent.end.getTime()) {
        this.movingEnd = true;
        currentEvent.start = this.getStartDate(currentEvent.end, timed);
      }

      if (this.movingEnd) {
        currentEvent.end = newEnd;
      } else {
        currentEvent.start = newStart;
      }
    },

    startCreateEvent(event) {
      const timed = !!event.time;

      const eventDate = this.getDateByEvent(event);
      const start = this.getStartDate(eventDate, timed);
      const end = this.getEndDate(eventDate, timed);

      if (event.time) {
        end.setMinutes(start.getMinutes() + 15);
      }

      this.newEvent = {
        name: this.$t('calendar.noTitle'),
        start,
        end,
        timed,
        data: {},
      };
      this.creating = true;
      this.movingEnd = true;
    },

    finishCreateEvent(event, nativeEvent) {
      this.showCreateEventPopover(this.newEvent, nativeEvent.target);

      this.creating = false;
    },

    startDragEvent({ event }, nativeEvent) {
      nativeEvent.stopPropagation();

      if (this.editEvent) {
        this.closeCreateEventPopover();
      }

      this.editEvent = { ...event };
      this.dragging = true;
    },

    handleMouseMoveEventDrag(event) {
      const mouseTime = this.getDateByEvent(event).getTime();

      if (!this.dragStartTime) {
        this.dragStartTime = mouseTime;
        return;
      }

      const diff = mouseTime - this.dragStartTime;

      if (diff) {
        this.dragStartTime = mouseTime;

        const start = new Date(this.editEvent.start.getTime() + diff);
        const end = new Date(this.editEvent.end.getTime() + diff);

        this.editEvent.start = start;
        this.editEvent.end = end;
      }
    },

    finishDragEvent(event, nativeEvent) {
      nativeEvent.stopPropagation();

      this.showCreateEventPopover(this.newEvent, nativeEvent.target);

      this.dragStartTime = 0;
      this.dragging = false;
    },

    startResize(event) {
      this.editEvent = { ...event };
      this.resizing = true;
    },

    finishResizeEvent(event, nativeEvent) {
      nativeEvent.stopPropagation();

      this.showCreateEventPopover(this.newEvent, nativeEvent.target);

      this.resizing = false;
    },
  },
};
</script>

<style lang="scss">
.c-calendar {
  position: relative;

  &__calendar {
    .v-calendar-weekly__week {
      min-height: 140px;
    }

    &--dragging {
      cursor: grab;
    }

    &--resizing-month {
      cursor: ew-resize;
    }

    &--resizing-day {
      cursor: ns-resize;
    }

    &--dragging, &--creating, &--resizing-day, &--resizing-month {
      .v-event, .v-event-timed {
        pointer-events: none;
        user-select: none;
      }
    }
  }

  &__event-drag-bottom {
    position: absolute;
    cursor: ns-resize;
    opacity: 0;
    transition: .3s cubic-bezier(.25, .8, .5,1);

    &:hover {
      opacity: 1;
    }

    &:not(&--right) {
      cursor: ns-resize;
      left: 0;
      right: 0;
      bottom: 4px;
      height: 4px;
      border-top: 1px solid white;
      border-bottom: 1px solid white;
    }

    &--right {
      cursor: ew-resize;
      right: 4px;
      top: 0;
      bottom: 0;
      width: 4px;
      border-left: 1px solid white;
      border-right: 1px solid white;
    }
  }

  &__popover-wrapper {
    max-height: 95%;
    max-width: 95% !important;
    width: 980px !important;
    top: 50% !important;
    transform: translate3d(0, -50%, 0);
  }
}
</style>
