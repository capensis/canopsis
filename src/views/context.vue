<template lang="pug">
  context-table(
    :contextProperty="$mq| mq("
                      +"{mobile: contextPropertyOnMobile,tablet: contextPropertyOnTablet,laptop: contextProperty}"
                      + ")",
  )
</template>

<script>
import ContextTable from '@/components/context-table.vue';

export default {
  name: 'context',
  components: { ContextTable },
  data() {
    return {
      contextProperty: [
        {
          text: 'name',
          value: 'name',
        },
        {
          text: 'type',
          value: 'type',
        },
      ],
      contextPropertyOnMobile: [
        {
          text: 'name',
          value: 'name',
        },
      ],
      contextPropertyOnTablet: [
        {
          text: 'name',
          value: 'name',
        },
      ],
    };
  },
};
</script>
