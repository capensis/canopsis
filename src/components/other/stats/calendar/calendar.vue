<template lang="pug">
  v-container.ds-expand.ds-calendar-app
    v-layout(row)
      v-flex(xs3)
        v-btn.ds-skinny-button(
        depressed,
        @click="setToday"
        )
          span {{ $t('calendar.today') }}
      v-flex(xs6)
        v-btn.ds-light-forecolor.mx-2(icon, depressed, @click="prev")
          v-icon keyboard_arrow_left
        span.subheading {{ summary }}
        v-btn.ds-light-forecolor.mx-2(icon, depressed, @click="next")
          v-icon keyboard_arrow_right
      v-flex(xs3)
        v-menu
          v-btn(slot="activator", flat)
            span {{ currentType.label }}
            v-icon arrow_drop_down
          v-list
            v-list-tile(
            v-for="type in types"
            :key="type.id"
            @click="currentType = type"
            )
              v-list-tile-content
                v-list-tile-title {{ type.label }}
    v-container.ds-calendar-container(fluid, fill-height)
      ds-gestures(@swipeleft="next", @swiperight="prev")
        .ds-expand
          ds-calendar(
          ref="calendar",
          v-bind="{$scopedSlots}",
          v-on="$listeners",
          :calendar="calendar",
          @edit="edit",
          @view-day="viewDay",
          )
</template>

<script>
import moment from 'moment';
import { rrulestr } from 'rrule';
import { Units, Sorts, Calendar, Day, Schedule } from 'dayspan';
import dsDefaults from 'dayspan-vuetify/src/defaults';

const c = Calendar.months();

const pbehavior = {
  _id: 'asd',
  author: 'User',
  comment: null,
  connector: 'canopsis',
  connector_name: 'canopsis',
  eids: ['something'],
  enabled: true,
  filter: '{"infos.scenario_name.value": "Sc_gspfo_aude_p"}',
  name: 'downtime',
  reason: '',
  rrule: 'FREQ=WEEKLY;COUNT=30;WKST=FR;BYDAY=MO,TU,WE;BYHOUR=15;BYMINUTE=12;BYSECOND=11',
  tstart: 1537947053000,
  tstop: 1539652000000,
  type_: 'Hors plage horaire de surveillance',
};

if (pbehavior.rrule) {
  const ruleSecond = rrulestr(pbehavior.rrule, {
    dtstart: new Date(pbehavior.tstart),
  });

  if (!ruleSecond.options.until) {
    ruleSecond.options.until = new Date(pbehavior.tstop);
  }

  const days = ruleSecond.all();

  const daysObjects = days.map(day => new Day(moment(day)));

  c.addEvents(daysObjects.map(dayObject => (
    {
      data: {
        title: 'PBEHAVIOR',
        description: 'Something',
        color: '#3F51B5',
        meta: pbehavior,
      },
      schedule: new Schedule({
        on: dayObject,
        times: [dayObject.asTime()],
      }),
    }
  )));
}

export default {
  props: {
    events: {
      type: Array,
    },
    calendar: {
      type: Calendar,
      default() {
        return Calendar.months();
      },
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

    todayDate() {
      return this.$dayspan.today.format(this.formats.today);
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

    edit() {
      // console.log(calendarEvent);
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
  .ds-day {
    min-height: 10em;
  }
</style>
