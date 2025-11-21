<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <alarms-list-table-with-pagination
        :widget="widget"
        :columns="widget.parameters.widgetColumns"
        :alarms="alarms"
        :meta="meta"
        :loading="pending"
        :query.sync="query"
        :refresh-alarms-list="fetchList"
        :columns-settings="columnsSettings"
        selectable
        expandable
        @update:columns-settings="updateColumnsSettings"
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
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

import AlarmsListTableWithPagination from '@/components/widgets/alarm/partials/alarms-list-table-with-pagination.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.alarmsList,
  components: { AlarmsListTableWithPagination, ModalWrapper },
  mixins: [
    modalInnerMixin,
    widgetColumnsAlarmMixin,
    entitiesUserPreferenceMixin,
  ],
  data() {
    const { config = {} } = this.modal;

    return {
      pending: false,
      alarms: [],
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

    columnsSettings() {
      return this.userPreference.content.alarms_columns_settings ?? {};
    },
  },
  watch: {
    query(query, prevQuery) {
      if (!isEqual(query, prevQuery)) {
        this.fetchList();
      }
    },
  },
  async mounted() {
    if (this.widget._id) {
      await this.fetchUserPreference({ id: this.widget._id });
    }

    this.fetchList();
  },
  methods: {
    updateColumnsSettings(columnsSettings) {
      this.updateContentInUserPreference({ alarms_columns_settings: columnsSettings });
    },

    async fetchList() {
      try {
        this.pending = true;

        if (this.config.fetchList) {
          const { data, meta } = await this.config.fetchList(convertWidgetQueryToRequest(this.query));

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
