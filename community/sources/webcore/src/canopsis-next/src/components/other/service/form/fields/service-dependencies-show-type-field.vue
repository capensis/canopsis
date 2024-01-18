<template>
  <v-radio-group
    class="mt-0"
    v-field="value"
    name="opened"
    hide-details
    mandatory
    row
  >
    <v-radio
      v-for="type in types"
      :key="type.value"
      :label="type.label"
      :value="type.value"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: Number,
      default: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
    },
  },
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
