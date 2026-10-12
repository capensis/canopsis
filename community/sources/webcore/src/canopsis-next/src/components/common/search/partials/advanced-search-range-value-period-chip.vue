<template>
  <v-chip
    :close="closable"
    :color="color"
    class="c-advanced-search__chip c-advanced-search__array-chip"
    @click.prevent=""
    @click:close="close"
  >
    <v-menu
      v-for="chip in chips"
      :key="chip.key"
      :close-on-content-click="true"
      max-width="290"
      offset-y
    >
      <template #activator="{ on }">
        <v-chip
          class="c-advanced-search__chip"
          v-on="on"
        >
          <span class="font-italic grey--text">
            {{ chip.label }}:
          </span>
          {{ chip.valueText }}
        </v-chip>
      </template>
      <variables-list
        v-field="chips[chip.key]"
        :items="chip.items"
      />
    </v-menu>
  </v-chip>
</template>

<script>
import { computed, toRef } from 'vue';

import { useQuickDateIntervalRange } from '@/hooks/form/quick-date-interval-range';
import { useI18n } from '@/hooks/i18n';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

export default {
  components: {
    VariablesList,
  },
  props: {
    value: {
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
  setup(props, { emit }) {
    const { t } = useI18n();

    const {
      fromQuickRanges,
      toQuickRanges,
      itemFromDisabled,
      itemToDisabled,
    } = useQuickDateIntervalRange({
      fromRanges: toRef(props, 'intervalRanges'),
      toRanges: toRef(props, 'intervalRanges'),
      value: toRef(props, 'value'),
    });

    const chipConfigs = computed(() => ({
      from: { ranges: fromQuickRanges, itemDisabled: itemFromDisabled },
      to: { ranges: toQuickRanges, itemDisabled: itemToDisabled },
    }));

    /**
     * Computed chips for 'from' and 'to' with all attributes needed for rendering.
     *
     * @returns {Array<Object>} Chips with key, label, value, valueText, items, updateValue.
     */
    const chips = computed(() => (
      Object.entries(chipConfigs.value).map(([key, { ranges, itemDisabled }]) => {
        const chipRanges = ranges.value;
        const chipValue = props.value?.[key];

        const valueText = !chipValue
          ? '---------'
          : (chipRanges.find(v => v.value === chipValue)?.text ?? t(`quickRanges.types.${chipValue}`));

        const items = chipRanges.map(item => ({
          ...item,
          disabled: itemDisabled(item),
        }));

        return {
          key,
          label: t(`common.${key}`),
          value: chipValue,
          valueText,
          items,
        };
      })
    ));

    /**
     * Emits the close event when the chip is closed.
     */
    const close = () => emit('close');

    return {
      chips,
      close,
    };
  },
};
</script>
