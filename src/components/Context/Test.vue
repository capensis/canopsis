<template lang="pug">
  basic-list(:items="items")
      h1(slot="row" slot-scope="item") {{item.props}}
      p(slot="expandedRow" slot-scope="item") Bonjour
</template>

<script>
import BasicList from '../BasicComponent/BasicList.vue';

export default {
  name: 'Test',
  components: { BasicList },
  data() {
    return {
      items: ['a', 'b', 'c'],
    };
  },
};
</script>

<style scoped>

</style>
