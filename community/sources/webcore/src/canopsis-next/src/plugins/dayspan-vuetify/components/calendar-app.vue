<template lang="pug">
  div.ds-expand.ds-calendar-app
    v-layout.pa-3
      slot(name="today", v-bind="{ setToday, todayDate, calendar }")
        v-flex
          v-tooltip(bottom)
            template(#activator="{ on }")
              v-btn.ds-skinny-button.ds-calendar-app-action(
                v-on="on",
                :icon="$vuetify.breakpoint.smAndDown",
                depressed,
                @click="setToday"
              )
                span(v-if="$vuetify.breakpoint.mdAndUp") {{ labels.today }}
                v-icon(v-else) {{ labels.todayIcon }}
            span {{ todayDate }}
      slot(name="pagination", v-bind="{ prev, prevLabel, next, nextLabel, summary, calendar }")
        v-flex.text-sm-center
          v-tooltip(bottom)
            template(#activator="{ on }")
              v-btn.mx-2.ds-light-forecolor.ds-calendar-app-action(
                v-on="on",
                icon,
                depressed,
                @click="prev"
              )
                v-icon keyboard_arrow_left
            span {{ prevLabel }}
          calendar-app-period-picker(:calendar="calendar", @change="selectPeriod")
          v-tooltip(bottom)
            template(#activator="{ on }")
              v-btn.mx-2.ds-light-forecolor.ds-calendar-app-action(
                v-on="on",
                icon,
                depressed,
                @click="next"
              )
                v-icon keyboard_arrow_right
            span {{ nextLabel }}
      slot(name="view", v-bind="{ currentType, types }")
        v-flex.text-sm-right
          v-menu
            template(#activator="{ on }")
              v-btn.ds-calendar-app-action(v-on="on", flat) {{ currentType.label }}
                v-icon arrow_drop_down
            v-list
              v-list-tile(
                v-for="type in types",
                :key="type.id",
                @click="currentType = type"
              )
                v-list-tile-content
                  v-list-tile-title {{ type.label }}
      slot(name="menuRight")
    div.ds-expand
      v-container.ds-calendar-container(
        :fluid="fluid",
        :fill-height="fillHeight"
      )
        slot(name="calendarAppLoader")
        ds-gestures(
          @swipeleft="next",
          @swiperight="prev"
        )
          div.ds-expand(v-if="currentType.schedule")
            slot(name="calendarAppAgenda", v-bind="{ $scopedSlots, $listeners, calendar, viewDay }")
              ds-agenda(
                v-bind="{ $scopedSlots }",
                v-on="$listeners",
                :calendar="calendar",
                :read-only="readOnly",
                @view-day="viewDay"
              )
          div.ds-expand(v-else)
            slot(name="calendarAppCalendar", v-bind="{ $scopedSlots, $listeners, calendar, viewDay }")
              ds-calendar(
                ref="calendar",
                v-bind="{ $scopedSlots }",
                v-on="$listeners",
                :calendar="calendar",
                :read-only="readOnly",
                :current-time-for-today="currentTimeForToday",
                @view-day="viewDay"
              )

        slot(name="calendarAppEventDialog", v-bind="{ $scopedSlots, $listeners, calendar, eventFinish }")
        slot(name="calendarAppOptions", v-bind="{ optionsVisible, optionsDialog, options, chooseOption }")
        slot(name="calendarAppPrompt", v-bind="{ promptVisible, promptDialog, promptQuestion, choosePrompt }")
        slot(name="containerInside", v-bind="{ events, calendar }")
</template>

<script>
import { Calendar, Sorts } from 'dayspan';

import calendarOptionsMixin from '../mixins/calendar-options';

import CalendarAppPeriodPicker from './calendar-app-period-picker.vue';

export default {
  components: { CalendarAppPeriodPicker },
  mixins: [calendarOptionsMixin],
  props: {
    events: {
      type: Array,
      default: () => [],
    },
    calendar: {
      type: Calendar,
      default: () => Calendar.months(),
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    types: {
      type: Array,
      default() {
        return this.$dsDefaults().types.filter(({ id }) => ['M', 'W', 'D'].includes(id));
      },
    },
    formats: {
      type: Object,
      validate(x) {
        return this.$dsValidate(x, 'formats');
      },
      default() {
        return this.$dsDefaults().formats;
      },
    },
    labels: {
      type: Object,
      validate(x) {
        return this.$dsValidate(x, 'labels');
      },
      default() {
        return this.$dsDefaults().labels;
      },
    },
    optionsDialog: {
      type: Object,
      validate(x) {
        return this.$dsValidate(x, 'optionsDialog');
      },
      default() {
        return this.$dsDefaults().optionsDialog;
      },
    },
    promptDialog: {
      type: Object,
      validate(x) {
        return this.$dsValidate(x, 'promptDialog');
      },
      default() {
        return this.$dsDefaults().promptDialog;
      },
    },
    fluid: {
      type: Boolean,
      default: false,
    },
    fillHeight: {
      type: Boolean,
      default: false,
    },
    currentTimeForToday: {
      type: Boolean,
      default: false,
    },
    removeEventsBeforeMove: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      drawer: null,
      optionsVisible: false,
      options: [],
      promptVisible: false,
      promptQuestion: '',
      promptCallback: null,
    };
  },
  computed: {
    currentType: {
      get() {
        return this.types.find(
          type => type.type === this.calendar.type && type.size === this.calendar.size,
        ) || this.types[0];
      },
      set(type) {
        this.rebuild(undefined, true, type);
      },
    },

    summary() {
      const small = this.$vuetify.breakpoint.xs;

      if (small) {
        return this.calendar.start.format(this.formats.xs);
      }

      const large = this.$vuetify.breakpoint.mdAndUp;

      return this.calendar.summary(false, !large, false, !large);
    },

    todayDate() {
      return this.$dayspan.today.format(this.formats.today);
    },

    nextLabel() {
      return this.labels.next(this.currentType);
    },

    prevLabel() {
      return this.labels.prev(this.currentType);
    },
  },

  watch: {
    calendar: 'applyEvents',
    events: {
      immediate: true,
      handler: 'applyEvents',
    },
    readOnly: {
      immediate: true,
      handler: 'applyReadOnly',
    },
  },
  created() {
    /**
     * We've added that for 'eventSorter' field initialization on the calendar
     */
    this.rebuild(undefined, true, this.currentType, true);
  },
  mounted() {
    if (!this.$dayspan.promptOpen) {
      this.$dayspan.promptOpen = (question, callback) => {
        this.promptVisible = false;
        this.promptQuestion = question;
        this.promptCallback = callback;
        this.promptVisible = true;
      };
    }
  },
  methods: {
    setState(state, ignoreTriggerChange) {
      state.eventSorter = state.listTimes
        ? Sorts.List([Sorts.FullDay, Sorts.Start])
        : Sorts.Start;

      this.calendar.set(state);

      if (!ignoreTriggerChange) {
        this.triggerChange();
      }
    },

    applyEvents() {
      if (this.events) {
        this.calendar.removeEvents();
        this.calendar.addEvents(this.events);
      }
    },

    applyReadOnly(value) {
      this.$dayspan.readOnly = value;
    },

    isType(type, aroundDay) {
      const cal = this.calendar;

      return (
        cal.type === type.type
        && cal.size === type.size
        && (
          !aroundDay
          || cal.span.matchesDay(aroundDay)
        )
      );
    },

    rebuild(aroundDay, force, forceType, ignoreTriggerChange) {
      const type = forceType || this.currentType || this.types[2];

      if (this.isType(type, aroundDay) && !force) {
        return;
      }

      const input = {
        type: type.type,
        size: type.size,
        around: aroundDay,
        eventsOutside: true,
        preferToday: false,
        listTimes: type.listTimes,
        updateRows: type.updateRows,
        updateColumns: type.listTimes,
        fill: !type.listTimes,
        otherwiseFocus: type.focus,
        repeatCovers: type.repeat,
      };

      this.setState(input, ignoreTriggerChange);
    },

    selectPeriod(diff) {
      if (this.removeEventsBeforeMove) {
        this.calendar.removeEvents(null, true);
      }

      this.calendar.move(diff);

      this.triggerChange();
    },

    next() {
      if (this.removeEventsBeforeMove) {
        this.calendar.removeEvents(null, true);
      }

      this.calendar.unselect().next();

      this.triggerChange();
    },

    prev() {
      if (this.removeEventsBeforeMove) {
        this.calendar.removeEvents(null, true);
      }

      this.calendar.unselect().prev();

      this.triggerChange();
    },

    setToday() {
      this.rebuild(this.$dayspan.today);
    },

    viewDay(day) {
      this.rebuild(day, false, this.types[0]);
    },

    chooseOption(option) {
      if (option) {
        option.callback();
      }

      this.optionsVisible = false;
    },

    choosePrompt(yes) {
      this.promptCallback(yes);
      this.promptVisible = false;
    },

    eventFinish() {
      this.triggerChange();
    },

    triggerChange() {
      this.$emit('change', {
        calendar: this.calendar,
      });
    },
  },
};
</script>

<style lang="scss">
  .ds-week-view-container {
    max-height: 75vh;
  }

  .ds-week-view {
    position: relative !important;
    overflow: hidden;
    max-height: 75vh;

    .ds-week-header .ds-day {
      border-top: #e0e0e0 1px solid;
    }
  }

  .ds-day {
    min-height: 10em;
    user-select: none;

    .ds-month & {
      padding-bottom: 22px;
    }
  }

  .ds-day:first-child, .ds-week-header-day:first-child {
    border-left: #e0e0e0 1px solid;
  }

  .ds-week-header-day {
    border-top: #e0e0e0 1px solid;
  }

  .theme--dark {
    .ds-month {
      background: var(--v-secondary) !important;
    }

    .ds-week-header-day,
    .ds-dom,
    .ds-week-date,
    .ds-week-weekday,
    .ds-hour-text {
      color: white !important;
    }
  }
</style>
