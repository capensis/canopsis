<template>
  <div class="c-clickable-tooltip">
    <div
      ref="activator"
      @focusin="showTooltip"
      @focusout="hideTooltip"
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
      ignore-content-leave="ignore-content-leave"
    >
      <span
        @focusin="showTooltip"
        @focusout="hideTooltip"
      >
        <slot /></span>
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
