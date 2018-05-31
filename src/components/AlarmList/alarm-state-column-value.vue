<template lang="pug">
  div
    span(:class="[`bg-${color}`, 'badge']") {{ text }}
    v-icon(:color="color") account_circle
</template>

<script>
import { STATES } from '@/config';

export default {
  props: {
    stateId: {
      required: true,
    },
  },
  computed: {
    color() {
      switch (this.stateId) {
        case STATES.info:
          return 'green';
        case STATES.minor:
          return 'yellow';
        case STATES.major:
          return 'orange';
        case STATES.critical:
          return 'red';
        default:
          return 'purple';
      }
    },
    text() {
      switch (this.stateId) {
        case STATES.info:
          return 'Info';
        case STATES.minor:
          return 'Minor';
        case STATES.major:
          return 'Major';
        case STATES.critical:
          return 'Critical';
        default:
          return 'Unknown';
      }
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
