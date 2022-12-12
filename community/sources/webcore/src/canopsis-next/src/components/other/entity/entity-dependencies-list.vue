<template lang="pug">
  entities-list-table-with-pagination(
    :widget="widget",
    :entities="entities",
    :pending="pending",
    :meta="meta",
    :query.sync="query",
    selectable
  )
    template(#toolbar="")
      v-flex
        c-advanced-search-field(
          :query.sync="query",
          :columns="columns",
          :tooltip="$t('search.contextAdvancedSearch')"
        )
      v-flex(v-if="hasAccessToCategory")
        c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
</template>

<script>
import { omit } from 'lodash';

import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query-local/query';
import { widgetColumnsContextMixin } from '@/mixins/widget/columns';
import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';
import { permissionsWidgetsContextCategory } from '@/mixins/permissions/widgets/context/category';

import EntitiesListTableWithPagination from '../../widgets/context/partials/entities-list-table-with-pagination.vue';

export default {
  components: {
    EntitiesListTableWithPagination,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    widgetColumnsContextMixin,
    entitiesEntityDependenciesMixin,
    permissionsWidgetsContextCategory,
  ],
  props: {
    entityId: {
      type: String,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    impact: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      entities: [],
      meta: {},
    };
  },
  computed: {
    headers() {
      if (this.hasColumns) {
        return [
          ...this.columns,
          { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
        ];
      }

      return [];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    updateCategory(category) {
      const categoryId = category && category._id;

      this.query = {
        ...this.query,

        page: 1,
        category: categoryId,
      };
    },

    getQuery() { // TODO: move this logic to helpers
      const query = omit(this.query, [
        'sortKey',
        'sortDir',
      ]);

      query.with_flags = true;

      if (this.query.sortKey) {
        query.sort = this.query.sortDir.toLowerCase();
        query.sort_by = this.query.sortKey;
      }

      return query;
    },

    async fetchList() {
      try {
        this.pending = true;

        const { data, meta } = await this.fetchDependenciesList({
          id: this.entityId,
          params: this.getQuery(),
        });

        this.entities = data;
        this.meta = meta;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
