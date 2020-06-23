<template lang="pug">
  ds-calendar-app(
    :events="events",
    :readOnly="readOnly",
    fluid,
    fillHeight,
    @moved="moveHandler",
    @resized="resizeHandler",
    @added="addHandler"
  )
    ds-calendar-event-popover(
      slot="eventCreatePopover",
      slot-scope="props",
      v-bind="props"
    )
      pbehavior-create-event(slot-scope="{ close }", :close="close")
</template>

<script>
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
  computed: {
    events() {
      return [{
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
      }, {
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
      }];
    },
  },
  methods: {
    moveHandler(event) {
      this.$emit('moved', event.calendarEvent);
      event.clearPlaceholder();
    },
    resizeHandler(event) {
      this.$emit('resized', event.calendarEvent);
      event.clearPlaceholder();
    },
    addHandler(event) {
      this.$emit('added', event.calendarEvent);
      event.clearPlaceholder();
    },
  },
};
</script>
