<template lang="pug">
  v-list-group(data-test="periodicRefresh")
    v-list-tile(slot="activator") {{$t('settings.periodicRefresh')}}
      .font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      v-layout
        v-flex
          v-switch(
            v-field="value.enabled",
            data-test="periodicRefreshSwitch",
            color="primary",
            hide-details
          )
        v-flex
          v-text-field.pt-0(
            v-field="value.interval",
            :disabled="!value.enabled",
            data-test="periodicRefreshField",
            type="number",
            hide-details
          )
</template>

<script>
export default {
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
