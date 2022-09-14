<template lang="pug">
  tree-of-dependencies-preview.tree-of-dependencies-expand-panel(
    :map="map"
  )
</template>

<script>
import { TREE_OF_DEPENDENCIES_TYPES } from '@/constants';

import TreeOfDependenciesPreview from './tree-of-dependencies-preview.vue';

export default {
  components: { TreeOfDependenciesPreview },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
  computed: {
    impact() {
      return this.map?.parameters.type === TREE_OF_DEPENDENCIES_TYPES.impactChain;
    },

    entities() {
      return (this.map?.parameters?.entities ?? [])
        .map(entity => ({ data: entity.entity, pinned: entity.pinned_entities }));
    },
  },
};
</script>

<style lang="scss">
.tree-of-dependencies-expand-panel {
  height: 700px;
}
</style>
