<template lang="pug">
  v-tabs.filter-editor(v-model="activeTab" slider-color="blue darken-4" centered)
    v-tab(:disabled="isRequestStringChanged") {{$t('filterEditor.tabs.visualEditor')}}
    v-tab-item
      v-container
        filter-group(
        :group="filter",
        :possibleFields="possibleFields",
        @update:group="updateFilter",
        isInitial
        )
    v-tab(@click="openAdvancedTab") {{$t('filterEditor.tabs.advancedEditor')}}
    v-tab-item
      v-textarea(
      v-model="requestString",
      :label="$t('filterEditor.tabs.advancedEditor')",
      @input="updateRequestString",
      rows="10",
      )
      v-layout(justify-center)
        v-flex(xs10 md-6)
          v-alert(:value="parseError", type="error") {{ parseError }}
      v-btn(@click="parse", :disabled="!isRequestStringChanged") {{$t('common.parse')}}
</template>


<script>
import cloneDeep from 'lodash/cloneDeep';
import isEmpty from 'lodash/isEmpty';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import parseGroupToFilter from '@/helpers/filter-editor/parse-group-to-filter';
import parseFilterToRequest from '@/helpers/filter-editor/parse-filter-to-request';

import FilterGroup from '@/components/other/filter-editor/partial/filter-group.vue';
import FilterResults from '@/components/other/filter-editor/partial/filter-results.vue';

/**
 * Component to create new MongoDB filter
 *
 * @prop {string} value - Initial value for filter
 *
 * @event input
 */
export default {
  components: {
    FilterGroup,
    FilterResults,
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  data() {
    let filter = cloneDeep(FILTER_DEFAULT_VALUES.group);

    if (this.value !== '') {
      const valueObject = JSON.parse(this.value);

      if (!isEmpty(valueObject)) {
        try {
          filter = parseGroupToFilter(valueObject);
        } catch (err) {
          console.warn(err);
        }
      }
    }

    return {
      filter,
      activeTab: 0,
      requestString: '',
      parseError: '',
      isRequestStringChanged: false,
      possibleFields: ['connector', 'connector_name', 'component', 'resource'],
    };
  },
  computed: {
    request() {
      try {
        return parseFilterToRequest(this.filter);
      } catch (err) {
        console.error(err);

        return {};
      }
    },
  },
  methods: {
    updateFilter(value) {
      this.filter = value;
      this.$emit('input', JSON.stringify(parseFilterToRequest(value)));
    },

    updateRequestString() {
      this.isRequestStringChanged = true;
    },

    openAdvancedTab() {
      if (!this.isRequestStringChanged) {
        this.requestString = JSON.stringify(this.request, undefined, 4);
      }
    },

    parse() {
      this.parseError = '';

      try {
        if (this.requestString !== '') {
          this.updateFilter(parseGroupToFilter(JSON.parse(this.requestString)));
          this.isRequestStringChanged = false;
        } else {
          this.requestString = JSON.stringify(this.request, undefined, 4);
          this.isRequestStringChanged = false;
        }
      } catch (err) {
        this.parseError = this.$t('filterEditor.errors.invalidJSON');
      }
    },
  },
};
</script>

<style lang="scss">
  .filter-editor {
    .v-card {
      box-shadow: 0 0 0 -1px rgba(0, 0, 0, 0.5), 0 1px 5px 0 rgba(0, 0, 0, 0.44), 0 1px 3px 0 rgba(0, 0, 0, 0.42);
    }
  }
</style>
