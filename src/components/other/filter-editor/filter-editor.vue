<template lang="pug">
  v-tabs(v-model="activeTab" slider-color="blue darken-4" centered)
    v-tab(:disabled="isSimpleTabDisabled") {{$t('filterEditor.tabs.visualEditor')}}
    v-tab-item
      v-container
        filter-group(
        :condition.sync="filter.condition",
        :groups.sync="filter.groups",
        :rules.sync="filter.rules",
        :possibleFields="possibleFields",
        initialGroup
        )
    v-tab {{$t('filterEditor.tabs.advancedEditor')}}
    v-tab-item
      v-text-field(
      v-model="inputValue",
      rows="10",
      :label="$t('filterEditor.tabs.advancedEditor')",
      textarea,
      )
      v-layout(justify-center)
        v-flex(xs10 md-6)
          v-alert(:value="parseError", type="error") {{ parseError }}
      v-btn(@click="parse", :disabled="!isRequestChanged") {{$t('common.parse')}}
    v-tab(@click="fetchList", :disabled="isRequestChanged") {{$t('filterEditor.tabs.results')}}
    v-tab-item
      v-data-table.elevation-1(
      :headers="resultsTableHeaders",
      :items="items",
      :pagination.sync="pagination",
      :total-items="meta.total",
      :rows-per-page-items="rowsPerPageItems",
      :loading="pending"
      )
        template(slot="items", slot-scope="props")
          td {{props.item.v.connector}}
          td {{props.item.v.connector_name}}
          td {{props.item.v.component}}
          td {{props.item.v.resource}}
</template>


<script>
import { createNamespacedHelpers } from 'vuex';
import cloneDeep from 'lodash/cloneDeep';

import { PAGINATION_PER_PAGE_VALUES } from '@/config';
import { FILTER_DEFAULT_VALUES } from '@/constants';
import parseGroupToFilter from '@/helpers/filter-editor/parse-group-to-filter';
import parseFilterToRequest from '@/helpers/filter-editor/parse-filter-to-request';
import FilterGroup from '@/components/other/filter-editor/filter-group.vue';

const { mapActions: alarmsMapActions } = createNamespacedHelpers('alarm');

/**
 * Component to create new MongoDB filter
 */
export default {
  components: {
    FilterGroup,
  },
  data() {
    return {
      items: [],
      meta: {},
      pending: false,
      pagination: {},
      activeTab: 0,
      newRequest: '',
      parseError: '',
      filter: cloneDeep(FILTER_DEFAULT_VALUES.group),
      isRequestChanged: false,
      possibleFields: ['component_name', 'connector_name', 'connector', 'resource'],
      resultsTableHeaders: [
        {
          text: this.$t('filterEditor.resultsTableHeaders.connector'),
          align: 'left',
          sortable: false,
          value: 'connector',
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.connectorName'),
          align: 'left',
          sortable: false,
          value: 'connector_name',
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.component'),
          align: 'left',
          sortable: false,
          value: 'component',
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.resource'),
          align: 'left',
          sortable: false,
          value: 'resource',
        },
      ],
    };
  },
  computed: {
    rowsPerPageItems() {
      return PAGINATION_PER_PAGE_VALUES;
    },
    request() {
      try {
        return parseFilterToRequest(this.filter);
      } catch (err) {
        return err;
      }
    },
    /**
     * @description Value of the input field of the advanced editor.
     * Prettify the value of the parsed filter
     */
    inputValue: {
      get() {
        return JSON.stringify(this.request, undefined, 4);
      },
      set(value) {
        this.newRequest = value;
        this.isRequestChanged = true;
      },
    },

    isSimpleTabDisabled() {
      return this.isRequestChanged || this.parseError !== '';
    },
  },
  watch: {
    pagination(value, oldValue) {
      if (value.page !== oldValue.page || value.rowsPerPage !== oldValue.rowsPerPage) {
        this.fetchList();
      }
    },
    activeTab() {
      this.newRequest = '';
    },
    filter: {
      handler(value) {
        this.$emit('update:filter', JSON.stringify(parseFilterToRequest(value)));
      },
      deep: true,
    },
  },
  methods: {
    ...alarmsMapActions({
      fetchAlarmListWithoutStore: 'fetchListWithoutStore',
    }),

    updateFilter(value) {
      try {
        this.filter = parseGroupToFilter(value);
      } catch (err) {
        this.parseError = err.message;
      }
    },

    async fetchList() {
      const { rowsPerPage: limit, page } = this.pagination;

      this.pending = true;
      const { alarms, total } = await this.fetchAlarmListWithoutStore({
        params: {
          limit,
          skip: limit * (page - 1),
          filter: this.request,
        },
      });

      this.pending = false;
      this.items = alarms;
      this.meta = { total };
    },

    parse() {
      this.parseError = '';

      try {
        if (this.newRequest === '') {
          this.isRequestChanged = false;
          return this.updateFilter(JSON.parse(JSON.stringify(this.request)));
        }

        this.isRequestChanged = false;

        return this.updateFilter(JSON.parse(this.newRequest));
      } catch (err) {
        this.parseError = this.$t('filterEditor.errors.invalidJSON');
        return this.isRequestChanged = true;
      }
    },
  },
};
</script>
