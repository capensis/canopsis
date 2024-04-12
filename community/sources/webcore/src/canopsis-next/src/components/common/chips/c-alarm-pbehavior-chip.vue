<template>
  <c-alarm-extra-details-chip
    :icon="preparedIcon"
    :color="preparedColor"
    :icon-color="iconColor"
    v-on="$listeners"
  />
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';

import { getMostReadableTextColor } from '@/helpers/color';

export default {
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
    const iconColor = computed(() => getMostReadableTextColor(preparedColor.value, { level: 'AA', size: 'large' }));

    return {
      preparedIcon,
      preparedColor,
      iconColor,
    };
  },
};
</script>
