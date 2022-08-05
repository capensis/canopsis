<template lang="pug">
  v-select(
    v-field="value",
    v-validate="'required'",
    :items="themes",
    :error-messages="errors.collect(name)",
    :label="$t('mermaid.theme')",
    :name="name",
    hide-details
  )
</template>

<script>
import { MERMAID_THEMES } from '@/constants';

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
    themes() {
      return Object.values(MERMAID_THEMES).map(value => ({
        value,
        text: this.$t(`mermaid.themes.${value}`),
      }));
    },
  },
};
</script>
