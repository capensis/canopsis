<template>
  <span
    class="point-icon"
    v-on="$listeners"
    :style="pointStyles"
  >
    <v-icon
      :size="icon.size || size"
      :color="icon.color"
    >
      {{ icon.name }}
    </v-icon>
  </span>
</template>

<script>
import { isNumber } from 'lodash';

import { ENTITIES_STATES } from '@/constants';
import { CSS_COLORS_VARS } from '@/config';

import { isNotActivePbehaviorType } from '@/helpers/entities/pbehavior/form';
import { getEntityColor, getEntityStateColor } from '@/helpers/entities/entity/color';

export default {
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
      if (this.entity.state === ENTITIES_STATES.ok) {
        return {
          name: 'check_circle_outline',
          color: 'white',
          backgroundColor: getEntityStateColor(ENTITIES_STATES.ok),
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

    icon() {
      if (!this.entity) {
        return {
          name: 'link',
          color: 'grey darken-2',
        };
      }

      if (this.pbehaviorEnabled && this.isNotActivePbehavior) {
        return {
          name: this.entity.pbehavior_info.icon_name,
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
        color: 'grey darken-2',
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
