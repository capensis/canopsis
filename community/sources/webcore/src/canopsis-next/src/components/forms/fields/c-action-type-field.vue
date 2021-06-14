<template lang="pug">
  v-select(
    v-field="value",
    v-validate="'required'",
    :items="actionTypes",
    :error-messages="errors.collect(name)",
    :label="label || $t('common.type')",
    :name="name"
  )
</template>

<script>
import { SCENARIO_ACTION_TYPES, CAT_SCENARIO_ACTION_TYPES } from '@/constants';

import entitiesInfoMixin from '@/mixins/entities/info';

export default {
  inject: ['$validator'],
  mixins: [entitiesInfoMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
  },
  computed: {
    actionTypes() {
      return Object.values(SCENARIO_ACTION_TYPES)
        .filter(type => !CAT_SCENARIO_ACTION_TYPES.includes(type) || this.isCatVersion)
        .map(type => ({
          value: type,
          text: this.$t(`scenario.actions.${type}`),
        }));
    },
  },
};
</script>
