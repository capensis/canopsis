<template>
  <c-layout class="c-calendar" column>
    <c-layout class="c-calendar__toolbar py-3" justify-space-around>
      <c-layout class="c-calendar__toolbar-column" gap="2" align-center>
        <calendar-today-btn class="c-calendar__today-btn" @click="setToday" />

        <c-timezone-field
          v-if="!noTimezone"
          v-model="localTimezone"
          class="c-timezone-field"
          server
          hide-details
          @update:model-value="updateTimezone"
        />
      </c-layout>

      <c-layout class="c-calendar__toolbar-column" justify-center>
        <calendar-pagination
          :focus="focus"
          :type="internalType"
          class="c-calendar__pagination"
          @prev="prev"
          @next="next"
          @update:focus="setFocusDate"
        />
      </c-layout>

      <c-layout class="c-calendar__toolbar-column c-calendar__menu-right" justify-end>
        <calendar-view-mode v-model="internalType" />

        <slot name="menu-right" />
      </c-layout>
    </c-layout>

    <div class="position-relative">
      <v-calendar
        ref="calendar"
        v-bind="$attrs"
        :class="calendarClasses"
        :model-value="focus"
        :events="availableEvents"
        :type="internalType"
        :weekdays="weekdays"
        :event-color="getEventColor"
        :weekday-format="formatWeekday"
        :interval-format="formatInterval"
        :event-margin-bottom="3"
        v-on="calendarListeners"
        @update:model-value="setCalendarFocus"
        @change="updateRange"
      >
        <template #day-body="day">
          <div
            v-if="isCurrentTimeLineVisible && day.present"
            :style="{ top: currentTimeYPixels }"
            class="c-calendar__current-time"
          />
        </template>

        <template #day-label="day">
          <v-btn
            :color="day.present ? 'primary' : 'transparent'"
            variant="flat"
            size="x-small"
            icon
            @click.stop.prevent="moveToDay(day)"
            @mousedown.stop
          >
            {{ day.day }}
          </v-btn>

          <template v-if="day.day === 1">
            {{ formatMonth(day) }}
          </template>
        </template>

        <template #day-label-header="day">
          <v-btn
            :color="day.present ? 'primary' : 'transparent'"
            class="c-calendar__week-day-label-btn"
            variant="flat"
            size="x-small"
            icon
            @click.stop.prevent="moveToDay(day)"
            @mousedown.stop
          >
            {{ day.day }}
          </v-btn>
        </template>

        <template #event="{ event, start, end }">
          <slot
            :event="event"
            :start="start"
            :end="end"
            name="event"
          >
            <c-layout
              :class="['c-calendar__event-content pl-1', getMenuClassByCalendarEvent(event.id)]"
              :style="getEventStyle(event)"
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

              <span :style="{ color: event.iconColor }">
                {{ event.name }}
              </span>

              <div
                v-if="!readonly"
                :class="[
                  'c-calendar__event-drag-bottom',
                  { 'c-calendar__event-drag-bottom--right': isMonthType || !event.timed },
                ]"
                @mousedown.stop="startResize(event)"
                @mouseup.stop="finishResizeEvent(event, $event)"
              />
            </c-layout>
          </slot>
        </template>
      </v-calendar>

      <div class="c-calendar__loader">
        <slot name="loader">
          <c-progress-overlay :pending="loading" />
        </slot>
      </div>
    </div>

    <v-menu
      v-model="popoverOpen"
      :close-on-content-click="false"
      :close-on-back="false"
      :target="[positionX, positionY]"
      attach=".pbehavior-calendar-event-menu"
      content-class="c-calendar__popover-wrapper"
      persistent
    >
      <v-card v-if="popoverOpen" class="c-calendar__popover-card">
        <slot
          :close="clearPlaceholder"
          :event="newEvent || editEvent || popoverEvent"
          :timezone="localTimezone"
          :no-timezone="noTimezone"
          name="form-event"
        />
      </v-card>
    </v-menu>
  </c-layout>
</template>

<script>
import { LOCALES } from '@/config';
import {
  CALENDAR_POPOVER_WIDTH,
  CALENDAR_START_DRAG_DELAY,
  CALENDAR_TYPES,
  LLM_AI_CHAT_WIDTH,
  TIME_UNITS,
} from '@/constants';

