<template>
  <div class="ds-expand ds-calendar-app">
    <div>
      <v-btn slot="activator"
             class="ds-skinny-button"
             depressed
             :icon="$vuetify.breakpoint.smAndDown"
             @click="setToday">

        <span>{{ labels.today }}</span>
      </v-btn>

      <v-btn slot="activator"
             icon depressed class="ds-light-forecolor ds-skinny-button"
             @click="prev">
        <v-icon>keyboard_arrow_left</v-icon>
      </v-btn>
      <v-btn slot="activator"
             icon depressed
             class="ds-light-forecolor ds-skinny-button"
             @click="next">
        <v-icon>keyboard_arrow_right</v-icon>
      </v-btn>

      <h1 class="title ds-light-forecolor">
        {{ summary }}
      </h1>

      <v-spacer></v-spacer>

      <v-menu>
        <v-btn flat slot="activator">
          {{ currentType.label }}
          <v-icon>arrow_drop_down</v-icon>
        </v-btn>
        <v-list>
          <v-list-tile v-for="type in types"
                       :key="type.id"
                       @click="currentType = type">
            <v-list-tile-content>
              <v-list-tile-title>{{ type.label }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list>
      </v-menu>
    </div>
    <v-container fluid fill-height class="ds-calendar-container">

      <ds-gestures
        @swipeleft="next"
        @swiperight="prev">

        <div v-if="currentType.schedule" class="ds-expand">

          <ds-agenda
            v-bind="{$scopedSlots}"
            v-on="$listeners"
            :calendar="calendar"
            @add="add"
            @edit="edit"
            @view-day="viewDay"
          ></ds-agenda>

        </div>

        <div v-else class="ds-expand">
          <ds-calendar ref="calendar"
                       v-bind="{$scopedSlots}"
                       v-on="$listeners"
                       :calendar="calendar"
                       @add="add"
                       @add-at="addAt"
                       @edit="edit"
                       @view-day="viewDay"
                       @added="handleAdd"
                       @moved="handleMove"
          ></ds-calendar>
        </div>

      </ds-gestures>

      <ds-event-dialog ref="eventDialog"
                       v-bind="{$scopedSlots}"
                       v-on="$listeners"
                       :calendar="calendar"
                       @saved="eventFinish"
                       @actioned="eventFinish"
      ></ds-event-dialog>

      <v-dialog ref="optionsDialog"
                v-model="optionsVisible"
                v-bind="optionsDialog"
                :fullscreen="$dayspan.fullscreenDialogs">
        <v-list>
          <template v-for="option in options">
            <v-list-tile :key="option.text" @click="chooseOption( option )">
              {{ option.text }}
            </v-list-tile>
          </template>
        </v-list>
      </v-dialog>
      <v-dialog ref="promptDialog"
                v-model="promptVisible"
                v-bind="promptDialog">
        <v-card>
          <v-card-title>{{ promptQuestion }}</v-card-title>
          <v-card-actions>
            <v-btn color="primary" flat @click="choosePrompt( true )">
              {{ labels.promptConfirm }}
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn color="secondary" flat @click="choosePrompt( false )">
              {{ labels.promptCancel }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </div>
</template>

<script>
/* eslint-disable */ // TODO: REMOVE IT

import { Units, Sorts, Calendar, Op, Month, Weekday } from 'dayspan';

const c = Calendar.months();

c.addEvents([
  {
    data: {
      title: 'Weekly Meeting',
      color: '#3F51B5',
    },
    schedule: {
      dayOfWeek: [Weekday.MONDAY],
      times: [9],
      duration: 30,
      durationUnit: 'minutes',
    },
  },
  {
    data: {
      title: 'First Weekend',
      color: '#4CAF50',
    },
    schedule: {
      weekspanOfMonth: [0],
      dayOfWeek: [Weekday.FRIDAY],
      duration: 3,
      durationUnit: 'days',
    },
  },
  {
    data: {
      title: 'End of Month',
      color: '#000000',
    },
    schedule: {
      lastDayOfMonth: [1],
      duration: 1,
      durationUnit: 'hours',
    },
  },
  {
    data: {
      title: 'Mother\'s Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.MAY],
      dayOfWeek: [Weekday.SUNDAY],
      weekspanOfMonth: [1],
    },
  },
  {
    data: {
      title: 'New Year\'s Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.JANUARY],
      dayOfMonth: [1],
    },
  },
  {
    data: {
      title: 'Inauguration Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.JANUARY],
      dayOfMonth: [20],
    },
  },
  {
    data: {
      title: 'Martin Luther King, Jr. Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.JANUARY],
      dayOfWeek: [Weekday.MONDAY],
      weekspanOfMonth: [2],
    },
  },
  {
    data: {
      title: 'George Washington\'s Birthday',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.FEBRUARY],
      dayOfWeek: [Weekday.MONDAY],
      weekspanOfMonth: [2],
    },
  },
  {
    data: {
      title: 'Memorial Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.MAY],
      dayOfWeek: [Weekday.MONDAY],
      lastWeekspanOfMonth: [0],
    },
  },
  {
    data: {
      title: 'Independence Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.JULY],
      dayOfMonth: [4],
    },
  },
  {
    data: {
      title: 'Labor Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.SEPTEMBER],
      dayOfWeek: [Weekday.MONDAY],
      lastWeekspanOfMonth: [0],
    },
  },
  {
    data: {
      title: 'Columbus Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.OCTOBER],
      dayOfWeek: [Weekday.MONDAY],
      weekspanOfMonth: [1],
    },
  },
  {
    data: {
      title: 'Veterans Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.NOVEMBER],
      dayOfMonth: [11],
    },
  },
  {
    data: {
      title: 'Thanksgiving Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.NOVEMBER],
      dayOfWeek: [Weekday.THURSDAY],
      weekspanOfMonth: [3],
    },
  },
  {
    data: {
      title: 'Christmas Day',
      color: '#2196F3',
      calendar: 'US Holidays',
    },
    schedule: {
      month: [Month.DECEMBER],
      dayOfMonth: [25],
    },
  },
]);

export default {

  name: 'dsCalendarApp',
  props:
    {
      events: {
        type: Array,
      },
      calendar: {
        type: Calendar,
        default() {
          return c;
        },
      },
      types: {
        type: Array,
        default: () => [
          {
            id: 'M',
            label: 'Month',
            shortcut: 'M',
            type: Units.MONTH,
            size: 1,
            focus: 0.4999,
            repeat: true,
            listTimes: false,
            updateRows: true,
            schedule: false,
          },
          {
            id: 'W',
            label: 'Week',
            shortcut: 'W',
            type: Units.WEEK,
            size: 1,
            focus: 0.4999,
            repeat: true,
            listTimes: true,
            updateRows: true,
            schedule: false,
          },
          {
            id: 'D',
            label: 'Day',
            shortcut: 'D',
            type: Units.DAY,
            size: 1,
            focus: 0.4999,
            repeat: true,
            listTimes: true,
            updateRows: true,
            schedule: false,
          },
        ],
      },
      allowsAddToday: {
        type: Boolean,
        default() {
          return this.$dsDefaults().allowsAddToday;
        },
      },
      formats: {
        validate(x) {
          return this.$dsValidate(x, 'formats');
        },
        default() {
          return this.$dsDefaults().formats;
        },
      },
      labels: {
        validate(x) {
          return this.$dsValidate(x, 'labels');
        },
        default() {
          return this.$dsDefaults().labels;
        },
      },
      styles: {
        validate(x) {
          return this.$dsValidate(x, 'styles');
        },
        default() {
          return this.$dsDefaults().styles;
        },
      },
      optionsDialog: {
        validate(x) {
          return this.$dsValidate(x, 'optionsDialog');
        },
        default() {
          return this.$dsDefaults().optionsDialog;
        },
      },
      promptDialog: {
        validate(x) {
          return this.$dsValidate(x, 'promptDialog');
        },
        default() {
          return this.$dsDefaults().promptDialog;
        },
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
    currentType:
      {
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

    nextLabel() {
      return this.labels.next(this.currentType);
    },

    prevLabel() {
      return this.labels.prev(this.currentType);
    },

    hasCreatePopover() {
      return !!this.$scopedSlots.eventCreatePopover;
    },
  },

  watch: {
    events: 'applyEvents',
    calendar: 'applyEvents',
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

    edit(calendarEvent) {
      this.$refs.eventDialog.edit(calendarEvent);
    },

    editPlaceholder(createEdit) {
      const placeholder = createEdit.calendarEvent;
      const { details } = createEdit;
      const { eventDialog } = this.$refs;
      const { calendar } = this.$refs;

      eventDialog.addPlaceholder(placeholder, details);
      eventDialog.$once('close', calendar.clearPlaceholder);
    },

    add(day) {
      if (!this.$dayspan.features.addDay) {
        return;
      }

      const { eventDialog } = this.$refs;
      const { calendar } = this.$refs;
      const useDialog = !this.hasCreatePopover;

      calendar.addPlaceholder(day, true, useDialog);

      if (useDialog) {
        eventDialog.add(day);
        eventDialog.$once('close', calendar.clearPlaceholder);
      }
    },

    addAt(dayHour) {
      if (!this.$dayspan.features.addTime) {
        return;
      }

      const { eventDialog } = this.$refs;
      const { calendar } = this.$refs;
      const useDialog = !this.hasCreatePopover;
      const at = dayHour.day.withHour(dayHour.hour);

      calendar.addPlaceholder(at, false, useDialog);

      if (useDialog) {
        eventDialog.addAt(dayHour.day, dayHour.hour);
        eventDialog.$once('close', calendar.clearPlaceholder);
      }
    },

    addToday() {
      if (!this.$dayspan.features.addDay) {
        return;
      }

      const { eventDialog } = this.$refs;
      const { calendar } = this.$refs;
      const useDialog = !this.hasCreatePopover || !calendar;

      let day = this.$dayspan.today;

      if (!this.calendar.filled.matchesDay(day)) {
        const first = this.calendar.days[0];
        const last = this.calendar.days[this.calendar.days.length - 1];
        const firstDistance = Math.abs(first.currentOffset);
        const lastDistance = Math.abs(last.currentOffset);

        day = firstDistance < lastDistance ? first : last;
      }

      if (calendar) {
        calendar.addPlaceholder(day, true, useDialog);
      }

      if (useDialog) {
        eventDialog.add(day);

        if (calendar) {
          eventDialog.$once('close', calendar.clearPlaceholder);
        }
      }
    },

    handleAdd(addEvent) {
      const { eventDialog } = this.$refs;
      const { calendar } = this.$refs;

      addEvent.handled = true;

      if (!this.hasCreatePopover) {
        if (addEvent.placeholder.fullDay) {
          eventDialog.add(addEvent.span.start, addEvent.span.days(Op.UP));
        } else {
          eventDialog.addSpan(addEvent.span);
        }

        eventDialog.$once('close', addEvent.clearPlaceholder);
      } else {
        calendar.placeholderForCreate = true;
      }
    },

    handleMove(moveEvent) {
      const calendarEvent = moveEvent.calendarEvent;
      const target = moveEvent.target;
      const targetStart = target.start;
      const sourceStart = calendarEvent.time.start;
      const schedule = calendarEvent.schedule;
      const options = [];

      moveEvent.handled = true;

      const callbacks = {
        cancel: () => {
          moveEvent.clearPlaceholder();
        },
        single: () => {
          calendarEvent.move(targetStart);
          this.eventsRefresh();
          moveEvent.clearPlaceholder();

          this.$emit('event-update', calendarEvent.event);
        },
        instance: () => {
          calendarEvent.move(targetStart);
          this.eventsRefresh();
          moveEvent.clearPlaceholder();

          this.$emit('event-update', calendarEvent.event);
        },
        duplicate: () => {
          schedule.setExcluded(targetStart, false);
          this.eventsRefresh();
          moveEvent.clearPlaceholder();

          this.$emit('event-update', calendarEvent.event);
        },
        all: () => {
          schedule.moveTime(sourceStart.asTime(), targetStart.asTime());
          this.eventsRefresh();
          moveEvent.clearPlaceholder();

          this.$emit('event-update', calendarEvent.event);
        },
      };

      options.push({
        text: this.labels.moveCancel,
        callback: callbacks.cancel,
      });

      if (schedule.isSingleEvent()) {
        options.push({
          text: this.labels.moveSingleEvent,
          callback: callbacks.single,
        });

        if (this.$dayspan.features.moveDuplicate) {
          options.push({
            text: this.labels.moveDuplicate,
            callback: callbacks.duplicate,
          });
        }
      } else {
        if (this.$dayspan.features.moveInstance) {
          options.push({
            text: this.labels.moveOccurrence,
            callback: callbacks.instance,
          });
        }

        if (this.$dayspan.features.moveDuplicate) {
          options.push({
            text: this.labels.moveDuplicate,
            callback: callbacks.duplicate,
          });
        }

        if (this.$dayspan.features.moveAll &&
          !schedule.isFullDay() &&
          targetStart.sameDay(sourceStart)) {
          options.push({
            text: this.labels.moveAll,
            callback: callbacks.all,
          });
        }
      }

      this.options = options;
      this.optionsVisible = true;
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

    eventsRefresh() {
      this.calendar.refreshEvents();

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
  .ds-app-calendar-toolbar {

    .v-toolbar__content {
      border-bottom: 1px solid rgb(224, 224, 224);
    }
  }

  .ds-skinny-button {
    margin-left: 2px !important;
    margin-right: 2px !important;
  }

  .ds-expand {
    width: 100%;
    height: 100%;
  }

  .ds-calendar-container {
    padding: 0 !important;
    position: relative;
  }

  .v-btn--floating.ds-add-event-today {
    .v-icon {
      width: 24px;
      height: 24px;
    }
  }

</style>
