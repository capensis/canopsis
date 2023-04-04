<template lang="pug">
  v-autocomplete(
    v-field="value",
    v-validate="rules",
    :items="availableParameters",
    :label="label",
    :name="name",
    :loading="pending",
    :multiple="isMultiple",
    :hide-details="hideDetails"
  )
    template(#selection="{ item, index }")
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
import { createNamespacedHelpers } from 'vuex';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

export default {
  inject: ['$validator'],
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
  },
  data() {
    return {
      pending: false,
      externalMetrics: [],
    };
  },
  computed: {
    isMultiple() {
      return isArray(this.value);
    },

    isMinValueLength() {
      return this.value.length === this.min;
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
        if (value) {
          this.fetchList();
        }
      },
    },
  },
  methods: {
    ...mapMetricsActions({
      fetchExternalMetricsListWithoutStore: 'fetchExternalMetricsListWithoutStore',
    }),

    isActiveValue(value) {
      return this.isMultiple ? this.value.includes(value) : this.value === value;
    },

    getSelectionLabel(item) {
      if (this.isMinValueLength) {
        return item.text;
      }

      return this.$t('common.parametersToDisplay', { count: this.value.length });
    },

    async fetchList() {
      this.pending = true;

      try {
        const { data: externalMetrics } = await this.fetchExternalMetricsListWithoutStore({
          params: { paginate: false },
        });

        this.externalMetrics = externalMetrics;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
