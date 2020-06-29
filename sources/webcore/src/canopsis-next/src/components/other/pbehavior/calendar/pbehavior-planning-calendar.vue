<template lang="pug">
  ds-calendar-app(
    :events="events",
    :readOnly="readOnly",
    fluid,
    fillHeight,
    @changed="changeHandler",
    @added="addHandler"
  )
    ds-calendar-event-popover(
      slot="eventPopover",
      slot-scope="props",
      v-bind="props"
    )
      pbehavior-create-event(
        slot-scope="{ placeholder, close, edit }",
        :calendarEvent="placeholder",
        @close="close",
        @submit="edit"
      )
    ds-calendar-event-popover(
      slot="eventCreatePopover",
      slot-scope="props",
      v-bind="props"
    )
      pbehavior-create-event(
        slot-scope="{ placeholder, close, add }",
        :calendarEvent="placeholder",
        @close="close",
        @submit="add"
      )
</template>

<script>
import Vue from 'vue';

import uuid from '@/helpers/uuid';

import PbehaviorCreateEvent from './partials/pbehavior-create-event.vue';

export default {
  components: { PbehaviorCreateEvent },
  props: {
    pbehaviors: {
      type: Array,
      default: () => [],
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      eventsById: {
        'e3b23564-71d7-4d59-8105-5ece80032100': {
          id: 'e3b23564-71d7-4d59-8105-5ece80032100',
          data: {
            title: 'Meeting',
            description: '',
            location: '',
            color: '#F44336',
            forecolor: '#ffffff',
            calendar: '',
            busy: true,
            icon: '',
          },
          schedule: {
            dayOfMonth: [16],
            duration: 270,
            durationUnit: 'minutes',
            month: [5],
            times: ['16:15'],
            year: [2020],
          },
        },
        '7c5db84d-cd0a-4af6-b4fd-c590cba585bb': {
          id: '7c5db84d-cd0a-4af6-b4fd-c590cba585bb',
          data: {
            title: 'Vocation',
            description: '',
            location: '',
            color: '#FFC107',
            forecolor: '#ffffff',
            calendar: '',
            busy: true,
            icon: '',
          },
          schedule: {
            dayOfMonth: [7],
            duration: 7,
            month: [5],
            year: [2020],
          },
        },
      },
    };
  },
  computed: {
    events() {
      return Object.values(this.eventsById);
    },
  },
  methods: {
    changeHandler(event) {
      this.saveEvent(event.calendarEvent);
      event.clearPlaceholder();
    },
    addHandler(event) {
      this.saveEvent(event.calendarEvent);
      event.clearPlaceholder();
    },
    saveEvent({ id = uuid(), data, schedule }) {
      Vue.set(this.eventsById, [id], {
        id,
        data,
        schedule,
      });
    },
  },
};
</script>
