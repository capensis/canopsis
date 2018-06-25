<template lang="pug">
  div
    filter-editor
    filter-selector
</template>

<script>
// COMPONENTS
import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import FilterSelector from '@/components/other/filter/filter-selector.vue';

export default {
  components: { FilterEditor, FilterSelector },
};
</script>
