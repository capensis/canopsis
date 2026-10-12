<template>
  <c-advanced-data-table
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    search
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <div>
        <template-testing-test-type-field
          :value="options.type"
          clearable
          @input="updateType"
        />
      </div>
    </template>
    <template #type="{ item }">
      <span>{{ $tc(`templateTesting.testTypes.${item.type}`) }}</span>
    </template>
    <template #actions="{ item }">
      <c-action-btn
        v-if="updatable"
        type="edit"
        @click="edit(item)"
      />
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="remove(item)"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import TemplateTestingTestTypeField from '../form/fields/template-testing-test-type-field.vue';

export default {
  components: { TemplateTestingTestTypeField },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.testName'),
        value: 'name',
      },
      {
        text: t('common.ruleType'),
        value: 'type',
      },
      {
        text: t('common.ruleName'),
        value: 'rule.name',
      },
      {
        text: t('common.description'),
        value: 'description',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    /**
     * Emits an edit event for the specified item
     *
     * @param {Object} item - The item object to be edited
     */
    const edit = item => emit('edit', item);

    /**
     * Emits a remove event for the specified item
     *
     * @param {Object} item - The item object to be removed
     */
    const remove = item => emit('remove', item);

    /**
     * Updates the table options and emits the changes to the parent component
     *
     * @param {Object} options - The updated options object
     */
    const updateOptions = options => emit('update:options', options);

    /**
     * Updates the type filter and emits the changes to the parent component
     *
     * @param {string} type - The new type filter value
     */
    const updateType = type => updateOptions({ ...props.options, type });

    return {
      headers,

      edit,
      remove,
      updateOptions,
      updateType,
    };
  },
};
</script>
