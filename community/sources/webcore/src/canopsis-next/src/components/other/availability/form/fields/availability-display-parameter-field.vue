<template>
  <v-select
    v-field="value"
    :items="types"
    :name="name"
    :label="label"
    hide-details
  />
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
        text: this.$tc('common.uptime'),
      }, {
        value: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        text: this.$t('common.downtime'),
      }];
    },
  },
};
</script>
