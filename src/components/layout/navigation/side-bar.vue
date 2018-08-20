<template lang="pug">
  v-navigation-drawer.side-bar.grey.darken-4(
  v-model="isOpen",
  :clipped="$mq | mq({ mobile: false, laptop: true })",
  :width="width",
  absolute,
  app,
  )
    v-card(flat)
    v-expansion-panel(
    class="panel",
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content.grey.darken-4.white--text
        div(slot="header") View Group 1
        v-card.grey.darken-3.white--text
          v-card-text View 1
        v-card.grey.darken-3.white--text
          v-card-text View 2
    v-expansion-panel(
    class="panel",
    expand,
    focusable,
    dark
    )
      v-expansion-panel-content.grey.darken-4.white--text
        div(slot="header") View Group 2
        v-card.grey.darken-3.white--text
          v-card-text View 1
        v-card.grey.darken-3.white--text
          v-card-text View 2
    v-divider
    v-btn.addBtn(
    fab,
    dark,
    fixed,
    bottom,
    right,
    color="green darken-4"
    )
      v-icon(dark) add
</template>

<script>
import { SIDE_BAR_WIDTH } from '@/config';

/**
 * Component for the side-bar, on the left of the application
 *
 * @prop {bool} [value=false] - visibility control
 *
 * @event input#update
 */
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      width: SIDE_BAR_WIDTH,
    };
  },
  computed: {
    isOpen: {
      get() {
        return this.value;
      },
      set(value) {
        if (value !== this.value) {
          this.$emit('input', value);
        }
      },
    },
  },
};
</script>

<style scoped>
  a {
    color: inherit;
    text-decoration: none;
  }
  .panel {
    box-shadow: none;
  }

  .side-bar {
    position: fixed;
    height: 100vh;
    overflow-y: auto;
  }
</style>
