<template>
  <v-menu
    v-model="menu"
    :activator="activator"
    :close-on-content-click="false"
    nudge-width="200"
    max-width="500"
    transition="fade-transition"
    ignore-click-upper-outside
    offset-overflow
    offset-x
  >
    <v-card>
      <v-card-text>
        <v-layout
          v-for="(event, index) in includedEvents"
          :key="`popover-event-${index}`"
          wrap
        >
          <v-flex xs12>
            <div
              :style="{ backgroundColor: getStyleColor(event) }"
              class="stats-calendar-menu__item"
              @click="$emit('click:event', event)"
            >
              <strong>{{ event.name }}</strong>
              <p>{{ event.description }}</p>
            </div>
          </v-flex>
        </v-layout>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script>
import { get } from 'lodash';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { isDateBefore } from '@/helpers/date/date';
import { colorToRgba } from '@/helpers/color';

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
      const past = isDateBefore(event.end, new Date());

      return past ? colorToRgba(event.color, 0.5) : event.color;
    },
  },
};
</script>

<style lang="scss">
.stats-calendar-menu {
  &__item {
    color: white;
    margin: 1px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    padding-left: 0.5em;
    cursor: pointer;
    border-radius: 2px;
  }
}
</style>
