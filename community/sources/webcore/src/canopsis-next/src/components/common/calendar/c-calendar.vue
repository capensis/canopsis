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
          :type="internalType"
          @prev="prev"
          @next="next"
        />
      </v-layout>
      <div>
        <calendar-view-mode v-model="internalType" />
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
      :type="internalType"
      :weekdays="weekdays"
      :event-height="calendarEventHeight"
      :event-color="getEventColor"
      :weekday-format="formatWeekday"
    >
      <template #day-label="event">
        <v-btn
          :color="event.present ? 'primary' : 'transparent'"
          fab
          depressed
          small
          @click.stop.prevent="moveToDay(event)"
          @mousedown.stop=""
        >
          {{ event.day }}
        </v-btn>
        <template v-if="event.day === 1">
          {{ formatMonth(event) }}
        </template>
      </template>

      <template #day-label-header="event">
        <v-btn
          :color="event.present ? 'primary' : 'transparent'"
          fab
          depressed
          small
          @click.stop.prevent="moveToDay(event)"
          @mousedown.stop=""
        >
          {{ event.day }}
        </v-btn>
      </template>

      <template #event="{ event, start, end }">
        <slot
          name="event"
          :event="event"
          :start="start"
          :end="end"
        >
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
              :class="['c-calendar__event-drag-bottom', { 'c-calendar__event-drag-bottom--right': !event.timed }]"
              @mousedown.stop="startResize(event)"
            />
          </v-layout>
        </slot>
      </template>
    </v-calendar>

    <div class="c-calendar__loader">
      <slot name="loader">
        <c-progress-overlay :pending="loading" />
      </slot>
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
            :close="clearPlaceholder"
            :event="newEvent || editEvent || popoverEvent"
          />
        </v-card-text>
      </v-card>
    </v-menu>
  </v-layout>
</template>

<script>
import { CALENDAR_TYPES, CALENDAR_START_DRAG_DELAY } from '@/constants';
import { LOCALES } from '@/config';

