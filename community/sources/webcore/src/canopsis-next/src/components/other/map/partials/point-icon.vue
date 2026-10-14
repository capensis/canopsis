<template>
  <span v-bind="$attrs" :style="pointStyles" class="point-icon">
    <v-icon :size="icon.size || size" :color="icon.color">
      {{ icon.name }}
    </v-icon>
  </span>
</template>

<script>
import { isNumber } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { ALARM_STATES, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { getAlarmStateColor } from '@/helpers/entities/alarm/color';
import { getEntityColor } from '@/helpers/entities/entity/color';
import { isNotActivePbehaviorType } from '@/helpers/entities/pbehavior/form';

export default {
  inheritAttrs: false,

  props: {
    entity: {
      type: [String, Object],
      required: false,
    },
    size: {
      type: Number,
      default: 24,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      required: false,
    },
    pauseEnabled: {
      type: Boolean,
      required: false,
    },
  },

  computed: {
    pointStyles() {
      return {
        width: `${this.size}px`,
        height: `${this.size}px`,
        backgroundColor: this.icon.backgroundColor,
      };
    },

    entityIcon() {
      if (this.entity.state === ALARM_STATES.ok) {
        return {
          name: 'check_circle_outline',
          color: 'white',
          backgroundColor: getAlarmStateColor(ALARM_STATES.ok),
        };
      }

      return {
        backgroundColor: getEntityColor(this.entity, this.colorIndicator),
        color: 'white',
        size: this.size - 8,
        name: 'warning',
      };
    },

    isNotActivePbehavior() {
      return isNotActivePbehaviorType(this.entity.pbehavior_info?.canonical_type);
    },

    hasNonPausePbehavior() {
      return this.isNotActivePbehavior &&
        this.entity.pbehavior_info?.canonical_type !== PBEHAVIOR_TYPE_TYPES.pause;
    },

    icon() {
      if (!this.entity) {
        return {
          name: 'link',
          color: 'grey-darken-2',
        };
      }

      if (this.pbehaviorEnabled && this.hasNonPausePbehavior) {
        return {
          name: this.entity.pbehavior_info.icon_name,
          color: 'white',
          size: this.size - 8,
          backgroundColor: CSS_COLORS_VARS.secondary,
        };
      }

      if (this.pauseEnabled && this.entity.pause) {
        return {
          name: 'pause',
          color: 'white',
          size: this.size - 8,
          backgroundColor: CSS_COLORS_VARS.secondary,
        };
      }

      if (this.pbehaviorEnabled && this.isNotActivePbehavior) {
        return {
          name: '$vuetify.icons.pbehavior',
          color: 'white',
          size: this.size - 8,
          backgroundColor: CSS_COLORS_VARS.secondary,
        };
      }

      if (isNumber(this.entity?.state)) {
        return this.entityIcon;
      }

      return {
        name: 'location_on',
        color: 'grey-darken-2',
      };
    },
  },
};
</script>

<style lang="scss">
.point-icon {
  word-break: initial;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}
</style>
