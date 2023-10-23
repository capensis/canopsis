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
      :value="focus"
      :events="events"
      :type="type"
      @click:event="$emit('click:event', $event)"
      @input="setCalendarFocus"
      @change="updateRange"
    >
      <template #event="{ event }">
        <v-layout align-center>
          <v-icon
            v-if="event.icon"
            :color="event.iconColor"
            size="14"
            class="pr-1"
          >
            {{ event.icon }}
          </v-icon>
          {{ event.name }}
        </v-layout>
      </template>
    </v-calendar>

    <div class="c-calendar__loader">
      <slot name="loader" />
    </div>
  </v-layout>
</template>

<script>
import { CALENDAR_TYPES } from '@/constants';

import { convertDateToDateObject } from '@/helpers/date/date';

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
    };
  },
  methods: {
    updateRange({ start, end }) {
      this.$emit('change', {
        start: convertDateToDateObject(start.date),
        end: convertDateToDateObject(end.date),
      });
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
  },
};
</script>

<style lang="scss">
.c-calendar {
  position: relative;
}
</style>
