<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ title }}</span>
    </template>
    <template #text>
      <alarms-list-table-with-pagination
        :widget="widget"
        :columns="widget.parameters.widgetColumns"
        :alarms="alarms"
        :meta="meta"
        :loading="pending"
        :query="query"
        :refresh-alarms-list="fetchList"
        :columns-settings="columnsSettings"
        selectable
        expandable
        @update:query="updateQuery"
        @update:columns-settings="updateColumnsSettings"
      />
    </template>
  </c-modal-wrapper>
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import { convertAlarmWidgetToQuery } from '@/helpers/entities/alarm/query';
import { convertWidgetQueryToRequest } from '@/helpers/entities/shared/query';

import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

import AlarmsListTableWithPagination from '@/components/widgets/alarm/partials/alarms-list-table-with-pagination.vue';

export default {
  name: MODALS.alarmsList,

  components: {
    AlarmsListTableWithPagination,
  },

  mixins: [
    entitiesUserPreferenceMixin,
  ],

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
      alarms: [],
      meta: {},
      query: convertAlarmWidgetToQuery(config.widget),
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

    columnsSettings() {
      return this.userPreference.content?.alarms_columns_settings ?? {};
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
    updateQuery(query) {
      this.query = query;
    },

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
