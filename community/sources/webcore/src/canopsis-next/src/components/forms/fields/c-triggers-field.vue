<template lang="pug">
  v-select(
    v-validate="'required'",
    v-field="value",
    :items="availableTriggers",
    :disabled="disabled",
    :label="label || $t('common.triggers')",
    :error-messages="errors.collect(name)",
    :name="name",
    multiple,
    chips
  )
</template>

<script>
import { SCENARIO_TRIGGERS, CAT_SCENARIO_TRIGGERS } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';

export default {
  inject: ['$validator'],
  mixins: [entitiesInfoMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'triggers',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableTriggers() {
      return Object.values(SCENARIO_TRIGGERS)
        .filter(type => !CAT_SCENARIO_TRIGGERS.includes(type) || this.isCatVersion);
    },
  },
};
</script>
