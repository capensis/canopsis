<template>
  <widget-settings-item :title="$tc('common.column', 2)">
    <c-alert
      :value="!columns.length"
      :type="errors.has(name) ? 'error' : 'info'"
    >
      {{ $t('widgetTemplate.errors.columnsRequired') }}
    </c-alert>
    <c-progress-overlay
      :pending="pending"
      transition
    />
    <c-movable-card-iterator-field
      v-field="columns"
      addable
      @add="add"
    >
      <template #item="{ item, index }">
        <statistics-column
          v-field="columns[index]"
          :type="type"
          :name="item.key"
          :rating-settings="ratingSettings"
        />
      </template>
    </c-movable-card-iterator-field>
  </widget-settings-item>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { KPI_RATING_SETTINGS_TYPES } from '@/constants';

import { statisticsWidgetColumnToForm } from '@/helpers/entities/widget/forms/statistics';

import { formArrayMixin } from '@/mixins/form';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

import StatisticsColumn from './statistics-column.vue';

const { mapActions } = createNamespacedHelpers('ratingSettings');

export default {
  inject: ['$validator'],
  components: {
    WidgetSettingsItem,
    StatisticsColumn,
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
          /**
           * We don't have rating settings for KPI_RATING_SETTINGS_TYPES.user type with `main` === false
           */
          type: KPI_RATING_SETTINGS_TYPES.entity,
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
