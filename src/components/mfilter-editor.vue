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
      v-btn(@click="handleParseClick", :disabled="!isRequestChanged") {{$t('common.parse')}}
      p(v-if="parseError !== ''") {{ parseError }}
    v-tab(:disabled="isRequestChanged", @click='handleResultTabClick') {{$t('mFilterEditor.tabs.results')}}
    v-tab-item
      v-data-table(
        :headers='resultsTableHeaders',
        :items="Object.values(byId)",
        hide-actions,
        class="elevation-1"
      )
        template(slot="items", slot-scope="props")
          td {{props.item.connector}}
          td {{props.item.connector_name}}
          td {{props.item.component}}
          td {{props.item.resource}}
      // div(class="text-xs-center my-2")
        v-pagination(:length="6" v-model="currentPage")
</template>


<script>
import { createNamespacedHelpers } from 'vuex';
import FilterGroup from '@/components/other/mfilter-editor/filter-group.vue';

const { mapActions: mFilterActions, mapGetters: mFilterGetters } = createNamespacedHelpers('mFilterEditor');
const { mapActions: eventsActions, mapGetters: eventsGetters } = createNamespacedHelpers('entities/events');

export default {
  name: 'mfilter-editor',

  components: {
    FilterGroup,
  },

  data() {
    return {
      activeTab: 0,
      newRequest: '',
      resultsTableHeaders: [
        {
          text: 'Connector',
          align: 'left',
          sortable: false,
          value: 'connector',
        },
        {
          text: 'Connector Name',
          align: 'left',
          sortable: false,
          value: 'connector_name',
        },
        {
          text: 'Component',
          align: 'left',
          sortable: false,
          value: 'component',
        },
        {
          text: 'Resource',
          align: 'left',
          sortable: false,
          value: 'resource',
        },
      ],
      isRequestChanged: false,
      currentPage: 1,
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
      set(newVal) {
        this.newRequest = newVal;
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
    ...eventsActions(['fetchList']),

    handleResultTabClick() {
      this.newRequest = '';
      this.currentPage = 1;
      this.fetchList({ limit: 5, filter: JSON.stringify(this.request), start: 0 });
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
        return this.isRequestChanged = true;
      }
    },
  },
};
</script>
