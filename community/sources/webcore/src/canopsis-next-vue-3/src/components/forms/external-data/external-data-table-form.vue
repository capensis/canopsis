<template>
  <c-layout column>
    <c-layout gap="4">
      <c-enabled-field
        v-if="optionally"
        :model-value="form.optional"
        :label="$t('common.optional')"
        :disabled="disabled"
        @update:model-value="updateField('optional', $event)"
      />
      <external-data-table-table-field
        :value="form.table"
        :label="$t('externalData.fields.collection')"
        :disabled="disabled"
        required
        @update:value="updateField('table', $event)"
        @update:selected-items="updateSelectedItems"
      />
    </c-layout>
    <v-row>
      <v-col cols="6">
        <c-select
          :model-value="form.sort_by"
          :items="preparedColumns"
          :label="$t('externalData.fields.sortBy')"
          :name="sortByFieldName"
          :disabled="disabled"
          item-title="name"
          item-value="name"
          @update:model-value="updateField('sort_by', $event)"
        />
      </v-col>
      <v-col class="ml-3" cols="6">
        <c-select
          :model-value="form.sort"
          :items="sortOrders"
          :label="$t('common.sort')"
          :name="sortFieldName"
          :disabled="disabled"
          clearable
          @update:model-value="updateField('sort', $event)"
        />
      </v-col>
    </v-row>
    <external-data-table-condition-form
      v-for="(condition, index) in form.conditions"
      :key="condition.key"
      :condition="form.conditions[index]"
      :name="`${name}.conditions.${condition.key}`"
      :disabled-remove="hasOnlyOneCondition"
      :disabled="disabled"
      :variables="variables"
      :columns="columns"
      @update:condition="updateCondition(index, $event)"
      @remove="removeCondition(index)"
    />
    <v-btn
v-if="!disabled"
class="ml-0 mb-0"
color="primary"
outlined
@click="addCondition">
      {{ $t('common.addMore') }}
    </v-btn>
  </c-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { SORT_ORDERS } from '@/constants';

import { addPriorityColumnToColumnsArray } from '@/helpers/entities/external-data-table/form';
import { externalDataItemConditionAttributeToForm } from '@/helpers/entities/shared/external-data/form';

import { useModelField } from '@/hooks/form/model-field';

import ExternalDataTableTableField from '@/components/other/external-data-table/form/fields/external-data-table-table-field.vue';

import ExternalDataTableConditionForm from './external-data-table-condition-form.vue';

export default {
  emits: ['update:form'],
  components: { ExternalDataTableTableField, ExternalDataTableConditionForm },
  props: {
    form: {
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
    variables: {
      type: Array,
      default: () => [],
    },
    optionally: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const sortOrders = Object.values(SORT_ORDERS).map(order => order.toLowerCase());
    const columns = ref([]);

    const hasOnlyOneCondition = computed(() => props.form.conditions.length === 1);
    const collectionFieldName = computed(() => `${props.name}.collection`);
    const sortFieldName = computed(() => `${props.name}.sort`);
    const sortByFieldName = computed(() => `${props.name}.sort_by`);

    const preparedColumns = computed(() => addPriorityColumnToColumnsArray(columns.value));

    /**
     * Updates one condition in the conditions array.
     *
     * @param {number} index
     * @param {Object} condition
     */
    const updateCondition = (index, condition) => {
      const conditions = [...props.form.conditions];

      conditions[index] = condition;

      return updateField('conditions', conditions);
    };

    /**
     * Adds a new condition to the external data table form.
     */
    const addCondition = () =>
      updateField('conditions', [
        ...props.form.conditions,

        externalDataItemConditionAttributeToForm(),
      ]);

    /**
     * Removes a condition from the external data table form by index.
     *
     * @param {number} index - The index of the condition to remove from the conditions array
     */
    const removeCondition = index =>
      updateField(
        'conditions',
        props.form.conditions.filter((condition, conditionIndex) => index !== conditionIndex),
      );

    /**
     * Updates the available columns based on selected external data table.
     *
     * @param {Array} [tables=[]] - Array of selected tables, expects first item to be the main table
     * @param {Object} [tables[0]={}] - The selected external data table object
     * @param {Array} [tables[0].column_configs] - Array of column configuration objects
     * @param {string} tables[0].column_configs[].type - Column data type
     * @param {string} tables[0].column_configs[].name - Column name identifier
     */
    const updateSelectedItems = ([table = {}] = []) =>
      (columns.value = table?.column_configs ?? []);

    return {
      columns,
      preparedColumns,

      sortOrders,
      hasOnlyOneCondition,
      collectionFieldName,
      sortFieldName,
      sortByFieldName,

      updateField,
      updateCondition,
      addCondition,
      removeCondition,
      updateSelectedItems,
    };
  },
};
</script>
