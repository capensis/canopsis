<template lang="pug">
  context-table(
  :contextProperties="$mq| mq(contextProperties)",
  )
</template>

<script>
import ContextEntityMixin from '@/mixins/context';
import ContextTable from '@/components/other/context-explorer/context-table.vue';
import { PAGINATION_LIMIT } from '@/config';

export default {
  name: 'context',
  components: { ContextTable },
  mixins: [
    ContextEntityMixin,
  ],
  data() {
    return {
      contextProperties: {
        laptop: [
          {
            text: this.$t('tables.contextEntities.columns._id'),
            value: '_id',
          },
          {
            text: this.$t('tables.contextEntities.columns.type'),
            value: 'type',
          },
          {
            text: this.$t('tables.contextEntities.columns.name'),
            value: 'name',
          },
        ],
        mobile: [
          {
            text: this.$t('tables.contextEntities.columns.name'),
            value: 'name',
          },
        ],
        tablet: [
          {
            text: this.$t('tables.contextEntities.columns.name'),
            value: 'name',
          },
        ],
      },
    };
  },
  mounted() {
    this.fetchContextEntities({
      params: {
        limit: PAGINATION_LIMIT,
      },
    });
  },
};
</script>
