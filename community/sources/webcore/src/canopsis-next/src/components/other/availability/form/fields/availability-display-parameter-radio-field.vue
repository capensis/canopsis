<template>
  <v-radio-group
    v-field="value"
    :name="name"
    :label="label"
    hide-details
    mandatory
  >
    <v-radio
      v-for="type in types"
      :key="type.value"
      :label="type.label"
      :value="type.value"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

export default {
  props: {
    value: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    name: {
      type: String,
      default: 'parameter',
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    types() {
      return [{
        value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
        label: this.$tc('common.uptime'),
      }, {
        value: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        label: this.$t('common.downtime'),
      }];
    },
  },
};
</script>