import { getMenuClassByCalendarEvent } from '@/helpers/calendar/calendar';
import { colorToRgba, getColorByHash, isValidColor } from '@/helpers/color';
import {
  convertDateToMoment,
  convertDateToEndOfDayDateObject,
  convertDateToStartOfDayDateObject,
  getLocalTimezone,
  getWeekdayNumber,
  isDateBefore,
} from '@/helpers/date/date';
import { uid } from '@/helpers/uid';

import CalendarPagination from './partials/calendar-pagination.vue';
import CalendarTodayBtn from './partials/calendar-today-btn.vue';
import CalendarViewMode from './partials/calendar-view-mode.vue';

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
    hideDetailsMenu: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      default: CALENDAR_TYPES.month,
    },
    timezone: {
      type: String,
      default: getLocalTimezone(),
    },
    noTimezone: {
      type: Boolean,
      default: false,
    },
  },

  emits: [
    'update:type',
    'update:timezone',
    'change:pagination',
    'click:event',
    'move:event',
    'resize:event',
  ],

  data() {
    return {
      ready: false,
      internalType: this.type,
      focus: new Date(),
      localTimezone: this.timezone,
      filled: {
        start: null,
        end: null,
      },
      currentTimeIntervalId: null,
      currentTimeYPixels: '-10px',

      popoverOpen: false,
      dragging: false,
      resizing: false,
      creating: false,

      dragStartTime: 0,
      startDraggingTimerId: null,

      newEvent: null,
      popoverEvent: null,
      editEvent: null,

      movingEnd: true,

      positionY: null,
      positionX: null,
    };
  },

  computed: {
    calendarElement() {
      return this.ready ? this.$refs.calendar : null;
    },

    weekdays() {
      return this.$i18n.locale === LOCALES.en ? [0, 1, 2, 3, 4, 5, 6] : [1, 2, 3, 4, 5, 6, 0];
    },

    currentEditingEvent() {
      return this.editEvent || this.newEvent;
    },

    isMonthType() {
      return this.internalType === CALENDAR_TYPES.month;
    },

    availableEvents() {
      if (this.editEvent) {
        return this.events.map(event =>
          event.id === this.editEvent.id ? this.editEvent : event,
        );
      }

      return this.newEvent ? [...this.events, this.newEvent] : this.events;
    },

    calendarClasses() {
      const classes = [
        'c-calendar__calendar',
        {
          'c-calendar__calendar--dragging': this.dragging,
          'c-calendar__calendar--creating': this.creating,
        },
      ];

      if (this.creating || this.resizing) {
        classes.push(
          this.isMonthType || !this.currentEditingEvent?.timed
            ? 'c-calendar__calendar--resizing-right'
            : 'c-calendar__calendar--resizing-bottom',
        );
      }

      return classes;
    },

    calendarListeners() {
      const listeners = {
        'click:event': this.handleClickEvent,
      };

      if (this.readonly) {
        return listeners;
      }

      if (!this.popoverOpen) {
        listeners['mousedown:time'] = this.handleMouseDownTime;
        listeners['mousedown:day'] = this.handleMouseDownDay;
      }

      listeners['mousedown:event'] = this.handleMouseDownEvent;

      listeners['mouseup:time'] = this.handleMouseUpTimeFromCalendar;
      listeners['mouseup:day'] = this.handleMouseUpDayFromCalendar;

      if (this.creating || this.resizing) {
        listeners['mousemove:day'] = this.handleMouseMoveDay;
        listeners['mousemove:time'] = this.handleMouseMoveTime;
      } else if (this.dragging) {
        listeners['mousemove:day'] = this.handleMouseMoveDayDrag;
        listeners['mousemove:time'] = this.handleMouseMoveTimeDrag;
      }

      return listeners;
    },

    isCurrentTimeLineVisible() {
      return !this.isMonthType;
    },
  },

  watch: {
    type(type) {
      if (type !== this.internalType) {
        this.internalType = type;
      }
    },

    timezone(timezone) {
      if (timezone !== this.localTimezone) {
        this.localTimezone = timezone;
      }
    },

    internalType() {
      this.$emit('update:type', this.internalType);
      this.$nextTick(this.updateCurrentTime);
    },
  },

  mounted() {
    this.ready = true;
    this.$refs.calendar.checkChange();
    this.startUpdateCurrentTime();
    this.updateCurrentTime();
  },

  beforeUnmount() {
    this.stopUpdateCurrentTime();
  },

  methods: {
    getMenuClassByCalendarEvent,

    updateCurrentTime() {
      this.$refs.calendar.updateTimes();

      const size = this.calendarElement?.timeToY?.(this.calendarElement.times.now);

      if (size) {
        this.currentTimeYPixels = `${size}px`;
      }
    },

    startUpdateCurrentTime() {
      this.currentTimeIntervalId = setInterval(this.updateCurrentTime, 1000);
    },

    stopUpdateCurrentTime() {
      clearInterval(this.currentTimeIntervalId);
    },

    formatWeekday(day) {
      const weekday = getWeekdayNumber(day.date);

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

    formatMonth(day) {
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
      ][day.month - 1];
    },

    formatInterval(interval) {
      return interval.time;
    },

    getEventSource(event) {
      return event?.input ?? event;
    },

    getEventColor(event) {
      const source = this.getEventSource(event);
      const color = isValidColor(source.color)
        ? source.color
        : getColorByHash(source.id);

      if (source.id === this.currentEditingEvent?.id) {
        return colorToRgba(color, 0.75);
      }

      const past = isDateBefore(source.start, new Date());

      if (!past) {
        return colorToRgba(color, 0.5);
      }

      return color;
    },

    getEventStyle(event) {
      return {
        backgroundColor: this.getEventColor(event),
      };
    },

    updateRange({ start, end }) {
      const parsedStart = this.$refs.calendar.getStartOfWeek(start);
      const parsedEnd = this.$refs.calendar.getEndOfWeek(end);
      const newFilled = {
        start: convertDateToStartOfDayDateObject(parsedStart.date),
        end: convertDateToEndOfDayDateObject(parsedEnd.date),
      };

      if (
        this.filled.start?.getTime() !== newFilled.start.getTime() ||
        this.filled.end?.getTime() !== newFilled.end.getTime()
      ) {
        this.filled = newFilled;

        this.$emit('change:pagination');
      }
    },

    setFocusDate(date) {
      this.focus = date instanceof Date ? date : new Date(date);
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
      const timed = !!event.hasTime;
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
        color: getColorByHash(uid()),
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
      this.setFocusDate(convertDateToStartOfDayDateObject(date));
    },

    moveToDay(day) {
      this.setCalendarFocus(day.date);
      this.internalType = CALENDAR_TYPES.day;
    },

    getPaginationUnit() {
      return {
        [CALENDAR_TYPES.day]: TIME_UNITS.day,
        [CALENDAR_TYPES.week]: TIME_UNITS.week,
        [CALENDAR_TYPES.month]: TIME_UNITS.month,
      }[this.internalType];
    },

    shiftFocus(value) {
      this.setFocusDate(convertDateToMoment(this.focus).add(value, this.getPaginationUnit()).toDate());
    },

    prev() {
      this.shiftFocus(-1);
    },

    next() {
      this.shiftFocus(1);
    },

    showCreateEventPopover(event, target) {
      const { top } = target.getBoundingClientRect();

      this.setPopoverEvent(event);

      const wrapperWidth = window.innerWidth * 0.95;
      const wrapperWidthWithoutAiChat = wrapperWidth - LLM_AI_CHAT_WIDTH;
      const left = (wrapperWidthWithoutAiChat - CALENDAR_POPOVER_WIDTH) / 2;

      this.positionX = left;
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
      event.nativeEvent.stopPropagation();

      if (this.$attrs['onClick:event']) {
        this.$emit('click:event', event);
        return;
      }

      if (!this.hideDetailsMenu) {
        this.showCreateEventPopover(
          event.event,
          event.nativeEvent.currentTarget || event.nativeEvent.target,
        );
      }
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
    },

    finishCreateEvent(event, nativeEvent) {
      this.showCreateEventPopover(this.newEvent, nativeEvent.currentTarget || nativeEvent.target);

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
      return (
        this.currentEditingEvent.start.getTime() !== this.currentEditingEvent.oldStart.getTime() ||
        this.currentEditingEvent.end.getTime() !== this.currentEditingEvent.oldEnd.getTime()
      );
    },

    finishDragEvent(event, nativeEvent) {
      nativeEvent.stopPropagation();

      if (this.newEvent) {
        this.showCreateEventPopover(this.newEvent, nativeEvent.currentTarget || nativeEvent.target);
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
        this.showCreateEventPopover(this.newEvent, nativeEvent.currentTarget || nativeEvent.target);
      } else if (this.isEditingEventChanged()) {
        this.$emit('resize:event', this.editEvent);
      } else {
        this.clearPlaceholder();
      }

      this.resizing = false;
      this.movingEnd = true;
    },

    handleMouseUpTime(event, nativeEvent) {
      if (this.startDraggingTimerId) {
        clearTimeout(this.startDraggingTimerId);
        this.startDraggingTimerId = null;
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

    updateTimezone(newTimezone) {
      this.$emit('update:timezone', newTimezone);
    },

    handleClickEvent(nativeEvent, event) {
      this.showEventDetails({ ...event, nativeEvent });
    },

    handleMouseDownTime(nativeEvent, event) {
      return this.startCreateEvent(event);
    },

    handleMouseDownDay(nativeEvent, event) {
      return this.startCreateEvent(event);
    },

    handleMouseDownEvent(nativeEvent, event) {
      return this.startDragEvent(event, nativeEvent);
    },

    handleMouseUpTimeFromCalendar(nativeEvent, event) {
      return this.handleMouseUpTime(event, nativeEvent);
    },

    handleMouseUpDayFromCalendar(nativeEvent, event) {
      return this.handleMouseUpTime(event, nativeEvent);
    },

    handleMouseMoveDay(nativeEvent, event) {
      return this.handleMouseMoveEvent(event);
    },

    handleMouseMoveTime(nativeEvent, event) {
      return this.handleMouseMoveEvent(event);
    },

    handleMouseMoveDayDrag(nativeEvent, event) {
      return this.handleMouseMoveEventDrag(event);
    },

    handleMouseMoveTimeDrag(nativeEvent, event) {
      return this.handleMouseMoveEventDrag(event);
    },
  },
};
</script>

<style lang="scss">
:root {
  --c-calendar-popover-width: 980px;
}

.c-calendar {
  position: relative;

  .theme--light & {
    --c-calendar-present-day-background-color: rgba(0, 0, 0, 0.04);
  }

  .theme--dark & {
    --c-calendar-present-day-background-color: rgba(255, 255, 255, 0.04);
  }

  &__toolbar-column {
    flex: 0 0 33.3333%;
    max-width: 33.3333%;
  }

  .v-event {
    font-weight: 700;

    box-shadow:
      0 2px 1px -1px rgba(0, 0, 0, 0.2),
      0px 1px 1px 0px rgba(0, 0, 0, 0.14),
      0px 1px 3px 0px rgba(0, 0, 0, 0.12) !important;
  }

  &__event-content {
    width: 100%;
    height: 100%;
    border-radius: inherit;
    overflow: hidden;
  }

  .v-calendar-weekly__head-weekday {
    text-align: unset;
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
    &.v-present {
      background-color: var(--c-calendar-present-day-background-color);
    }
  }

  .v-calendar-daily__day {
    &.v-present {
      background-color: var(--c-calendar-present-day-background-color);
    }
  }

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

    &--dragging,
    &--creating,
    &--resizing-bottom,
    &--resizing-right {
      .v-event,
      .v-event-timed {
        pointer-events: none;
        user-select: none;
      }
    }
  }

  &__event-drag-bottom {
    position: absolute;
    cursor: ns-resize;
    opacity: 0;
    transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);

    &:hover {
      opacity: 1;
    }

    &:not(&--right) {
      cursor: ns-resize;
      left: 0;
      right: 0;
      bottom: -3px;
      height: 6px;
    }

    &--right {
      cursor: ew-resize;
      right: -3px;
      top: 0;
      bottom: 0;
      width: 6px;
    }
  }

  &__popover-wrapper {
    max-height: 95%;
    max-width: 95% !important;
    width: var(--c-calendar-popover-width) !important;
    top: unset !important;
    bottom: unset !important;
    contain: none !important;
    position: relative !important;
    display: flex;
    flex-direction: column;
    overflow: visible !important;
    border-radius: 4px;
    box-shadow:
      0 5px 5px -3px rgba(0, 0, 0, 0.2),
      0 8px 10px 1px rgba(0, 0, 0, 0.14),
      0 3px 14px 2px rgba(0, 0, 0, 0.12);
  }

  &__popover-card {
    display: flex;
    flex-direction: column;
    flex: 1 1 auto;
    min-height: 0;
    max-height: 100%;
    overflow: hidden;
  }

  &__week-day-label-btn {
    font-size: 25px !important;
  }

  .v-present > &__current-time {
    height: 3px;
    background-color: rgb(var(--v-theme-secondary));
    position: absolute;
    left: -1px;
    right: 0;
    pointer-events: none;
  }

  .v-calendar-monthly,
  .v-calendar-weekly__week {
    overflow: hidden;
  }
}
</style>
