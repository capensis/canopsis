<template>
  <div
    :class="{ 'chip-container--small': small }"
    class="chip-container"
    @click="$emit('click', $event)"
  >
    <v-badge
      :value="!!badgeValue"
      color="secondary"
      overlap
    >
      <template #badge="">
        <span class="px-1">{{ badgeValue }}</span>
      </template>
      <span
        :style="{ backgroundColor: style.color }"
        class="chip"
      >
        {{ style.text }}
      </span>
    </v-badge>
  </div>
</template>

<script>
import { ENTITY_INFOS_TYPE } from '@/constants';

import { formatState, formatStatus } from '@/helpers/entities/entity/formatting';

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
    small: {
      type: Boolean,
      default: false,
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

    & ::v-deep .v-badge--overlap .v-badge__badge {
      padding: 0;
      font-size: 10px;
      border-radius: 5px;
      min-width: 16px;
      width: auto;
      height: 16px;
      inset: -14px -6px auto auto !important;
      line-height: 16px;
    }

    &--small {
      line-height: 1;

      .chip {
        display: block;
        padding: 1px 5px;
        font-size: 12px;
      }

      & ::v-deep .v-badge--overlap {
        display: block;

        .v-badge__badge {
          font-size: 8px;
          min-width: 12px;
          height: 12px;
          line-height: 15px;
          inset: -9px -6px auto auto !important;
        }
      }
    }
  }
</style>
