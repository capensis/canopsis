<template lang="pug">
  div(data-test="filterEditor")
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
          v-validate="'json'",
          :label="$t('filterEditor.tabs.advancedEditor')",
          :error-messages="errors.collect('requestString')",
          data-vv-validate-on="none",
          name="requestString",
          rows="10",
          @input="updateRequestString"
        )
        v-layout(justify-center)
          v-flex(xs10, md-6)
            v-alert(:value="parseError", type="error") {{ parseError }}
        v-btn(
          :disabled="!isRequestStringChanged || errors.has('requestString')",
          @click="parseRequestStringToFilter"
        ) {{ $t('common.parse') }}
    v-alert(:value="errors.has('filter')", type="error") {{ $t('filterEditor.errors.required') }}
</template>


<script>
import { cloneDeep, isEmpty, isString } from 'lodash';

import { ENTITIES_TYPES, FILTER_DEFAULT_VALUES } from '@/constants';

import parseGroupToFilter from '@/helpers/filter/editor/parse-group-to-filter';
import parseFilterToRequest from '@/helpers/filter/editor/parse-filter-to-request';
import { checkIfGroupIsEmpty } from '@/helpers/filter/editor/filter-check';

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
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity, ENTITIES_TYPES.pbehavior].includes(value),
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
      data.parseError = this.$t('filterEditor.errors.cantParseToVisualEditor');
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

        case ENTITIES_TYPES.pbehavior:
          return ['name', 'type', 'impact', 'depends'];

        default:
          return [];
      }
    },
  },
  created() {
    if (this.required) {
      this.$validator.extend('json', {
        getMessage: () => this.$t('filterEditor.errors.invalidJSON'),
        validate: (value) => {
          try {
            return !!JSON.parse(value);
          } catch (err) {
            return false;
          }
        },
      });

      this.$validator.attach({
        name: 'filter',
        rules: 'required:true',
        getter: () => {
          const isFilterNotEmpty = !checkIfGroupIsEmpty(this.filter);
          const isRequestStringNotEmpty = this.isRequestStringChanged && this.requestString !== '';

          return isFilterNotEmpty || isRequestStringNotEmpty;
        },
        context: () => this,
      });
    }
  },
  methods: {
    updateFilter(value) {
      const preparedFilter = parseFilterToRequest(value);

      this.filter = value;
      this.requestString = this.$options.filters.json(preparedFilter);

      this.$emit('input', isString(this.value) ? this.requestString : preparedFilter);

      if (this.required && this.errors.has('filter')) {
        this.$validator.validate('filter');
      }
    },

    updateRequestString(requestString) {
      try {
        this.errors.remove('requestString');

        if (!this.isRequestStringChanged) {
          this.isRequestStringChanged = true;
        }

        this.$emit('input', isString(this.value) ? requestString : JSON.parse(requestString));
      } catch (err) {
        console.warn(err);
      }
    },

    openAdvancedTab() {
      if (!this.isRequestStringChanged) {
        this.requestString = this.$options.filters.json(this.request);
      }
    },

    parseRequestStringToFilter() {
      try {
        this.parseError = '';
        this.errors.remove('requestString');

        if (this.requestString !== '') {
          this.updateFilter(parseGroupToFilter(this.parseRequestStringToObject()));
        } else {
          this.requestString = this.$options.filters.json(this.request);
        }

        this.isRequestStringChanged = false;
      } catch (err) {
        if (!this.errors.has('requestString')) {
          this.parseError = this.$t('filterEditor.errors.cantParseToVisualEditor');
        }
      }
    },

    parseRequestStringToObject() {
      try {
        return JSON.parse(this.requestString);
      } catch (err) {
        this.errors.add({
          field: 'requestString',
          msg: this.$t('filterEditor.errors.invalidJSON'),
          rule: 'json',
        });

        throw err;
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
