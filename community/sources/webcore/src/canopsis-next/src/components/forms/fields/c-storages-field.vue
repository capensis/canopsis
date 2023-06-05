<template lang="pug">
  v-layout(row, wrap)
    v-flex(v-show="label", xs12)
      v-layout
        h4.subheading.grey--text.text--darken-2 {{ label }}
        c-help-icon(v-if="helpText", :text="helpText", icon-class="ml-2 storage-help-tooltip", right)
    v-flex(xs12)
      v-layout(column)
        c-storage-field(
          v-for="(storage, index) in storages",
          :key="storage.key",
          :value="storage.directory",
          :disabled="disabled",
          @edit="$emit('edit', storage)",
          @remove="removeItemFromArray(index)"
        )
    v-flex(xs12)
      v-layout
        v-btn.ml-0(color="primary", :disabled="disabled", @click="$emit('add')") {{ $t('common.add') }}
</template>

<script>
import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'storages',
    event: 'input',
  },
  props: {
    storages: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    helpText: {
      type: String,
      required: false,
    },
  },
};
</script>

<style lang="scss">
.storage-help-tooltip {
  pointer-events: auto !important;
}
</style>
