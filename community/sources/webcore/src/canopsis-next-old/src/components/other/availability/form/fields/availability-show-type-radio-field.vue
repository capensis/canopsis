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
import { computed } from 'vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import { useI18n } from '@/hooks/i18n';

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
  setup() {
    const { t, tc } = useI18n();

    const types = computed(() => [{
      value: AVAILABILITY_SHOW_TYPE.percent,
      label: tc('common.percent'),
    }, {
      value: AVAILABILITY_SHOW_TYPE.duration,
      label: t('common.duration'),
    }]);

    return {
      types,
    };
  },
};
</script>
