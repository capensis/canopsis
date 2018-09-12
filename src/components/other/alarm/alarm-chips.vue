<template lang="pug">
  div
    span.badge(:style="{backgroundColor : style.color}") {{ style.text }}
</template>

<script>
import { formatState, formatStatus } from '@/helpers/state-and-status-formatting';
import { ENTITY_INFOS_TYPE } from '@/constants';

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

</style>
