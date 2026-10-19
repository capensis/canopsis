<template>
  <c-advanced-data-table
    :headers="headers"
    :items="entities"
    :loading="pending"
    :options.sync="options"
    :total-items="meta.total_count"
    item-key="_id"
    search
  />
</template>

<script>
import { localQueryMixin } from '@/mixins/query/query';
import entitiesPbehaviorEntitiesMixin from '@/mixins/entities/pbehavior/entities';

export default {
  mixins: [localQueryMixin, entitiesPbehaviorEntitiesMixin],
  props: {
    pbehavior: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      pending: false,
      meta: {},
      entities: [],
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.id'),
          sortable: false,
          value: '_id',
        },
        {
          text: this.$t('common.name'),
          sortable: false,
          value: 'name',
        },
        {
          text: this.$t('common.type'),
          sortable: false,
          value: 'type',
        },
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const { data: entities, meta } = await this.fetchPbehaviorEntitiesListWithoutStore({
        id: this.pbehavior._id,
        params: this.getQuery(),
      });

      this.meta = meta;
      this.entities = entities;

      this.pending = false;
    },
  },
};
</script>
