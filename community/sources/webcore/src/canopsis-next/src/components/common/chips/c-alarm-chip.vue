<template lang="pug">
  div.chip-container
    v-badge(:value="!!badgeValue", color="secondary", overlap)
      template(#badge="")
        span.px-1 {{ badgeValue }}
      span.chip(:style="{ backgroundColor: style.color }") {{ style.text }}
</template>

<script>
import { ENTITY_INFOS_TYPE } from '@/constants';

import { formatState, formatStatus } from '@/helpers/formatting';

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
      default: ENTITY_INFOS_TYPE.state,
    },
    badgeValue: {
      type: [Number, String],
      default: 0,
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
<style lang="scss" scoped>
  .chip-container {
    display: inline-block;

    .chip {
      padding: 3px 7px;
      font-size: 14px;
      color: #fff;
      white-space: nowrap;
      border-radius: 10px;
    }

    & /deep/ .v-badge--overlap .v-badge__badge {
      font-size: 10px;
      border-radius: 5px;
      min-width: 16px;
      width: auto;
      height: 16px;
      top: -12px;
      right: -6px;
    }
  }
</style>
