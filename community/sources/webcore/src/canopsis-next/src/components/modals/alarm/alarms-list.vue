<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ title }}
    template(#text="")
      alarms-list-table-with-pagination(
        :widget="widget",
        :alarms="alarms",
        :meta="meta",
        :columns="columns",
        :loading="pending",
        :query.sync="query",
        :refresh-alarms-list="fetchList",
        selectable,
        expandable
      )
</template>

<script>
import { isEqual } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { generateDefaultAlarmListWidget } from '@/helpers/entities';
import { convertAlarmsListQueryToRequest } from '@/helpers/query';
import { alarmsListColumnsToTableColumns } from '@/helpers/forms/widgets/alarm';

import { modalInnerMixin } from '@/mixins/modal/inner';

import AlarmsListTableWithPagination from '@/components/widgets/alarm/partials/alarms-list-table-with-pagination.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.alarmsList,
  components: { AlarmsListTableWithPagination, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      pending: false,
      alarms: [],
      meta: {},
      query: {
        page: 1,
        limit: PAGINATION_LIMIT,
      },
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.alarmsList.title');
    },

    widget() {
      return this.config.widget ?? generateDefaultAlarmListWidget();
    },

    columns() {
      return alarmsListColumnsToTableColumns(this.config.widget.parameters.widgetColumns);
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
          const { data, meta } = await this.config.fetchList(convertAlarmsListQueryToRequest(this.query));

          this.alarms = data;
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
