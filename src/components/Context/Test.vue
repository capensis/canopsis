<template lang="pug">
  basic-list(:items="items")
    template(slot="header")
      div B
      div A
    template(slot="row" scope="item")
      div {{ item.props.a }}
      div {{ item.props.b }}
    template(slot="expandedRow" scope="item")
      div {{ item.props.a }}
      div {{ item.props.b }}
</template>

<script>
import BasicList from '../BasicComponent/BasicList.vue';

export default {
  name: 'Test',
  components: { BasicList },
  data() {
    return {
      items: [{ a: 1, b: 2 }, { a: 1, b: 2}],
    };
  },
};
</script>

<style scoped>

</style>
