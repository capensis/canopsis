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

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import { useI18n } from '@/i18n';

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
      text: tc('common.percent'),
    }, {
      value: AVAILABILITY_SHOW_TYPE.duration,
      text: t('common.duration'),
    }]);

    return {
      types,
    };
  },
};
</script>
