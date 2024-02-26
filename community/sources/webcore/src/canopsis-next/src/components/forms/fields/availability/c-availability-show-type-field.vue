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
        label: this.$tc('common.percent'),
      }, {
        value: AVAILABILITY_SHOW_TYPE.duration,
        label: this.$t('common.duration'),
      }];
    },
  },
};
</script>
