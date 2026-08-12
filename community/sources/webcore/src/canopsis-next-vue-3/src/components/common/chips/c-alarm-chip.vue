<template>
  <div
    :class="{ 'chip-container--small': small }"
    class="chip-container"
    @click="$emit('click', $event)"
  >
    <v-badge :model-value="!!badgeValue" color="secondary">
      <template #badge>
        <span class="px-1">
          {{ badgeValue }}
        </span>
      </template>

      <span :style="{ backgroundColor: style.color }" :class="style.class" class="chip">
        {{ style.text }}
        <slot name="append" />
      </span>
    </v-badge>
  </div>
</template>

<script>
import { ENTITY_INFOS_TYPE } from '@/constants';
import { formatAlarmState, formatAlarmStatus } from '@/helpers/entities/alarm/formatting';

/**
 * Chips for the state/status of the alarm on timeline
 *
 * @module alarm
 *
 * @prop {Number,String} [value] - Value of the state or the status of the alarm
 * @prop {Boolean} [isStatus] - Boolean to determine if this is a state, or a status
 */
export default {
  emits: ['click'],

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
      return this.type === ENTITY_INFOS_TYPE.status
        ? formatAlarmStatus(this.value)
        : formatAlarmState(this.value);
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
    white-space: nowrap;
    border-radius: 10px;
  }

  :deep(.v-badge__badge) {
    inset: -14px -6px auto auto !important;
    width: auto;
    min-width: 16px;
    height: 16px;
    padding: 0;
    font-size: 10px;
    line-height: 16px;
    border-radius: 5px;
  }

  &--small {
    line-height: 1;

    .chip {
      display: block;
      padding: 1px 5px;
      font-size: 12px;
    }

    :deep(.v-badge__badge) {
      min-width: 12px;
      height: 12px;
      font-size: 8px;
      line-height: 15px;
      inset: -9px -6px auto auto !important;
    }
  }
}
</style>
