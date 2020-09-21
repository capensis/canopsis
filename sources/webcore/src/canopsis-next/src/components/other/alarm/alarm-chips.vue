<template lang="pug">
  div.badge-container
    span.badge(:style="{ backgroundColor: style.color }") {{ style.text }}
</template>

<script>
import { ENTITY_INFOS_TYPE } from '@/constants';

import { formatState, formatStatus } from '@/helpers/state-and-status-formatting';

/**
 * Chips for the state/status of the alarm on timeline
 *
 * @module alarm
 *
 * @prop {Number,String} [value] - Value of the state or the status of the alarm
 * @prop {Boolean} [isStatus] - Boolean to determine if this is a state, or a status
 */
export default {
  props: {
    value: {
      type: [Number, String],
      default: 0,
    },
    type: {
      type: String,
      default: null,
    },
  },
  computed: {
    style() {
      if (this.type === ENTITY_INFOS_TYPE.status) {
        return formatStatus(this.value);
      }

      return formatState(this.value);
    },
  },
};
</script>
<style scoped>
  .badge-container {
    display: inline-block;
  }
  .badge {
    padding: 3px 7px;
    font-size: 14px;
    color: #fff;
    white-space: nowrap;
    border-radius: 10px;
  }
</style>
