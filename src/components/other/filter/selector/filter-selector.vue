<template lang="pug">
  v-layout
    v-flex
      v-switch(
      :label="$t('filterSelector.fields.mixFilters')",
      :value="isMultiple",
      @change="updateIsMultipleFlag"
      )
    v-flex(v-show="isMultiple")
      v-radio-group(
      :value="condition",
      @change="updateCondition"
      )
        v-radio(label="AND", value="$and")
        v-radio(label="OR", value="$or")
    v-flex
      v-select(
      v-bind="$props",
      :multiple="isMultiple",
      return-object,
      clearable,
      @input="updateFilter"
      )
</template>

<script>
import { FILTER_DEFAULT_VALUES } from '@/constants';

export default {
  props: {
    value: {
      type: [Object, Array],
      default: () => null,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'title',
    },
    itemValue: {
      type: String,
      default: 'filter',
    },
    condition: {
      type: String,
      default: FILTER_DEFAULT_VALUES.condition,
    },
  },
  computed: {
    isMultiple() {
      return Array.isArray(this.value);
    },
  },
  methods: {
    updateIsMultipleFlag(value) {
      if (value && !Array.isArray(this.value)) {
        this.updateFilter(this.value ? [this.value] : []);
      } else if (!value && Array.isArray(this.value)) {
        this.updateFilter(this.value[0] || null);
      }
    },

    updateFilter(value) {
      this.$emit('input', value);
    },

    updateCondition(value) {
      this.$emit('update:condition', value);
    },
  },
};
</script>
