<template lang="pug">
  div
    p {{ impacts }}
    context-list

</template>

<script>
import ContextList from '@/components/other/context-explorer/context-list.vue';

export default {
  components: { ContextList },
  data() {
    return {
      impacts: [],
    };
  },
};
</script>

