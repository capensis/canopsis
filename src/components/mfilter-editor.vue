<template lang="pug">
  v-stepper(v-model="e1" non-linear @input="handleResultTabClick")
    v-stepper-header
      v-stepper-step(step="1" editable) {{$t('m_filter_editor.tabs.visual_editor')}}
      v-divider
      v-stepper-step(step="2" editable) {{$t('m_filter_editor.tabs.advanced_editor')}}
      v-divider
      v-stepper-step(step="3" editable) {{$t('m_filter_editor.tabs.results')}}
      v-divider
    v-stepper-items
      v-stepper-content(step="1")
        filter-group(
          initialGroup,
          :index = 0,
          :condition.sync="filter[0].condition",
          :possibleFields="possibleFields",
          :rules="filter[0].rules",
          :groups="filter[0].groups",
        )
      v-stepper-content(step="2")
        v-text-field(
          ref="input",
          v-model="inputValue",
          rows="20",
          :label="$t('m_filter_editor.tabs.advanced_editor')",
          textarea
        )
        v-btn(@click="handleParseClick") {{$t('common.parse')}}
        p(v-if="parseError !== ''") {{ parseError }}
      v-stepper-content(step="3")
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
      e1: 0,
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
  },
  methods: {
    ...mFilterActions(['deleteParseError', 'updateFilter', 'onParseError']),
    ...eventsActions(['fetchList']),

    handleResultTabClick() {
      this.newRequest = '';
      this.fetchList({ start: 0, filter: JSON.stringify(this.request), limit: 10 });
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
