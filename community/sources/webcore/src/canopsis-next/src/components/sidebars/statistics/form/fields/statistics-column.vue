<template>
  <v-layout column>
    <kpi-rating-metric-field
      v-field="column.metric"
      :metrics="availableMetrics"
      :type="type"
      :label="$tc('common.column')"
      :name="`${name}.column`"
      required
    />
    <c-enabled-field
      v-model="customLabel"
      :label="$t('settings.columns.customLabel')"
      class="pa-0 my-2"
      @change="updateCustomLabel"
    />
    <v-text-field
      v-if="customLabel"
      v-field="column.label"
      v-validate="'required'"
      :label="$t('common.label')"
      :error-messages="errors.collect(`${name}.label`)"
      :name="`${name}.label`"
    />
    <c-enabled-field
      v-field="column.split"
      :label="$t('settings.statisticsWidgetColumn.split')"
      :disabled="!hasPossibilityToSplit"
    />
    <c-select-field
      v-if="column.split"
      v-field="column.criteria"
      :items="ratingSettings"
      :label="$t('common.infos')"
      :name="`${name}.infos`"
      item-text="label"
      item-value="id"
      required
    />
  </v-layout>
</template>

<script>
import {
  KPI_RATING_SETTINGS_TYPES,
  STATISTICS_WIDGETS_ENTITY_METRICS,
  STATISTICS_WIDGETS_USER_METRICS,
  STATISTICS_WIDGETS_USER_METRICS_WITH_ENTITY_TYPE,
} from '@/constants';

import { formMixin } from '@/mixins/form';

import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';
import KpiRatingMetricField from '@/components/other/kpi/charts/form/fields/kpi-rating-metric-field.vue';

export default {
  inject: ['$validator'],
  components: {
    CEnabledField,
    KpiRatingMetricField,
  },
  mixins: [formMixin],
  model: {
    prop: 'column',
    event: 'input',
  },
  props: {
    column: {
      type: Object,
      default: () => ({}),
    },
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
    name: {
      type: String,
      default: 'column',
    },
    ratingSettings: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      customLabel: !!this.column.label,
    };
  },
  computed: {
    hasPossibilityToSplit() {
      return STATISTICS_WIDGETS_USER_METRICS_WITH_ENTITY_TYPE.includes(this.column?.metric);
    },

    availableMetrics() {
      const metrics = this.isEntityType
        ? STATISTICS_WIDGETS_ENTITY_METRICS
        : STATISTICS_WIDGETS_USER_METRICS;

      return metrics.map((value) => {
        const kpiKey = `kpi.statisticsWidgets.metrics.${value}`;
        const alarmKey = `alarm.metrics.${value}`;

        let text;

        if (this.$te(kpiKey)) {
          text = this.$t(kpiKey);
        } else if (this.$te(alarmKey)) {
          text = this.$t(alarmKey);
        }

        return {
          value,
          text: text ?? this.$t(`user.metrics.${value}`),
        };
      });
    },
  },
  watch: {
    hasPossibilityToSplit(value) {
      if (!value) {
        this.updateModel({
          ...this.column,

          split: false,
          criteria: '',
        });
      }
    },
  },
  methods: {
    updateCustomLabel(checked) {
      if (checked) {
        return;
      }

      this.updateField('label', '');
    },
  },
};
</script>
