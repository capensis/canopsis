<template>
  <span :style="style" class="c-alarm-extra-details-chip" v-on="$listeners">
    <slot>
      <v-icon v-if="icon" :color="iconColor" small>
        {{ icon }}
      </v-icon>
    </slot>
  </span>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';

export default {
  props: {
    color: {
      type: String,
      default: COLORS.secondary,
    },
    iconColor: {
      type: String,
      default: 'white',
    },
    icon: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const style = computed(() => ({ backgroundColor: props.color }));

    return { style };
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
