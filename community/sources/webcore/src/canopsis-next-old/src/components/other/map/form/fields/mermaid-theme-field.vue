<template>
  <v-select
    v-field="value"
    v-validate="'required'"
    :items="themes"
    :error-messages="errors.collect(name)"
    :label="$t('mermaid.theme')"
    :name="name"
    class="my-0"
    hide-details
    light
  />
</template>

<script>
import { MERMAID_THEMES } from '@/constants';

export default {
  inject: ['$validator'],
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
    themes() {
      return Object.values(MERMAID_THEMES).map(value => ({
        value,
        text: this.$t(`mermaid.themes.${value}`),
      }));
    },
  },
};
</script>
