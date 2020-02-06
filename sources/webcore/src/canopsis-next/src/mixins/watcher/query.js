import { PAGINATION_LIMIT } from '@/config';
import queryMixin from '@/mixins/query';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

/**
 * @mixin Add query logic
 */
export default {
  mixins: [queryMixin, entitiesWatcherEntityMixin],
  data() {
    const { itemsPerPage = PAGINATION_LIMIT } = this.modal.config;

    return {
      query: {
        page: 1,
        limit: itemsPerPage,
      },
    };
  },
  computed: {
    metaData() {
      return {
        page: this.query.page,
        limit: this.query.limit,
        first: (this.query.page - 1) * this.query.limit,
        last: this.query.page * this.query.limit,
        total: this.watcherEntities.length,
      };
    },
    watchers() {
      const { first, last } = this.metaData;

      return this.watcherEntities.slice(first, last);
    },
  },
  mounted() {
    this.fetchWatchersList();
  },
  methods: {
    fetchWatchersList() {
      this.fetchWatcherEntitiesList({ watcherId: this.watcher.entity_id });
    },
    changePage(page) {
      this.query.page = page;
    },
    changeLimit(limit) {
      this.query = {
        page: 1,
        limit,
      };
    },
  },
};
