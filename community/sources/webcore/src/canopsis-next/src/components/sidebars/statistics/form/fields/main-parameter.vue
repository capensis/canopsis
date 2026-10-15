<template>
  <widget-settings-item :title="t('settings.mainParameter.title')">
    <c-progress-overlay :pending="pending" transition />
    <c-radio-group
      :model-value="modelValue.criteria"
      name="mainParameter"
      required
      @update:model-value="updateField('criteria', $event)"
    >
      <c-layout v-for="{ id, label, tooltip } of mainRatingSettingsWithCustom" :key="id">
        <v-radio :value="id" class="my-0" color="primary">
          <template #label>
            <span class="text-error mr-1">*</span>
            <span>{{ getCriteriaLabel(label) }}</span>
          </template>
        </v-radio>
        <c-help-icon
          v-if="tooltip"
          :text="tooltip"
          icon="help"
          top
        />
      </c-layout>
    </c-radio-group>
    <template v-if="isCustomCriteria">
      <c-name-field
        :model-value="modelValue.columnName"
        :label="t('settings.columnName')"
        name="columnName"
        required
        @update:model-value="updateField('columnName', $event)"
      />
      <field-filters-list
        :filters="modelValue.patterns"
        :widget-id="widgetId"
        name="patterns"
        with-entity
        addable
        editable
        required
        @update:filters="updateField('patterns', $event)"
      />
    </template>
  </widget-settings-item>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import {
  KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA,
  KPI_RATING_CRITERIA,
  KPI_RATING_SETTINGS_TYPES,
} from '@/constants';

import { isCustomCriteria as checkCustomCriteria } from '@/helpers/entities/metric/form';

import { useModelField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useRatingSettings } from '@/hooks/store/modules/rating-settings';

import FieldFiltersList from '@/components/sidebars/form/fields/filters-list.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: {
    WidgetSettingsItem,
    FieldFiltersList,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
    widgetId: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const pending = ref(false);
    const mainRatingSettings = ref([]);

    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);
    const { fetchRatingSettingsListWithoutStore } = useRatingSettings();

    const isEntityType = computed(() => props.type === KPI_RATING_SETTINGS_TYPES.entity);
    const isCustomCriteria = computed(() => checkCustomCriteria(props.modelValue.criteria));
    const criteriaLabels = computed(() => ({
      [KPI_RATING_CRITERIA.user]: t('common.user'),
      [KPI_RATING_CRITERIA.role]: t('common.role', 1),
      [KPI_RATING_CRITERIA.category]: t('common.category'),
      [KPI_RATING_CRITERIA.impactLevel]: t('common.impactLevel'),
    }));
    const mainRatingSettingsWithCustom = computed(() => isEntityType.value && !pending.value
        ? [
            ...mainRatingSettings.value,
            {
              id: KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA,
              label: t('settings.mainParameter.custom.label'),
              tooltip: t('settings.mainParameter.custom.tooltip'),
            },
          ]
        : mainRatingSettings.value,
    );

    const getCriteriaLabel = label => criteriaLabels.value[label] || label;

    const fetchList = async () => {
      pending.value = true;

      const { data = [] } = await fetchRatingSettingsListWithoutStore({
        params: {
          type: props.type,
          main: true,
          paginate: false,
        },
      });

      mainRatingSettings.value = data;
      pending.value = false;
    };

    onMounted(fetchList);

    return {
      t,
      pending,
      isCustomCriteria,
      mainRatingSettingsWithCustom,
      getCriteriaLabel,
      updateField,
    };
  },
};
</script>
