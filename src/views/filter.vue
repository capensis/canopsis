<template lang="pug">
  div
    mfilter-editor
    selector
</template>

<script>
import MfilterEditor from '@/components/mfilter-editor.vue';
import Selector from '@/components/other/filter/selector.vue';

export default {
  name: 'filter-editor',
  components: { MfilterEditor, Selector },
};
</script>
