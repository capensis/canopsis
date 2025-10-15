<template>
  <v-layout>
    <v-flex
      class="pr-2"
      xs3
    >
      <v-select
        v-field="condition.type"
        :items="conditionTypes"
        :label="$t('common.type')"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex
      class="px-2"
      xs4
    >
      <v-select
        v-field="condition.attribute"
        v-validate="'required'"
        :items="columnsForSelectedConditionType"
        :label="$t('common.attribute')"
        :name="conditionFieldName"
        :error-messages="errors.collect(conditionFieldName)"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex
      class="pl-2"
      xs5
    >
      <v-layout align-center>
        <c-payload-text-field
          v-field="condition.value"
          :label="$t('common.value')"
          :disabled="disabled"
          :variables="variables"
          :name="valueFieldName"
          clearable
        />
        <v-btn
          v-if="!disabled"
          :disabled="disabledRemove"
          icon
          small
          @click="removeCondition"
        >
          <v-icon
            color="error"
            small
          >
            delete
          </v-icon>
        </v-btn>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_CONDITION_TYPES, EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  inject: ['$validator'],
  model: {
    prop: 'condition',
    event: 'input',
  },
  props: {
    condition: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledRemove: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => [],
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const conditionTypes = computed(() => Object.values(EXTERNAL_DATA_CONDITION_TYPES)
      .map(type => ({ text: t(`externalData.conditionTypes.${type}`), value: type })));

    const columnsForSelectedConditionType = computed(() => (
      props.condition.type === EXTERNAL_DATA_CONDITION_TYPES.regexp
        ? props.columns.filter(column => column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp)
        : props.columns
    ));

    const conditionFieldName = computed(() => `${props.name}.condition`);

    const valueFieldName = computed(() => `${props.name}.value`);

    const removeCondition = () => emit('remove', props.condition);

    return {
      conditionTypes,
      columnsForSelectedConditionType,
      conditionFieldName,
      valueFieldName,
      removeCondition,
    };
  },
};
</script>
