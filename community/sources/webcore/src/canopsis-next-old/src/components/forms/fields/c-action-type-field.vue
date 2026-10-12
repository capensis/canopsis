<template>
  <v-select
    v-field="value"
    v-validate="'required'"
    :items="actionTypes"
    :error-messages="errors.collect(name)"
    :label="label || $t('common.type')"
    :name="name"
  />
</template>

<script>
import { ACTION_TYPES, PRO_ACTION_TYPES } from '@/constants';

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
    types: {
      type: Array,
      required: false,
    },
  },
  computed: {
    actionTypes() {
      const types = this.types || Object.values(ACTION_TYPES);

      return types
        .filter(type => !PRO_ACTION_TYPES.includes(type) || this.isProVersion)
        .map(type => ({
          value: type,
          text: this.$t(`scenario.actions.${type}`),
        }));
    },
  },
};
</script>
