<template>
  <v-chip
    :close="closable && !disabled"
    :color="color"
    class="c-advanced-search__chip c-advanced-search__array-chip"
    @click.prevent=""
    @click:close="close"
  >
    <v-chip>
      <span v-if="disabled">{{ value.value }}</span>
      <input
        v-else
        :value="value.value"
        :placeholder="$t('common.duration')"
        type="number"
        @input="updateValue"
      >
    </v-chip>
    <v-menu :disabled="disabled" bottom>
      <template #activator="{ on }">
        <v-chip v-on="on">
          <span>{{ unitText }}</span>
        </v-chip>
      </template>
      <variables-list
        v-field="value.unit"
        :items="units"
      />
    </v-menu>
  </v-chip>
</template>

<script>
import { computed } from 'vue';

import { AVAILABLE_TIME_UNITS, SHORT_AVAILABLE_TIME_UNITS, TIME_UNITS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

export default {
  components: {
    VariablesList,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
      required: false,
    },
    long: {
      type: Boolean,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { tc } = useI18n();
    const { updateField } = useModelField(props, emit);

    const units = computed(() => {
      const preparedUnits = props.long
        ? AVAILABLE_TIME_UNITS
        : SHORT_AVAILABLE_TIME_UNITS;

      return Object.values(preparedUnits).map(({ value, text }) => ({
        value,
        text: tc(text, props.value.value || 0),
      }));
    });

    const unitText = computed(() => (
      units.value.find(({ value }) => value === props.value.unit)?.text || props.value.unit
    ));

    const updateValue = event => updateField('value', Number(event.target.value) || 1);

    const close = () => emit('close');

    return {
      units,
      unitText,

      updateValue,
      close,
    };
  },
};
</script>
