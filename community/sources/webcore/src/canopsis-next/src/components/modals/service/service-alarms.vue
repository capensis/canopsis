<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.alarmsList.title') }}
    template(#text="")
      alarms-list-table-with-pagination(
        :widget="widget",
        :alarms="alarms",
        :meta="meta",
        :columns="columns",
        :loading="pending",
        :query.sync="query"
      )
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import AlarmsListTableWithPagination from '@/components/widgets/alarm/partials/alarms-list-table-with-pagination.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.serviceAlarms,
  components: { AlarmsListTableWithPagination, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      alarms: [],
      meta: {},
      query: {},
      pending: false,
    };
  },
  computed: {
    widget() {
      return this.config.widget;
    },

    columns() {
      return this.config.columns;
    },
  },
  methods: {
    async fetchList() {
      this.pending = true;

      if (this.config.fetchList) {
        this.alarms = await this.config.fetchList({ query: this.query });
      }

      this.pending = false;
    },
  },
};
</script>
