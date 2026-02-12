<template>
  <v-layout column>
    <v-layout class="gap-4">
      <c-enabled-field
        v-if="optionally"
        v-field="form.optional"
        :label="$t('common.optional')"
        :disabled="disabled"
      />
      <external-data-table-table-field
        v-field="form.table"
        :label="$t('externalData.fields.collection')"
        :disabled="disabled"
        required
        @update:selected-items="updateSelectedItems"
      />
    </v-layout>
    <v-layout>
      <v-flex xs6>
        <v-select
          v-field="form.sort_by"
          :items="preparedColumns"
          :label="$t('externalData.fields.sortBy')"
          :name="sortByFieldName"
          :error-messages="errors.collect(sortByFieldName)"
          :disabled="disabled"
          item-text="name"
          item-value="name"
        />
      </v-flex>
      <v-flex
        class="ml-3"
        xs6
      >
        <v-select
          v-field="form.sort"
          :items="sortOrders"
          :label="$t('externalData.fields.sort')"
          :name="sortFieldName"
          :error-messages="errors.collect(sortFieldName)"
          :disabled="disabled"
          clearable
        />
      </v-flex>
    </v-layout>
    <external-data-table-condition-form
      v-for="(condition, index) in form.conditions"
      v-field="form.conditions[index]"
      :key="condition.key"
      :name="`${name}.conditions.${condition.key}`"
      :disabled-remove="hasOnlyOneCondition"
      :disabled="disabled"
      :variables="variables"
      :columns="columns"
      @remove="removeCondition(index)"
    />
    <v-flex v-if="!disabled">
      <v-btn
        class="ml-0 mb-0"
        color="primary"
        outlined
        @click="addCondition"
      >
        {{ $t('common.addMore') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { SORT_ORDERS } from '@/constants';

import { addPriorityColumnToColumnsArray } from '@/helpers/entities/external-data-table/form';
import { externalDataItemConditionAttributeToForm } from '@/helpers/entities/shared/external-data/form';

import { useModelField } from '@/hooks/form/model-field';

import ExternalDataTableTableField
  from '@/components/other/external-data-table/form/fields/external-data-table-table-field.vue';

import ExternalDataTableConditionForm from './external-data-table-condition-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataTableTableField, ExternalDataTableConditionForm },
  model: {
    prop: 'form',
    event: 'input',
  },
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
      default: () => ([]),
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
     * Adds a new condition to the external data table form.
     */
    const addCondition = () => updateField('conditions', [
      ...props.form.conditions,

      externalDataItemConditionAttributeToForm(),
    ]);

    /**
     * Removes a condition from the external data table form by index.
     *
     * @param {number} index - The index of the condition to remove from the conditions array
     */
    const removeCondition = index => updateField(
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
    const updateSelectedItems = ([table = {}] = []) => columns.value = table?.column_configs ?? [];

    return {
      columns,
      preparedColumns,

      sortOrders,
      hasOnlyOneCondition,
      collectionFieldName,
      sortFieldName,
      sortByFieldName,

      addCondition,
      removeCondition,
      updateSelectedItems,
    };
  },
};
</script>
