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
      v-data-table.elevation-1(
        :headers='resultsTableHeaders',
        :items="items",
        hide-actions,
      )
        template(slot="items", slot-scope="props")
          td {{props.item.v.connector}}
          td {{props.item.v.connector_name}}
          td {{props.item.v.component}}
          td {{props.item.v.resource}}
      pagination(:limit="limit", :meta="meta")
</template>


<script>
import { createNamespacedHelpers } from 'vuex';

import FilterGroup from '@/components/other/filter-editor/filter-group.vue';
import paginationMixin from '@/mixins/pagination';

const { mapActions: mFilterActions, mapGetters: mFilterGetters } = createNamespacedHelpers('mFilterEditor');
const { mapActions: alarmsMapActions, mapGetters: alarmsMapGetters } = createNamespacedHelpers('alarm');

/**
 * Component to create new MongoDB filter
 */
export default {
  name: 'mfilter-editor',

  components: {
    FilterGroup,
  },
  mixins: [paginationMixin],
  data() {
    return {
      paginationLength: 0,
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
    ...alarmsMapGetters(['items', 'meta']),

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
  methods: {
    ...mFilterActions(['deleteParseError', 'updateFilter', 'onParseError']),
    ...alarmsMapActions({ fetchListAction: 'fetchList' }),

    handleResultTabClick() {
      this.newRequest = '';
      /*this.$router.push({
        query: {
          limit: 5,
          filter: JSON.stringify(this.request),
          skip: 0,
        },
      });*/
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
  },
};
</script>
