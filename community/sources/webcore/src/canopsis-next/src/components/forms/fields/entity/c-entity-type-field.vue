<template lang="pug">
  v-select(
    v-field="value",
    v-validate="'required'",
    :items="actionTypes",
    :error-messages="errors.collect(name)",
    :label="label || $t('common.type')",
    :disabled="disabled",
    :name="name",
    :multiple="isMultiple",
    :deletable-chips="isMultiple",
    :small-chips="isMultiple"
  )
</template>

<script>
import { isArray } from 'lodash';

import { BASIC_ENTITY_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Array],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    types: {
      type: Array,
      required: false,
    },
  },
  computed: {
    actionTypes() {
      const types = this.types ?? Object.values(BASIC_ENTITY_TYPES);

      return types.map(type => ({
        value: type,
        text: this.$t(`entity.types.${type}`),
      }));
    },

    isMultiple() {
      return isArray(this.value);
    },
  },
};
</script>
