<template>
  <tree-of-dependencies-preview
    class="tree-of-dependencies-expand-panel"
    :map="map"
    :columns="columns"
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
