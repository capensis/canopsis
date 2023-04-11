<template lang="pug">
  component(
    v-validate="rules",
    :is="addable ? 'v-combobox' : 'v-autocomplete'",
    :value="value",
    :search-input.sync="searchInput",
    :items="availableParameters",
    :label="label",
    :name="name",
    :loading="externalMetricsPending",
    :multiple="isMultiple",
    :hide-details="hideDetails",
    :filter="filterComputedMetric",
    :return-object="false",
    :hide-no-data="addable",
    @change="updateParameters"
  )
    template(v-if="!addable", #selection="{ item, index }")
      template(v-if="isMultiple")
        span(v-if="!index") {{ getSelectionLabel(item) }}
      template(v-else)
        v-icon.mr-2(v-if="item.isExternal") language
        span {{ item.text }}
    template(#item="{ parent, item, tile }")
      v-list-tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-action(v-if="isMultiple")
          v-checkbox(
            :input-value="tile.props.value",
            :color="parent.color",
            :disabled="tile.props.disabled"
          )
        v-icon.mr-3(v-if="item.isExternal") language
        v-list-tile-content
          v-list-tile-title {{ item.text }}
</template>

<script>
import { isArray, omit } from 'lodash';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

import { formBaseMixin } from '@/mixins/form';
import { entitiesMetricsMixin } from '@/mixins/entities/metrics';

export default {
  inject: ['$validator'],
  mixins: [formBaseMixin, entitiesMetricsMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Array],
      required: true,
    },
    name: {
      type: String,
      default: 'parameters',
    },
    label: {
      type: String,
      required: false,
    },
    min: {
      type: Number,
      default: 1,
    },
    required: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    parameters: {
      type: Array,
      default: () => Object.values(omit(ALARM_METRIC_PARAMETERS, ['timeToAck', 'timeToResolve'])),
    },
    disabledParameters: {
      type: Array,
      default: () => [],
    },
    withExternal: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      searchInput: null,
    };
  },
  computed: {
    isMultiple() {
      return isArray(this.value);
    },

    isMinValueLength() {
      return this.value?.length === this.min;
    },

    availableParameters() {
      const parameters = this.parameters.map(value => ({
        value,
        disabled: (this.disabledParameters.includes(value) && !this.isActiveValue(value))
          || (this.isMinValueLength && this.isActiveValue(value)),
        text: this.$t(`alarm.metrics.${value}`),
      }));

      if (this.withExternal) {
        parameters.push(...this.externalMetrics.map(({ _id: value, name }) => ({
          text: name,
          value,
          isExternal: true,
        })));
      }

      return parameters;
    },

    rules() {
      return {
        required: this.required,
      };
    },
  },
  watch: {
    withExternal: {
      immediate: true,
      handler(value) {
        if (value && !this.externalMetricsPending && !this.externalMetrics.length) {
          this.fetchList();
        }
      },
    },
  },
  methods: {
    filterComputedMetric({ text }) {
      try {
        return text.includes(this.searchInput) || new RegExp(`${this.searchInput}`).test(text);
      } catch (err) {
        return false;
      }
    },

    updateParameters(value) {
      this.updateModel(value ?? '');
    },

    isActiveValue(value) {
      return this.isMultiple ? this.value.includes(value) : this.value === value;
    },

    getSelectionLabel(item) {
      if (this.isMinValueLength) {
        return item.text;
      }

      return this.$t('common.parametersToDisplay', { count: this.value.length });
    },

    fetchList() {
      this.fetchExternalMetricsList({
        params: { paginate: false },
      });
    },
  },
};
</script>
