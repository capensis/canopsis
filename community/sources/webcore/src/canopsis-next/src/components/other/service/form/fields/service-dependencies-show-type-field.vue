<template>
  <c-radio-group
    :model-value="modelValue"
    class="mt-0"
    name="opened"
    hide-details
    mandatory
    inline
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      v-for="type in types"
      :key="type.value"
      :label="type.label"
      :value="type.value"
      color="primary"
    />
  </c-radio-group>
</template>

<script>
import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

export default {
  props: {
    modelValue: {
      type: Number,
      default: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    types() {
      return [
        TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
        TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState,
      ].map(value => ({
        value,
        label: this.$t(`entity.treeOfDependenciesShowTypes.${value}`),
      }));
    },
  },
};
</script>
