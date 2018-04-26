<template lang="pug">
  div
    context-table
    //v-pagination(:length="nbEntitiesToDisplay" @input="changePage")
</template>

<script>
import ContextTable from '@/components/Context/context-table.vue';

export default {
  name: 'Context',
  components: { ContextTable },
};
</script>

<style scoped>
</style>
