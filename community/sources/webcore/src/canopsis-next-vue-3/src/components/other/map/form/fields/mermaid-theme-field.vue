<template>
  <c-select
    :model-value="value"
    :items="themes"
    :label="$t('mermaid.theme')"
    :name="name"
    item-title="text"
    item-value="value"
    class="my-0"
    hide-details
    required
    @update:model-value="updateModel"
  />
</template>

<script>
import { MERMAID_THEMES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
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
      return Object.values(MERMAID_THEMES).map((value) => ({
        value,
        text: this.$t(`mermaid.themes.${value}`),
      }));
    },
  },
  emits: ['input'],
};
</script>
