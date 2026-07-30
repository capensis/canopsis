<template>
  <c-select
    :model-value="selectValue"
    :items="themes"
    :label="$t('mermaid.theme')"
    :name="name"
    :menu-props="menuProps"
    item-title="text"
    item-value="value"
    class="mermaid-theme-field my-0"
    hide-details
    required
    @update:model-value="updateSelectValue"
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
    selectValue() {
      return this.value || null;
    },

    menuProps() {
      return {
        attach: '.modals-wrapper',
        location: 'bottom',
        locationStrategy: 'connected',
        scrollStrategy: 'reposition',
        zIndex: 3000,
      };
    },

    themes() {
      return Object.values(MERMAID_THEMES).map((value) => ({
        value,
        text: this.$t(`mermaid.themes.${value}`),
      }));
    },
  },
  emits: ['input'],
  methods: {
    updateSelectValue(value) {
      this.updateModel(value ?? '');
    },
  },
};
</script>

<style lang="scss" scoped>
.mermaid-theme-field {
  color: rgba(0, 0, 0, 0.87);

  :deep(.v-field__input),
  :deep(.v-select__selection),
  :deep(.v-select__selection-text),
  :deep(input) {
    color: rgba(0, 0, 0, 0.87);
  }

  :deep(.v-label),
  :deep(.v-icon) {
    color: rgba(0, 0, 0, 0.6);
    opacity: 1;
  }
}
</style>
