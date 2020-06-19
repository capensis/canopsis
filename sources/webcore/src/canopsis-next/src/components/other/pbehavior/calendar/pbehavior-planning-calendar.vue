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
      return this.pbehaviors;
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
