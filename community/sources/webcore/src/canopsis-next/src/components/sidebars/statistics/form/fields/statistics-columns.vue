<template>
  <widget-settings-item :title="$t('common.column', 2)">
    <c-progress-overlay :pending="pending" transition />
    <c-layout gap="3" column>
      <v-card v-for="(item, index) in modelValue" :key="item.key" class="lighten-2">
        <c-layout class="pt-2 px-2" align="center" justify="space-between">
          <c-layout>
            <v-btn
              :disabled="index === 0"
              icon="arrow_upward"
              variant="text"
              size="small"
              @click.prevent="moveColumn(index, -1)"
            />
            <v-btn
              :disabled="index === modelValue.length - 1"
              icon="arrow_downward"
              variant="text"
              size="small"
              @click.prevent="moveColumn(index, 1)"
            />
          </c-layout>
          <v-btn
            icon="close"
            color="error"
            variant="text"
            size="small"
            @click.prevent="removeItemFromArray(index)"
          />
        </c-layout>
        <v-card-text class="px-3 pb-3">
          <statistics-column
            :model-value="item"
            :type="type"
            :name="item.key"
            :rating-settings="ratingSettings"
            @update:model-value="updateItemInArray(index, $event)"
          />
        </v-card-text>
      </v-card>
      <c-layout justify="end">
        <c-btn-with-error
          :error="hasError ? errorMessage : ''"
          variant="outlined"
          @click.prevent="add"
        >
          {{ $t('common.add') }}
        </c-btn-with-error>
      </c-layout>
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { computed, onMounted, ref, toRef, watch } from 'vue';

import { KPI_RATING_SETTINGS_TYPES } from '@/constants';

import { useRatingSettingsStore } from '@/stores/modules/rating-settings';

import { statisticsWidgetColumnToForm } from '@/helpers/entities/widget/forms/statistics';

import { useArrayModelField, useField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

import StatisticsColumn from './statistics-column.vue';

export default {
  components: {
    WidgetSettingsItem,
    StatisticsColumn,
  },

  props: {
    modelValue: {
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

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const ratingSettingsStore = useRatingSettingsStore();
    const ratingSettings = ref([]);
    const pending = ref(false);
    const fieldName = toRef(props, 'name');
    const { addItemIntoArray, updateModel, updateItemInArray, removeItemFromArray } =
      useArrayModelField(props, emit);

    const { errorMessage, validate } = useField(
      fieldName,
      () => (props.modelValue.length > 0 ? true : t('widgetTemplate.errors.columnsRequired')),
      {
        initialValue: props.modelValue.length,
        validateOnValueUpdate: false,
      },
    );

    watch(
      () => props.modelValue.length,
      () => validate(),
      { flush: 'post' },
    );

    const hasError = computed(() => !!errorMessage.value);

    const add = () => addItemIntoArray(statisticsWidgetColumnToForm());

    const moveColumn = (index, offset) => {
      const targetIndex = index + offset;

      if (targetIndex < 0 || targetIndex >= props.modelValue.length) {
        return;
      }

      const columns = [...props.modelValue];
      [columns[index], columns[targetIndex]] = [columns[targetIndex], columns[index]];

      updateModel(columns);
    };

    const fetchList = async () => {
      pending.value = true;

      const { data = [] } = await ratingSettingsStore.fetchListWithoutStore({
        params: {
          /**
           * We don't have rating settings for KPI_RATING_SETTINGS_TYPES.user type with `main` === false
           */
          type: KPI_RATING_SETTINGS_TYPES.entity,
          main: false,
          paginate: false,
        },
      });

      ratingSettings.value = data;
      pending.value = false;
    };

    onMounted(fetchList);

    return {
      errorMessage,
      hasError,
      pending,
      ratingSettings,

      add,
      moveColumn,
      updateItemInArray,
      removeItemFromArray,
    };
  },
};
</script>
