<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.filterOnOpenResolved') }}
    v-container
      v-layout
        v-radio-group.mt-0(v-field="value", name="opened", hide-details, mandatory)
          v-radio(
            v-for="type in types",
            :key="type.value",
            :label="type.label",
            :value="type.value",
            color="primary"
          )
</template>

<script>
import { ALARMS_OPENED_VALUES } from '@/constants';

export default {
  props: {
    value: {
      type: Boolean,
      default: ALARMS_OPENED_VALUES.opened,
    },
  },
  computed: {
    types() {
      return Object.values(ALARMS_OPENED_VALUES).map(value => ({
        value,
        label: this.$t(`settings.openedTypes.${value}`),
      }));
    },
  },
};
</script>
