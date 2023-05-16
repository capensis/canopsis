<template lang="pug">
  widget-settings-item(:title="$tc('common.column', 2)")
    c-alert(:value="!columns.length", color="info") {{ $t('widgetTemplate.errors.columnsRequired') }}
    c-progress-overlay(:pending="pending", transition)
    c-movable-card-iterator-field(v-field="columns", addable, @add="add")
      template(#item="{ item, index }")
        v-layout(column)
          kpi-rating-metric-field(
            v-field="columns[index].column",
            :type="type",
            :label="$tc('common.column')",
            :name="`column-${item.key}.column`",
            required
          )
          c-enabled-field(
            v-field="columns[index].split",
            :label="$t('settings.statisticsWidgetColumn.split')"
          )
          c-select-field(
            v-if="columns[index].split",
            v-field="columns[index].infos",
            :items="ratingSettings",
            :label="$t('common.infos')",
            :name="`column-${item.key}.infos`",
            item-text="label",
            item-value="id",
            required
          )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { KPI_RATING_SETTINGS_TYPES } from '@/constants';

import { statisticsWidgetColumnToForm } from '@/helpers/forms/widgets/statistics';

import { formArrayMixin } from '@/mixins/form';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';
import KpiRatingMetricField from '@/components/other/kpi/charts/form/fields/kpi-rating-metric-field.vue';

const { mapActions } = createNamespacedHelpers('ratingSettings');

export default {
  inject: ['$validator'],
  components: {
    WidgetSettingsItem,
    CEnabledField,
    KpiRatingMetricField,
  },
  mixins: [formArrayMixin, validationAttachRequiredMixin],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    columns: {
      type: Array,
      default: () => [],
    },
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
    name: {
      type: String,
      default: 'columns',
    },
  },
  data() {
    return {
      ratingSettings: [],
      pending: false,
    };
  },
  watch: {
    columns() {
      this.validateRequiredRule();
    },
  },
  mounted() {
    this.fetchList();
    this.attachRequiredRule(() => this.columns.length > 0);
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
  methods: {
    ...mapActions({ fetchRatingSettingsWithoutStore: 'fetchListWithoutStore' }),

    add() {
      this.addItemIntoArray(statisticsWidgetColumnToForm());
    },

    async fetchList() {
      this.pending = true;

      const { data: ratingSettings = [] } = await this.fetchRatingSettingsWithoutStore({
        params: {
          type: this.type,
          main: false,
          paginate: false,
        },
      });

      this.ratingSettings = ratingSettings;
      this.pending = false;
    },
  },
};
</script>
