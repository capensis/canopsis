<template>
  <v-chip
    :closable="closable && !disabled"
    :color="color"
    class="c-advanced-search__chip c-advanced-search__array-chip"
    size="small"
    @click.prevent="$emit('click')"
    @click:close="close"
    @focusout="$emit('focusout')"
  >
    <v-chip size="x-small">
      <span v-if="disabled">
        {{ chipValue.value }}
      </span>
      <input
        v-else
        :value="chipValue.value"
        :placeholder="$t('common.duration')"
        type="number"
        @input="updateValue"
      />
    </v-chip>
    <v-menu :disabled="disabled" location="bottom">
      <template #activator="{ props: activatorProps }">
        <v-chip v-bind="activatorProps" size="x-small">
          <span>{{ unitText }}</span>
        </v-chip>
      </template>
      <variables-list
        :value="chipValue.unit"
        :items="units"
        @update:value="updateUnit"
      />
    </v-menu>
  </v-chip>
</template>

<script>
import { computed } from 'vue';

import VariablesList from '@/components/common/text-editor/variables-list.vue';
import {
  AVAILABLE_TIME_UNITS,
  MEDIUM_AVAILABLE_TIME_UNITS,
  SHORT_AVAILABLE_TIME_UNITS,
  TIME_UNITS,
} from '@/constants';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  components: {
    VariablesList,
  },

  inheritAttrs: false,

  props: {
    modelValue: {
      type: Object,
      default: () => ({ value: 1, unit: TIME_UNITS.second }),
    },
    closable: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    long: {
      type: Boolean,
      default: false,
    },
    middle: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'click', 'focusout', 'close'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const chipValue = computed(() => props.modelValue ?? { value: 1, unit: TIME_UNITS.second });

    const units = computed(() => {
      const preparedUnits = {
        true: SHORT_AVAILABLE_TIME_UNITS,
        [props.long]: AVAILABLE_TIME_UNITS,
        [props.middle]: MEDIUM_AVAILABLE_TIME_UNITS,
      }.true;

      return Object.values(preparedUnits).map(({ value, text }) => ({
        value,
        text: t(text, chipValue.value.value || 0),
      }));
    });

    const unitText = computed(() => (
      units.value.find(({ value }) => value === chipValue.value.unit)?.text || chipValue.value.unit
    ));

    const updateValue = event => updateField('value', Number(event.target.value) || 1);
    const updateUnit = unit => updateField('unit', unit);
    const close = () => emit('close');

    return {
      chipValue,
      units,
      unitText,

      updateValue,
      updateUnit,
      close,
    };
  },
};
</script>
