<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.colorIndicator.title') }}
    v-container
      c-color-indicator-field(v-field="value")
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      default: COLOR_INDICATOR_TYPES.state,
    },
  },
};
</script>
