<template lang="pug">
  div
    mfilter-editor
    filter-selector
</template>

<script>
// COMPONENTS
import MfilterEditor from '@/components/other/mfilter-editor/mfilter-editor.vue';
import FilterSelector from '@/components/other/filter/filter-selector.vue';

export default {
  components: { MfilterEditor, FilterSelector },
};
</script>
