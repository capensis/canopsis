<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.colorIndicator.title') }}
    v-container
      color-indicator-field(v-field="value")
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

import ColorIndicatorField from '../partials/color-indicator-field.vue';

export default {
  components: { ColorIndicatorField },
  props: {
    value: {
      type: String,
      default: COLOR_INDICATOR_TYPES.state,
    },
  },
};
</script>
