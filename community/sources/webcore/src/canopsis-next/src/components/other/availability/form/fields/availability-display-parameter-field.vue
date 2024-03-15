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
import { computed } from 'vue';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

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
  setup() {
    const { tc } = useI18n();

    const types = computed(() => [{
      value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      text: tc('common.uptime'),
    }, {
      value: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
      text: tc('common.downtime'),
    }]);

    return {
      types,
    };
  },
};
</script>
