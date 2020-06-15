<template lang="pug">
  div
    v-layout
      v-checkbox(
        v-model="readOnly",
        label="Read only",
        color="primary"
      )
    v-layout
      ds-calendar-app.example-calendar(
        :events="events",
        :readOnly="readOnly",
        @moved="moveHandler",
        @resized="resizeHandler",
        @added="addHandler"
      )
</template>

<script>
import Vue from 'vue';
import uuid from '@/helpers/uuid';

export default {
  data() {
    return {
      readOnly: false,
      eventsById: {},
    };
  },
  computed: {
    events() {
      return Object.values(this.eventsById);
    },
  },
  methods: {
    moveHandler(event) {
      this.saveEvent(event.calendarEvent);
      event.clearPlaceholder();
    },
    resizeHandler(event) {
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

<style lang="scss" scoped>
  .example-calendar {
    /deep/ .ds-calendar-container {
      margin: 0;
      max-width: 100%;
      height: 700px;
    }
  }
</style>
