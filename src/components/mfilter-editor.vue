<template lang="pug">
  v-container
    v-layout(row, wrap)
      v-flex(xs12)
        v-tabs
          v-tab(
            @click="handleTabClick(0)"
            :disabled="parseError == '' ? false : true"
          ) {{$t('m_filter_editor.tabs.visual_editor')}}
          v-tab(
            @click="handleTabClick(1)"
          ) {{$t('m_filter_editor.tabs.advanced_editor')}}
          v-tab(
            @click="handleResultTabClick"
          ) Resultats

    v-container(v-if="activeTab === 0")
      filter-group(
        initialGroup,
        :index = 0,
        :condition.sync="filter[0].condition",
        :possibleFields="possibleFields",
        :rules="filter[0].rules",
        :groups="filter[0].groups",
      )

    v-container(v-if="activeTab === 1")
      v-text-field(
        ref="input",
        v-model="inputValue",
        rows="20",
        :label="$t('m_filter_editor.tabs.advanced_editor')",
        textarea
      )
      v-btn(@click="handleParseClick") {{$t('common.parse')}}
      p(v-if="parseError !== ''") {{ parseError }}

    v-container(v-if="activeTab === 2")
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
    };
  },

  computed: {
    ...mFilterGetters(['request', 'filter', 'possibleFields', 'activeTab', 'parseError']),
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
  },
  methods: {
    ...mFilterActions(['changeActiveTab', 'deleteParseError', 'updateFilter', 'onParseError']),
    ...eventsActions(['fetchList']),

    handleTabClick(tab) {
      this.newRequest = '';
      this.changeActiveTab(tab);
    },

    handleResultTabClick() {
      this.fetchList({ start: 0, filter: JSON.stringify(this.request), limit: 10 });
      this.changeActiveTab(2);
    },

    handleParseClick() {
      this.deleteParseError();
      try {
        if (this.newRequest === '') {
          this.updateFilter(JSON.parse(JSON.stringify(this.request)));
          return this;
        }
        this.updateFilter(JSON.parse(this.newRequest));
        return this;
      } catch (e) {
        this.onParseError(e.message);
        return e;
      }
    },
  },
};
</script>
