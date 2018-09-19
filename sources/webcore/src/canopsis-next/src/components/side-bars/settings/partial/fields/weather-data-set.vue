<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") Data set
    filter-editor
</template>

<script>
import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';

export default {
  components: {
    FilterEditor,
  },
};
</script>
