<template lang="pug">
  v-menu(
    v-model="menu",
    :activator="activator",
    :close-on-content-click="false",
    nudge-width="200",
    max-width="500",
    transition="fade-transition",
    ignore-click-upper-outside,
    offset-overflow,
    offset-x
  )
    v-card
      v-card-text
        v-layout(
          v-for="(event, index) in includedEvents",
          :key="`popover-event-${index}`",
          row,
          wrap
        )
          v-flex(xs12)
            div.ds-calendar-event-popover-item(
              :style="{ backgroundColor: getStyleColor(event) }",
              @click="$emit('event-click', event)"
            )
              strong {{ event.data.title }}
              p {{ event.data.description }}
</template>

<script>
import { get } from 'lodash';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

export default {
  props: {
    activator: {
      type: HTMLElement,
      required: true,
    },
    calendarEvent: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      menu: true,
      timer: null,
    };
  },
  computed: {
    includedEvents() {
      return get(this.calendarEvent, 'data.meta.events', []);
    },
  },
  watch: {
    menu(menu) {
      if (this.timer) {
        clearTimeout(this.timer);

        this.timer = null;
      }

      if (!menu) {
        this.timer = setTimeout(() => this.$emit('closed'), VUETIFY_ANIMATION_DELAY);
      }
    },
  },
  methods: {
    getStyleColor(event) {
      const past = event.schedule.end.isBefore(new Date());

      return this.$dayspan.getStyleColor(this.calendarEvent.data, event, past);
    },
  },
};
</script>

<style lang="scss">
.ds-calendar-event-popover-item {
  color: white;
  margin: 1px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding-left: 0.5em;
  cursor: pointer;
  border-radius: 2px;
}
</style>
