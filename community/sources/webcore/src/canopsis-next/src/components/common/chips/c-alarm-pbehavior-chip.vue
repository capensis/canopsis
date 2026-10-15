<template>
  <c-alarm-extra-details-chip
    v-bind="$attrs"
    :icon="preparedIcon"
    :color="preparedColor"
    :icon-color="iconColor"
  />
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';

import { getMostReadableTextColor } from '@/helpers/color';

export default {
  inheritAttrs: false,

  props: {
    icon: {
      type: String,
      default: 'pause',
    },
    color: {
      type: String,
      default: COLORS.secondary,
    },
  },

  setup(props) {
    const preparedIcon = computed(() => props.icon || 'pause');
    const preparedColor = computed(() => props.color || COLORS.secondary);
    const iconColor = computed(() => getMostReadableTextColor(preparedColor.value, { level: 'AA', size: 'large' }),
    );

    return {
      iconColor,
      preparedColor,
      preparedIcon,
    };
  },
};
</script>
