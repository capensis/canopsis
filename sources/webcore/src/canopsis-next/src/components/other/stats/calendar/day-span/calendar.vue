<template lang="pug">
  .ds-expand.ds-calendar-app
    v-layout.pa-3(row)
      v-flex(md3)
        v-btn.ds-skinny-button(
          depressed,
          @click="setToday"
        )
          span {{ $t('calendar.today') }}
      v-flex.text-sm-center(md6)
        v-btn.ds-light-forecolor.mx-2(icon, depressed, @click="prev")
          v-icon keyboard_arrow_left
        span.subheading {{ summary }}
        v-btn.ds-light-forecolor.mx-2(icon, depressed, @click="next")
          v-icon keyboard_arrow_right
      v-flex.text-sm-right(md3)
        v-menu
          v-btn(slot="activator", flat)
            span {{ currentType.label }}
            v-icon arrow_drop_down
          v-list
            v-list-tile(
              v-for="type in types",
              :key="type.id",
              @click="currentType = type"
            )
              v-list-tile-content
                v-list-tile-title {{ type.label }}
    v-container.ds-calendar-container(fluid, fill-height)
      ds-gestures(@swipeleft="next", @swiperight="prev")
        ds-calendar(
          ref="calendar",
          v-on="$listeners",
          v-bind="{$scopedSlots}",
          :calendar="calendar",
          @view-day="viewDay"
        )
</template>

<script>
import { Units, Sorts, Calendar } from 'dayspan';
import dsDefaults from 'dayspan-vuetify/src/defaults';

/**
 * Dayspan Vuetify calendar wrapper
 *
 * @see https://github.com/ClickerMonkey/dayspan-vuetify/tree/master/docs
 *
 * @prop {Array} events
 * @prop {Object} widget
 */
export default {
  props: {
    events: {
      type: Array,
      default: () => [],
    },
    calendar: {
      type: Calendar,
      default: () => Calendar.months(),
    },
    formats: {
      validate(x) {
        return this.$dsValidate(x, 'formats');
      },
      default() {
        return dsDefaults.dsCalendarApp.formats;
      },
    },
    labels: {
      validate(x) {
        return this.$dsValidate(x, 'labels');
      },
      default() {
        return dsDefaults.dsCalendarApp.labels;
      },
    },
  },
  computed: {
    types() {
      const defaultTypeValues = {
        size: 1,
        focus: 0.4999,
        repeat: true,
        listTimes: true,
        updateRows: true,
        schedule: false,
      };

      return [
        {
          ...defaultTypeValues,

          id: 'M',
          label: this.$t('calendar.month'),
          shortcut: 'M',
          type: Units.MONTH,
          listTimes: false,
        },
        {
          ...defaultTypeValues,

          id: 'W',
          label: this.$t('calendar.week'),
          shortcut: 'W',
          type: Units.WEEK,
        },
        {
          ...defaultTypeValues,

          id: 'D',
          label: this.$t('calendar.day'),
          shortcut: 'D',
          type: Units.DAY,
        },
      ];
    },
    currentType: {
      get() {
        return this.types.find(type =>
          type.type === this.calendar.type &&
          type.size === this.calendar.size) || this.types[0];
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
  },

  watch: {
    events: 'applyEvents',
    calendar: 'applyEvents',
  },
  created() {
    this.rebuild(undefined, true);
  },
  methods: {
    setState(state) {
      state.eventSorter = state.listTimes
        ? Sorts.List([Sorts.FullDay, Sorts.Start])
        : Sorts.Start;

      this.calendar.set(state);

      this.triggerChange();
    },

    applyEvents() {
      if (this.events) {
        this.calendar.removeEvents();
        this.calendar.addEvents(this.events);
      }
    },

    isType(type, aroundDay) {
      const cal = this.calendar;

      return (cal.type === type.type && cal.size === type.size &&
        (!aroundDay || cal.span.matchesDay(aroundDay)));
    },

    rebuild(aroundDay, force, forceType) {
      const type = forceType || this.currentType || this.types[0];

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

      this.setState(input);
    },

    next() {
      this.calendar.unselect().next();

      this.triggerChange();
    },

    prev() {
      this.calendar.unselect().prev();

      this.triggerChange();
    },

    setToday() {
      this.rebuild(this.$dayspan.today);
    },

    viewDay(day) {
      this.rebuild(day, false, this.types[2]);
    },

    triggerChange() {
      this.$emit('change', {
        calendar: this.calendar,
      });
    },
  },
};
</script>

<style>
  .ds-calendar {
    min-height: 100vh;
  }

  .ds-day {
    min-height: 10em;
  }
</style>
