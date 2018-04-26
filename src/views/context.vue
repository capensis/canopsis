<template lang="pug">
  div
    context-table
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
  th {
    overflow: hidden;
    text-overflow: ellipsis;
  }
  td, th {
    padding: 1%;
  }
  td {
    overflow-wrap: break-word;
  }
  .container {
    width: 100%;
    display: flex;
    flex-wrap: nowrap;
    justify-content: space-around;
    height: 100%;
  }
  .container>td,.container>th{
    flex: 5 0 auto;
    width: 10%;
  }
  .actions{
    position: absolute;
    flex: 0;
    width: 10%;
    right: 0;
  }
</style>
