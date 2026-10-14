<template>
  <v-chip
    :closable="closable"
    :color="color"
    class="c-advanced-search__chip c-advanced-search__array-chip"
    @click.prevent="$emit('click')"
    @click:close="close"
    @focusout="$emit('focusout')"
  >
    <v-menu
      v-for="chip in chips"
      :key="chip.key"
      :close-on-content-click="true"
      max-width="290"
      location="bottom"
    >
      <template #activator="{ props: activatorProps }">
        <v-chip
          class="c-advanced-search__chip"
          v-bind="activatorProps"
        >
          <span class="font-italic text-medium-emphasis">
            {{ chip.label }}:
          </span>
          {{ chip.valueText }}
        </v-chip>
      </template>
      <variables-list
        :value="chipValue[chip.key]"
        :items="chip.items"
        @update:value="updateField(chip.key, $event)"
      />
    </v-menu>
  </v-chip>
</template>

<script>
import { computed, toRef } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { useQuickDateIntervalRange } from '@/hooks/form/quick-date-interval-range';
import { useI18n } from '@/hooks/i18n';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

export default {
  components: {
    VariablesList,
  },

  inheritAttrs: false,

  props: {
    modelValue: {
      type: Object,
      default: () => ({ from: '', to: '' }),
    },
    closable: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
    intervalRanges: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue', 'click', 'focusout', 'close'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const chipValue = computed(() => props.modelValue ?? { from: '', to: '' });

    const {
      fromQuickRanges,
      toQuickRanges,
      itemFromDisabled,
      itemToDisabled,
    } = useQuickDateIntervalRange({
      fromRanges: toRef(props, 'intervalRanges'),
      toRanges: toRef(props, 'intervalRanges'),
      value: chipValue,
    });

    const chipConfigs = computed(() => ({
      from: { ranges: fromQuickRanges, itemDisabled: itemFromDisabled },
      to: { ranges: toQuickRanges, itemDisabled: itemToDisabled },
    }));

    const chips = computed(() => Object.entries(chipConfigs.value).map(([key, config]) => {
      const chipRanges = config.ranges.value;
      const value = chipValue.value?.[key];

      const valueText = !value
        ? '---------'
        : (chipRanges.find(item => item.value === value)?.text ?? t(`quickRanges.types.${value}`));

      const items = chipRanges.map(item => ({
        ...item,
        disabled: config.itemDisabled(item),
      }));

      return {
        key,
        label: t(`common.${key}`),
        value,
        valueText,
        items,
      };
    }));

    const close = () => emit('close');

    return {
      chipValue,
      chips,

      updateField,
      close,
    };
  },
};
</script>
