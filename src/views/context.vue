<template lang="pug">
  context-table(
  :contextProperties="$mq| mq(contextProperties)",
  )
</template>

<script>
import ContextTable from '@/components/tables/context-table/context-table.vue';

export default {
  name: 'context',
  components: { ContextTable },
  data() {
    return {
      contextProperties: {
        laptop: [
          {
            text: 'name',
            value: 'name',
          },
          {
            text: 'type',
            value: 'type',
          },
        ],
        mobile: [
          {
            text: 'name',
            value: 'name',
          },
        ],
        tablet: [
          {
            text: 'name',
            value: 'name',
          },
        ],
      },
    };
  },
};
</script>