import {
  convertDateToEndOfDayDateObject,
  convertDateToStartOfDayDateObject,
  getWeekdayNumber,
} from '@/helpers/date/date';
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
    loading: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      default: CALENDAR_TYPES.month,
    },
  },
  data() {
    return {
      internalType: this.type,
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
    weekdays() {
      return this.$i18n.locale === LOCALES.en
        ? [0, 1, 2, 3, 4, 5, 6]
        : [1, 2, 3, 4, 5, 6, 0];
    },

    currentEditingEvent() {
      return this.editEvent || this.newEvent;
    },

    isMonthType() {
      return this.internalType === CALENDAR_TYPES.month;
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
      }[this.internalType];
    },

    calendarClasses() {
      const classes = ['c-calendar__calendar', {
        'c-calendar__calendar--dragging': this.dragging,
        'c-calendar__calendar--creating': this.creating,
      }];

      if (this.creating || this.resizing) {
        classes.push(
          this.editEvent?.timed || this.newEvent?.timed
            ? 'c-calendar__calendar--resizing-bottom'
            : 'c-calendar__calendar--resizing-right',
        );
      }

      return classes;
    },

    calendarListeners() {
      const listeners = {
        input: this.setCalendarFocus,
        change: this.updateRange,

        'click:event': this.showEventDetails,
      };

      if (this.readonly) {
        return listeners;
      }

      if (!this.popoverOpen) {
        listeners['mousedown:time'] = this.startCreateEvent;
        listeners['mousedown:day'] = this.startCreateEvent;
      }

      listeners['mousedown:event'] = this.startDragEvent;

      listeners['mouseup:time'] = this.handleMouseUpTime;
      listeners['mouseup:day'] = this.handleMouseUpTime;

      if (this.creating || this.resizing) {
        listeners['mousemove:day'] = this.handleMouseMoveEvent;
        listeners['mousemove:time'] = this.handleMouseMoveEvent;
      } else if (this.dragging) {
        listeners['mousemove:day'] = this.handleMouseMoveEventDrag;
        listeners['mousemove:time'] = this.handleMouseMoveEventDrag;
      }

      return listeners;
    },
  },
  watch: {
    type(type) {
      if (type !== this.internalType) {
        this.internalType = type;
      }
    },

    internalType() {
      this.$emit('update:type', this.internalType);
    },
  },
  mounted() {
    this.$refs.calendar.checkChange();
  },
  methods: {
    formatWeekday(event) {
      const weekday = getWeekdayNumber(event.date);

      return [
        this.$t('common.shortWeekDays.monday'),
        this.$t('common.shortWeekDays.tuesday'),
        this.$t('common.shortWeekDays.wednesday'),
        this.$t('common.shortWeekDays.thursday'),
        this.$t('common.shortWeekDays.friday'),
        this.$t('common.shortWeekDays.saturday'),
        this.$t('common.shortWeekDays.sunday'),
      ][weekday - 1];
    },

    formatMonth(event) {
      const { month } = event;

      return [
        this.$t('common.shortMonths.january'),
        this.$t('common.shortMonths.february'),
        this.$t('common.shortMonths.march'),
        this.$t('common.shortMonths.april'),
        this.$t('common.shortMonths.may'),
        this.$t('common.shortMonths.june'),
        this.$t('common.shortMonths.july'),
        this.$t('common.shortMonths.august'),
        this.$t('common.shortMonths.september'),
        this.$t('common.shortMonths.october'),
        this.$t('common.shortMonths.november'),
        this.$t('common.shortMonths.december'),
      ][month - 1];
    },

    getEventColor(event) {
      if (event.id === this.currentEditingEvent?.id) {
        return colorToRgba(event.color, 0.5);
      }

      return event.color;
    },

    updateRange({ start, end }) {
      this.filled = {
        start: convertDateToStartOfDayDateObject(start.date),
        end: convertDateToEndOfDayDateObject(end.date),
      };

      this.$emit('change:pagination');
    },

    setFocusDate(date) {
      this.focus = date;
    },

    setToday() {
      this.setFocusDate(new Date());
    },

    setEditEvent(event) {
      this.editEvent = { ...event, oldStart: event.start, oldEnd: event.end };
    },

    resetEditEvent() {
      this.editEvent = null;
    },

    setNewEvent(event) {
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
    },

    resetNewEvent() {
      this.newEvent = null;
    },

    setPopoverEvent(event) {
      this.popoverEvent = event;
    },

    resetPopoverEvent() {
      this.popoverEvent = null;
    },

    closePopover() {
      this.popoverOpen = false;
    },

    setCalendarFocus(date) {
      this.setFocusDate(new Date(date));
    },

    moveToDay(event) {
      this.setCalendarFocus(event.date);
      this.internalType = CALENDAR_TYPES.day;
    },

    prev() {
      this.$refs.calendar.prev();
    },

    next() {
      this.$refs.calendar.next();
    },

    showCreateEventPopover(event, target) {
      const { top, left, width } = target.getBoundingClientRect();

      this.setPopoverEvent(event);

      this.positionX = left + width / 2;
      this.positionY = top;

      this.popoverOpen = true;
    },

    clearPlaceholder() {
      this.closePopover();
      this.resetPopoverEvent();
      this.resetNewEvent();
      this.resetEditEvent();
    },

    showEventDetails(event) {
      if (this.$listeners['click:event']) {
        this.$emit('click:event', event);
        return;
      }

      this.showCreateEventPopover(event.event, event.nativeEvent.target);

      event.nativeEvent.stopPropagation();
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

      if (this.movingEnd && newEnd.getTime() < this.currentEditingEvent.start.getTime()) {
        this.movingEnd = false;
        this.currentEditingEvent.end = this.getEndDate(this.currentEditingEvent.start, timed);
      } else if (!this.movingEnd && newStart.getTime() > this.currentEditingEvent.end.getTime()) {
        this.movingEnd = true;
        this.currentEditingEvent.start = this.getStartDate(this.currentEditingEvent.end, timed);
      }

      if (this.movingEnd) {
        this.currentEditingEvent.end = newEnd;
      } else {
        this.currentEditingEvent.start = newStart;
      }
    },

    startCreateEvent(event) {
      this.setNewEvent(event);

      this.creating = true;
      this.movingEnd = true;
    },

    finishCreateEvent(event, nativeEvent) {
      this.showCreateEventPopover(this.newEvent, nativeEvent.target);

      this.creating = false;
      this.movingEnd = true;
    },

    startDragEvent({ event }, nativeEvent) {
      nativeEvent.stopPropagation();

      if (this.currentEditingEvent) {
        if (this.currentEditingEvent.id === event.id) {
          this.closePopover();

          this.dragging = true;
          return;
        }

        this.clearPlaceholder();
      }

      this.startDraggingTimerId = setTimeout(() => {
        this.startDraggingTimerId = null;

        this.setEditEvent(event);
        this.dragging = true;
      }, CALENDAR_START_DRAG_DELAY);
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

        const start = new Date(this.currentEditingEvent.start.getTime() + diff);
        const end = new Date(this.currentEditingEvent.end.getTime() + diff);

        this.currentEditingEvent.start = start;
        this.currentEditingEvent.end = end;
      }
    },

    isEditingEventChanged() {
      return this.currentEditingEvent.start.getTime() !== this.currentEditingEvent.oldStart.getTime()
        || this.currentEditingEvent.end.getTime() !== this.currentEditingEvent.oldEnd.getTime();
    },

    finishDragEvent(event, nativeEvent) {
      nativeEvent.stopPropagation();

      if (this.newEvent) {
        this.showCreateEventPopover(this.newEvent, nativeEvent.target);
      } else if (this.isEditingEventChanged()) {
        this.$emit('move:event', this.editEvent);
      } else {
        this.clearPlaceholder();
      }

      this.dragStartTime = 0;
      this.dragging = false;
    },

    startResize(event) {
      if (this.newEvent) {
        if (this.newEvent.id === event.id) {
          this.closePopover();
        }
      } else {
        this.setEditEvent(event);
      }

      this.resizing = true;
    },

    finishResizeEvent(event, nativeEvent) {
      if (this.newEvent) {
        this.showCreateEventPopover(this.newEvent, nativeEvent.target);
      } else if (this.isEditingEventChanged()) {
        this.$emit('resize:event', this.editEvent);
      } else {
        this.clearPlaceholder();
      }

      this.resizing = false;
    },

    handleMouseUpTime(event, nativeEvent) {
      if (this.startDraggingTimerId) {
        clearTimeout(this.startDraggingTimerId);
        return;
      }

      if (this.creating) {
        this.finishCreateEvent(event, nativeEvent);
        return;
      }

      if (this.dragging) {
        this.finishDragEvent(event, nativeEvent);
        return;
      }

      if (this.resizing) {
        this.finishResizeEvent(event, nativeEvent);
      }
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

    &--resizing-right {
      cursor: ew-resize;
    }

    &--resizing-bottom {
      cursor: ns-resize;
    }

    &--dragging, &--creating, &--resizing-bottom, &--resizing-right {
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

  .v-calendar-weekly__head-weekday {
    text-align: unset;
  }

  .v-calendar-weekly__head-weekday {
    text-transform: unset;
  }

  .v-calendar-daily_head-weekday {
    text-align: unset;
    padding: 0 12px;
  }

  .v-calendar-daily_head-day-label {
    text-align: unset;
    padding: 0 5px;
  }

  .v-calendar-weekly__day-label {
    text-align: unset;

    .v-btn {
      width: 24px;
      height: 24px;
    }
  }

  .v-calendar-daily_head-day {
    .v-event {
      min-height: 20px;
    }
  }
}
</style>
