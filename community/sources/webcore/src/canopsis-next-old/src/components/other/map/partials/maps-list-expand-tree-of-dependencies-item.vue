<template>
  <tree-of-dependencies-preview
    :map="map"
    :columns="columns"
    class="tree-of-dependencies-expand-panel"
  />
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
    columns: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

<style lang="scss">
.tree-of-dependencies-expand-panel {
  height: 700px;
}
</style>
