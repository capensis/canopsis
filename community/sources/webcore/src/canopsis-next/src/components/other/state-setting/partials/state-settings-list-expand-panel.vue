<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $tc('common.information') }}
    v-tab-item
      v-layout.pa-3
        v-flex(xs12)
          v-card.pa-3
            v-layout
              v-flex(offset-xs1, xs10)
                state-setting-information(:state-setting="stateSetting")
</template>

<script>
import StateSettingInformation from './state-setting-information.vue';

export default {
  components: { StateSettingInformation },
  props: {
    stateSetting: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
