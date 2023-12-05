<template>
  <div class="c-clickable-tooltip">
    <div
      ref="activator"
      @mouseenter="showTooltip"
      @mouseleave="hideTooltip"
    >
      <slot name="activator" />
    </div>
    <v-tooltip
      :value="show"
      :top="top"
      :right="right"
      :bottom="bottom"
      :left="left"
      :close-delay="transitionDelay"
      :open-delay="transitionDelay"
      :activator="$refs.activator"
      ignore-content-leave
      content-class="c-clickable-tooltip__content"
    >
      <div
        @mouseenter="showTooltip"
        @mouseleave="hideTooltip"
      >
        <slot />
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { debounce } from 'lodash';

const DEFAULT_TRANSITION_DELAY = 200;

export default {
  props: {
    transitionDelay: {
      type: Number,
      default: DEFAULT_TRANSITION_DELAY,
    },
    top: {
      type: Boolean,
      required: false,
    },
    right: {
      type: Boolean,
      required: false,
    },
    bottom: {
      type: Boolean,
      required: false,
    },
    left: {
      type: Boolean,
      required: false,
    },
  },
  data() {
    return {
      show: false,
    };
  },
  created() {
    this.toggleTooltip = debounce(function toggleShow(value) {
      this.show = value;
    }, DEFAULT_TRANSITION_DELAY);
  },
  methods: {
    showTooltip() {
      this.toggleTooltip(true);
    },

    hideTooltip() {
      this.toggleTooltip(false);
    },
  },
};
</script>

<style lang="scss">
.c-clickable-tooltip {
  &__content {
    pointer-events: initial;
    padding: 0;

    & > div {
      padding: 5px 16px;
    }
  }
}
</style>
