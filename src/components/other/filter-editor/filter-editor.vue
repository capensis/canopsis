<template lang="pug">
  v-tabs(v-model="activeTab" slider-color="blue darken-4" centered)
    v-tab(:disabled="isSimpleTabDisabled") {{$t('mFilterEditor.tabs.visualEditor')}}
    v-tab-item
      v-container
        filter-group(
        :condition.sync="filter[0].condition",
        :groups="filter[0].groups",
        :rules="filter[0].rules",
        :possibleFields="possibleFields",
        initialGroup
        )
    v-tab {{$t('mFilterEditor.tabs.advancedEditor')}}
    v-tab-item
      v-text-field(
      v-model="inputValue",
      rows="10",
      :label="$t('mFilterEditor.tabs.advancedEditor')",
      textarea,
      )
      v-layout(justify-center)
        v-flex(xs10 md-6)
          v-alert(:value="parseError", type="error") {{ parseError }}
      v-btn(@click="parse", :disabled="!isRequestChanged") {{$t('common.parse')}}
    v-tab(@click="fetchList", :disabled="isRequestChanged") {{$t('mFilterEditor.tabs.results')}}
    v-tab-item
      v-data-table.elevation-1(
      :headers="resultsTableHeaders",
      :items="items",
      :pagination.sync="pagination",
      :total-items="meta.total",
      :rows-per-page-items="[5, 10, 25, 50]",
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

import parseGroupToFilter from '@/services/mfilter-editor/parseRequestToFilter';
import parseFilterToRequest from '@/services/mfilter-editor/parseFilterToRequest';
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
      filter: [{
        condition: '$or',
        groups: [],
        rules: [],
      }],
      isRequestChanged: false,
      possibleFields: ['component_name', 'connector_name', 'connector', 'resource'],
      resultsTableHeaders: [
        {
          text: this.$t('mFilterEditor.resultsTableHeaders.connector'),
          align: 'left',
          sortable: false,
          value: 'connector',
        },
        {
          text: this.$t('mFilterEditor.resultsTableHeaders.connectorName'),
          align: 'left',
          sortable: false,
          value: 'connector_name',
        },
        {
          text: this.$t('mFilterEditor.resultsTableHeaders.component'),
          align: 'left',
          sortable: false,
          value: 'component',
        },
        {
          text: this.$t('mFilterEditor.resultsTableHeaders.resource'),
          align: 'left',
          sortable: false,
          value: 'resource',
        },
      ],
    };
  },
  computed: {
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
        this.filter = [parseGroupToFilter(value)];
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
        this.parseError = 'Invalid JSON'; // TODO: translate
        return this.isRequestChanged = true;
      }
    },
  },
};
</script>
