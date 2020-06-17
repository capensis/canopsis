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
</template>

<script>
export default {
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
