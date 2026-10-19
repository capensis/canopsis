<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <entities-list-table-with-pagination
        :widget="widget"
        :columns="widget.parameters.widgetColumns"
        :entities="entities"
        :meta="meta"
        :loading="pending"
        :query.sync="query"
        selectable
      />
    </template>
  </modal-wrapper>
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import { convertWidgetQueryToRequest } from '@/helpers/entities/shared/query';
import { convertAlarmWidgetToQuery } from '@/helpers/entities/alarm/query';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { widgetColumnsAlarmMixin } from '@/mixins/widget/columns';

import EntitiesListTableWithPagination from '@/components/widgets/context/partials/entities-list-table-with-pagination.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entitiesList,
  components: { EntitiesListTableWithPagination, ModalWrapper },
  mixins: [
    modalInnerMixin,
    widgetColumnsAlarmMixin,
  ],
  data() {
    const { config = {} } = this.modal;

    return {
      pending: false,
      entities: [],
      meta: {},
      query: convertAlarmWidgetToQuery(config.widget),
    };
  },
  computed: {
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
