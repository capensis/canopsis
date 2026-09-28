<template>
  <span v-bind="$attrs" :style="style" class="c-alarm-extra-details-chip">
    <slot :icon-color="preparedIconColor">
      <v-icon v-if="icon" :color="preparedIconColor" size="small">
        {{ icon }}
      </v-icon>
    </slot>
  </span>
</template>

<script>
import { computed, inject } from 'vue';

import { COLORS } from '@/config';

import { getMostReadableTextColor } from '@/helpers/color';

export default {
  inheritAttrs: false,

  props: {
    color: {
      type: String,
      default: COLORS.secondary,
    },
    themeColor: {
      type: String,
      default: '',
    },
    iconColor: {
      type: String,
      default: '',
    },
    icon: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const system = inject('$system', {});
    const preparedColor = computed(
      () =>
        (props.themeColor && system.theme?.colors?.icons?.[props.themeColor]) ||
        COLORS.alarmExtraDetails[props.themeColor] ||
        props.color,
    );
    const preparedIconColor = computed(
      () =>
        props.iconColor ||
        getMostReadableTextColor(preparedColor.value, { level: 'AA', size: 'large' }),
    );
    const style = computed(() => ({ backgroundColor: preparedColor.value }));

    return { preparedIconColor, style };
  },
};
</script>

<style lang="scss" scoped>
.c-alarm-extra-details-chip {
  display: flex;
  height: 22px;
  padding: 3px 7px;
  text-align: center;
  white-space: nowrap;
  vertical-align: baseline;
  border-radius: 10px;
  line-height: 16px;
}
</style>
