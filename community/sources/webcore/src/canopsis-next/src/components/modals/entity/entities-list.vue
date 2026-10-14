<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ title }}</span>
    </template>
    <template #text>
      <entities-list-table-with-pagination
        :widget="widget"
        :columns="widget.parameters.widgetColumns"
        :entities="entities"
        :meta="meta"
        :pending="pending"
        :query="query"
        selectable
        @update:query="updateQuery"
      />
    </template>
  </c-modal-wrapper>
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import { convertContextWidgetToQuery } from '@/helpers/entities/entity/query';
import { convertWidgetQueryToRequest } from '@/helpers/entities/shared/query';

import EntitiesListTableWithPagination from '@/components/widgets/context/partials/entities-list-table-with-pagination.vue';

export default {
  name: MODALS.entitiesList,

  components: {
    EntitiesListTableWithPagination,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    const { config = {} } = this.modal;

    return {
      pending: false,
      entities: [],
      meta: {},
      query: convertContextWidgetToQuery(config.widget),
    };
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },

    title() {
      return this.config.title ?? this.$t('modals.alarmsList.title');
    },

    widget() {
      return this.config.widget;
    },
  },

  watch: {
    query(query, prevQuery) {
      if (!isEqual(query, prevQuery)) {
        this.fetchList();
      }
    },
  },

  mounted() {
    this.fetchList();
  },

  methods: {
    updateQuery(query) {
      this.query = query;
    },

    async fetchList() {
      try {
        this.pending = true;

        if (this.config.fetchList) {
          const { data, meta } = await this.config.fetchList(convertWidgetQueryToRequest(this.query));

          this.entities = data;
          this.meta = meta;
        }
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
