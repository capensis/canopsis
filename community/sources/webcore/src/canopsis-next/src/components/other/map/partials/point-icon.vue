<template lang="pug">
  v-icon.point-icon(
    v-on="$listeners",
    :style="{ backgroundColor: icon.backgroundColor }",
    :size="size",
    :color="icon.color"
    ) {{ icon.name }}
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

import { getEntityColor } from '@/helpers/color';

import { COLORS } from '@/config';

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
  },
  computed: {
    icon() {
      if (!this.entity) {
        return {
          name: 'link',
          color: 'grey darken-2',
        };
      }

      if (this.entity?.state) {
        const state = this.entity.state.val;

        if (state === ENTITIES_STATES.ok) {
          return {
            name: 'check_circle_outline',
            color: 'white',
            backgroundColor: COLORS.primary,
          };
        }

        return {
          color: getEntityColor(this.entity, this.colorIndicator),
          name: 'warning',
        };
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
  border-radius: 50%;
}
</style>
