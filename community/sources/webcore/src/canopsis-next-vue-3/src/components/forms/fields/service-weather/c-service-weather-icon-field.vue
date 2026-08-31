<template>
  <c-select
    :model-value="modelValue"
    :items="availableIcons"
    :label="label"
    :loading="pending"
    :name="name"
    :disabled="disabled"
    item-title="text"
    item-value="icon"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #selection="{ item }">
      <v-icon :class="disabled ? 'text-disabled' : ''">
        {{ item.raw.icon }}
      </v-icon>
      <span class="ml-2">
        {{ item.raw.text }}
      </span>
    </template>

    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps">
        <template #prepend>
          <v-icon>{{ item.raw.icon }}</v-icon>
        </template>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import {
  MAX_LIMIT,
  PBEHAVIOR_TYPE_TYPES,
  SERVICE_STATES,
  WEATHER_ICONS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'icon',
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();
    const { fetchPbehaviorTypesListWithoutStore } = usePbehaviorType();

    const pbehaviorTypes = ref([]);
    const pending = ref(false);

    const serviceWeatherIcons = computed(() => [
      {
        icon: WEATHER_ICONS[SERVICE_STATES.ok],
        text: t('serviceWeather.iconTypes.ok'),
      },
      {
        icon: WEATHER_ICONS[SERVICE_STATES.minor],
        text: t('serviceWeather.iconTypes.minorOrMajor'),
      },
      {
        icon: WEATHER_ICONS[SERVICE_STATES.critical],
        text: t('serviceWeather.iconTypes.critical'),
      },
    ]);

    const pbehaviorTypeIcons = computed(() => pbehaviorTypes.value.map(pbehaviorType => ({
        icon: pbehaviorType.icon_name,
        text: pbehaviorType.name,
      })),
    );

    const availableIcons = computed(() => [
      ...serviceWeatherIcons.value,
      ...pbehaviorTypeIcons.value,
    ]);

    const fetchTypesList = async () => {
      try {
        pending.value = true;

        const { data: types = [] } = await fetchPbehaviorTypesListWithoutStore({
          params: {
            types: [
              PBEHAVIOR_TYPE_TYPES.inactive,
              PBEHAVIOR_TYPE_TYPES.maintenance,
              PBEHAVIOR_TYPE_TYPES.pause,
            ],
            limit: MAX_LIMIT,
          },
        });

        pbehaviorTypes.value = types;
      } finally {
        pending.value = false;
      }
    };

    onMounted(fetchTypesList);

    return {
      pending,
      availableIcons,
    };
  },
};
</script>
