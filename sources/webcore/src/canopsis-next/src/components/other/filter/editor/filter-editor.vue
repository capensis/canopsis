<template lang="pug">
  div
    v-tabs.filter-editor(v-model="activeTab" slider-color="blue darken-4" centered)
      v-tab(:disabled="isRequestStringChanged") {{ $t('filterEditor.tabs.visualEditor') }}
      v-tab-item
        v-container.pa-1
          filter-group(
          :group="filter",
          :possibleFields="possibleFields",
          isInitial,
          @update:group="updateFilter"
          )
      v-tab(@click="openAdvancedTab") {{ $t('filterEditor.tabs.advancedEditor') }}
      v-tab-item
        v-textarea(
        v-model="requestString",
        :label="$t('filterEditor.tabs.advancedEditor')",
        rows="10",
        @input="updateRequestString"
        )
        v-layout(justify-center)
          v-flex(xs10 md-6)
            v-alert(:value="parseError", type="error") {{ parseError }}
        v-btn(@click="parse", :disabled="!isRequestStringChanged") {{ $t('common.parse') }}
    v-alert(:value="errors.has('filter')", type="error") {{ $t('filterEditor.errors.required') }}
</template>


<script>
import { cloneDeep, isEmpty, isString } from 'lodash';

import { ENTITIES_TYPES, FILTER_DEFAULT_VALUES } from '@/constants';

import parseGroupToFilter from '@/helpers/filter/editor/parse-group-to-filter';
import parseFilterToRequest from '@/helpers/filter/editor/parse-filter-to-request';

import FilterGroup from './partial/filter-group.vue';
import FilterResultsAlarm from './partial/results/alarm.vue';
import FilterResultsEntity from './partial/results/entity.vue';

/**
 * Component to create new MongoDB filter
 *
 * @prop {string} value - Initial value for filter
 *
 * @event input
 */
export default {
  inject: ['$validator'],
  components: {
    FilterGroup,
    FilterResultsAlarm,
    FilterResultsEntity,
  },
  props: {
    value: {
      type: [String, Object],
      default: '',
    },
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity].includes(value),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const data = {
      filter: cloneDeep(FILTER_DEFAULT_VALUES.group),
      activeTab: 0,
      parseError: '',
      requestString: '',
      isRequestStringChanged: false,
    };

    try {
      if (this.value !== '') {
        const parsedFilter = isString(this.value) ? JSON.parse(this.value) : this.value;

        if (!isEmpty(parsedFilter)) {
          data.filter = parseGroupToFilter(parsedFilter);
        }
      }
    } catch (err) {
      data.activeTab = 1;
      data.requestString = isString(this.value) ? this.value : JSON.stringify(this.value);
      data.isRequestStringChanged = true;
      data.parseError = this.$t('filterEditor.errors.invalidJSON');
    }

    return data;
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

    possibleFields() {
      switch (this.entitiesType) {
        case ENTITIES_TYPES.alarm:
          return ['connector', 'connector_name', 'component', 'resource'];

        case ENTITIES_TYPES.entity:
          return ['name', 'type'];

        default:
          return [];
      }
    },
  },
  created() {
    if (this.required && this.$validator) {
      this.$validator.attach('filter', 'required:true', {
        getter: () => {
          const firstRule = Object.values(this.filter.rules)[0];

          return firstRule && firstRule.field !== '' && firstRule.operator !== '' && firstRule.input !== '';
        },
        context: () => this,
      });
    }
  },
  methods: {
    updateFilter(value) {
      const preparedFilter = parseFilterToRequest(value);

      this.filter = value;

      this.$emit('input', isString(this.value) ? JSON.stringify(preparedFilter) : preparedFilter);

      if (this.required && this.$validator && this.errors.has('filter')) {
        this.$validator.validate('filter');
      }
    },

    updateRequestString(requestString) {
      if (!this.isRequestStringChanged) {
        this.isRequestStringChanged = true;
      }

      this.$emit('input', requestString);
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
        } else {
          this.requestString = JSON.stringify(this.request, undefined, 4);
        }

        this.isRequestStringChanged = false;
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
