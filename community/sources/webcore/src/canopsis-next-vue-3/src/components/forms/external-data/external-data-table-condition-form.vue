<template>
  <v-row>
    <v-col class="pr-2" cols="3">
      <c-select
        :model-value="condition.type"
        :items="conditionTypes"
        :label="$t('common.type')"
        :disabled="disabled"
        item-title="text"
        item-value="value"
        @update:model-value="updateField('type', $event)"
      />
    </v-col>
    <v-col class="px-2" cols="4">
      <c-select
        :model-value="condition.attribute"
        :items="columnsForSelectedConditionType"
        :label="$t('common.attribute')"
        :name="conditionFieldName"
        :disabled="disabled"
        item-title="name"
        item-value="name"
        required
        @update:model-value="updateField('attribute', $event)"
      />
    </v-col>
    <v-col class="pl-2" cols="5">
      <c-layout align="center">
        <c-payload-text-field
          :value="condition.value"
          :label="$t('common.value')"
          :disabled="disabled"
          :variables="variables"
          :name="valueFieldName"
          clearable
          @update:value="updateField('value', $event)"
        />
        <v-btn v-if="!disabled" :disabled="disabledRemove" icon small @click="removeCondition">
          <v-icon color="error" small> delete </v-icon>
        </v-btn>
      </c-layout>
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_CONDITION_TYPES, EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
  emits: ['update:condition'],
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
    const { updateField } = useModelField(props, emit);

    const conditionTypes = computed(() =>
      Object.values(EXTERNAL_DATA_CONDITION_TYPES).map((type) => ({
        text: t(`externalData.conditionTypes.${type}`),
        value: type,
      })),
    );

    const columnsForSelectedConditionType = computed(() =>
      props.condition.type === EXTERNAL_DATA_CONDITION_TYPES.regexp
        ? props.columns.filter(
            (column) => column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp,
          )
        : props.columns,
    );

    const conditionFieldName = computed(() => `${props.name}.condition`);
    const valueFieldName = computed(() => `${props.name}.value`);
    const removeCondition = () => emit('remove', props.condition);

    return {
      conditionTypes,
      columnsForSelectedConditionType,
      conditionFieldName,
      valueFieldName,
      updateField,
      removeCondition,
    };
  },
};
</script>
