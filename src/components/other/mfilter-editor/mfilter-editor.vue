<template lang="pug">
  v-tabs(v-model="activeTab" slider-color="blue darken-4" centered)
    v-tab(:disabled="isSimpleTabDisabled") {{$t('mFilterEditor.tabs.visualEditor')}}
    v-tab-item
      v-container
        filter-group(
          initialGroup,
          :index = 0,
          :condition.sync="filter[0].condition",
          :possibleFields="possibleFields",
          :rules="filter[0].rules",
          :groups="filter[0].groups",
        )
    v-tab {{$t('mFilterEditor.tabs.advancedEditor')}}
    v-tab-item
      v-text-field(
          v-model="inputValue",
          rows="20",
          :label="$t('mFilterEditor.tabs.advancedEditor')",
          textarea,
          @input="handleInputChange"
        )
      v-layout(justify-center)
        v-flex(xs10 md-6)
          v-alert(:value="parseError !== ''", type="error") {{ parseError }}
      v-btn(@click="handleParseClick", :disabled="!isRequestChanged") {{$t('common.parse')}}
    v-tab(:disabled="isRequestChanged", @click='handleResultTabClick') {{$t('mFilterEditor.tabs.results')}}
    v-tab-item
      v-data-table(
        :headers='resultsTableHeaders',
        :items="Object.values(byId)",
        class="elevation-1"
        :total-items="meta.total",
        pagination.sync="pagination"
        hide-actions
      )
        template(slot="items", slot-scope="props")
          td {{props.item.connector}}
          td {{props.item.connector_name}}
          td {{props.item.component}}
          td {{props.item.resource}}
      div(class="text-xs-center pt-2")
        v-pagination(@input="handleChangePage", v-model="pagination.page", :length="paginationLength")
</template>


<script>
import { createNamespacedHelpers } from 'vuex';
import FilterGroup from '@/components/other/mfilter-editor/filter-group.vue';
import { FILTER_EDITOR_PAGINATION_LIMIT } from '@/config';

const { mapActions: mFilterActions, mapGetters: mFilterGetters } = createNamespacedHelpers('mFilterEditor');
const { mapActions: eventsActions, mapGetters: eventsGetters } = createNamespacedHelpers('events');

export default {
  name: 'mfilter-editor',

  components: {
    FilterGroup,
  },

  data() {
    return {
      paginationLength: 0,
      pagination: {},
      activeTab: 0,
      newRequest: '',
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
      isRequestChanged: false,
    };
  },

  computed: {
    ...mFilterGetters(['request', 'filter', 'possibleFields', 'parseError']),
    ...eventsGetters(['byId', 'allIds', 'meta', 'fetchComplete', 'fetchError']),

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
      },
    },

    isSimpleTabDisabled() {
      if (this.isRequestChanged || this.parseError !== '') {
        return true;
      }
      return false;
    },
  },
  watch: {
    meta(val) {
      if (val.total) {
        this.paginationLength = Math.ceil(val.total / FILTER_EDITOR_PAGINATION_LIMIT);
      }
    },
  },
  methods: {
    ...mFilterActions(['deleteParseError', 'updateFilter', 'onParseError']),
    ...eventsActions(['fetchList']),

    handleResultTabClick() {
      this.newRequest = '';
      this.fetchList({ limit: FILTER_EDITOR_PAGINATION_LIMIT, filter: JSON.stringify(this.request), start: 0 });
    },

    handleInputChange() {
      this.isRequestChanged = true;
    },

    handleParseClick() {
      this.deleteParseError();
      try {
        if (this.newRequest === '') {
          this.isRequestChanged = false;
          return this.updateFilter(JSON.parse(JSON.stringify(this.request)));
        }
        this.isRequestChanged = false;
        return this.updateFilter(JSON.parse(this.newRequest));
      } catch (e) {
        this.onParseError('Invalid JSON');
        return this.isRequestChanged = true;
      }
    },

    handleChangePage(e) {
      const start = (e - 1) * FILTER_EDITOR_PAGINATION_LIMIT;
      this.fetchList({ limit: FILTER_EDITOR_PAGINATION_LIMIT, filter: JSON.stringify(this.request), start });
    },
  },
};
</script>
