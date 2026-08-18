<template>
  <v-menu
    v-model="menu"
    :target="activator"
    :close-on-content-click="false"
    max-width="500"
    transition="fade-transition"
    location="end"
  >
    <v-card>
      <v-card-text>
        <c-layout
          v-for="(event, index) in includedEvents"
          :key="`popover-event-${index}`"
          wrap
        >
          <div
            :style="{ backgroundColor: getStyleColor(event) }"
            class="stats-calendar-menu__item"
            role="button"
            tabindex="0"
            @click="emitClick(event)"
            @keydown.enter="emitClick(event)"
          >
            <strong>{{ event.name }}</strong>
            <p>{{ event.description }}</p>
          </div>
        </c-layout>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script>
import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { colorToRgba } from '@/helpers/color';
import { isDateBefore } from '@/helpers/date/date';

const getCalendarEventSource = event => event?.input ?? event;

export default {
  props: {
    activator: {
      type: Object,
      required: true,
    },
    calendarEvent: {
      type: Object,
      required: true,
    },
  },

  emits: ['click:event', 'closed'],

  data() {
    return {
      menu: true,
      timer: null,
    };
  },

  computed: {
    includedEvents() {
      const source = getCalendarEventSource(this.calendarEvent);

      return source?.data?.meta?.events ?? [];
    },
  },

  watch: {
    menu(value) {
      if (this.timer) {
        clearTimeout(this.timer);
        this.timer = null;
      }

      if (!value) {
        this.timer = setTimeout(() => this.$emit('closed'), VUETIFY_ANIMATION_DELAY);
      }
    },
  },

  beforeUnmount() {
    clearTimeout(this.timer);
  },

  methods: {
    getStyleColor(event) {
      const source = getCalendarEventSource(event);
      const past = isDateBefore(source.end, new Date());

      return past ? colorToRgba(source.color, 0.5) : source.color;
    },

    emitClick(event) {
      this.$emit('click:event', event);
    },
  },
};
</script>

<style lang="scss" scoped>
.stats-calendar-menu__item {
  width: 100%;
  margin: 1px;
  overflow: hidden;
  padding-left: 0.5em;
  border-radius: 2px;
  color: white;
  cursor: pointer;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
