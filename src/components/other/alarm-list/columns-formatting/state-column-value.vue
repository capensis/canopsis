<template lang="pug">
  div
    span(:class="[`bg-${color}`, 'badge']") {{ text }}
    v-icon(:color="color", v-if="showIcon") account_circle
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

export default {
  props: {
    stateId: {
      required: true,
    },
    showIcon: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      colors: {
        [ENTITIES_STATES.ok]: 'green',
        [ENTITIES_STATES.minor]: 'yellow',
        [ENTITIES_STATES.major]: 'orange',
        [ENTITIES_STATES.critical]: 'red',
      },
    };
  },
  computed: {
    color() {
      const color = this.colors[this.stateId];

      if (color) {
        return color;
      }

      return 'purple';
    },
    text() {
      const title = this.$t(`tables.alarmStates.${this.stateId}`);

      if (title) {
        return title;
      }

      return 'Unknown';
    },
  },
};
</script>

<style scoped>
  .badge {
    display: inline-block;
    min-width: 10px;
    padding: 3px 7px;
    font-size: 12px;
    font-weight: 700;
    line-height: 1;
    color: #fff;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    background-color: #777;
    border-radius: 10px;
  }

  .badge.bg-green {
    background-color: #00a65a
  }

  .badge.bg-red {
    background-color: #CF0000
  }

  .badge.bg-yellow {
    background-color: #FFE000
  }

  .badge.bg-orange {
    background-color: #FF9900
  }
</style>
