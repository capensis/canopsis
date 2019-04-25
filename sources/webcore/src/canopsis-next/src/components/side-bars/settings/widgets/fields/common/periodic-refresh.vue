<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.periodicRefresh')}}
      .font-italic.caption.ml-1 ({{ $t('common.optionnal') }})
    v-container
      v-layout
        v-flex
          v-switch(
          v-model="value.enabled",
          @change="updateField('enabled', $event)",
          color="primary",
          hide-details
          )
        v-flex
          v-text-field.pt-0(
          type="number",
          :value="value.interval",
          :disabled="!value.enabled",
          @input="updateField('interval', $event)",
          hide-details
          )
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      default: () => ({ enabled: false, interval: '60' }),
    },
  },
  data() {
    return {
      isEnabled: false,
    };
  },
};
</script>
