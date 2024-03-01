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
import { AVAILABILITY_SHOW_TYPE } from '@/constants';

export default {
  props: {
    value: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    name: {
      type: String,
      default: 'availability',
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    types() {
      return [{
        value: AVAILABILITY_SHOW_TYPE.percent,
        text: this.$tc('common.percent'),
      }, {
        value: AVAILABILITY_SHOW_TYPE.duration,
        text: this.$t('common.duration'),
      }];
    },
  },
};
</script>
