<template lang="pug">
  tree-of-dependencies-preview.tree-of-dependencies-expand-panel(:map="map")
</template>

<script>
import TreeOfDependenciesPreview from './tree-of-dependencies-preview.vue';

export default {
  components: { TreeOfDependenciesPreview },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss">
.tree-of-dependencies-expand-panel {
  height: 700px;
}
</style>
