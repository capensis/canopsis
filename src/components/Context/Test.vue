<template>
  <basic-list :items="items">
      <h1 slot="row">BON</h1>
      <h2 slot="expanded">JOUR</h2>
  </basic-list>
</template>

<script>
import BasicList from '../BasicComponent/BasicList.vue';
import ListBrick from '../BasicComponent/ListBrick';

export default {
  name: 'Test',
  components: { ListBrick, BasicList },
  data() {
    return {
      items: ['a', 'b', 'c'],
    };
  },
};
</script>

<style scoped>

</style>
