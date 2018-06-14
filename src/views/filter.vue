<template lang="pug">
  div
    mfilter-editor
    filter-selector
</template>

<script>
import MfilterEditor from '@/components/other/mfilter-editor/mfilter-editor.vue';
import FilterSelector from '@/components/other/filter/filter-selector.vue';

export default {
  name: 'filter-editor',
  components: { MfilterEditor, FilterSelector },
};
</script>
