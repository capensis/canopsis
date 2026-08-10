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
        <template-testing-data-type-field
          :value="options.type"
          hide-details
          clearable
          @update:value="updateType"
        />
      </div>
    </template>

    <template #type="{ item }">
      <span>{{ t(`templateTesting.dataTypes.${item.type}`) }}</span>
    </template>

    <template #actions="{ item }">
      <c-action-btn v-if="updatable" type="edit" @click="edit(item)" />
      <c-action-btn v-if="removable" type="delete" @click="remove(item)" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import TemplateTestingDataTypeField from '../form/fields/template-testing-data-type-field.vue';

export default {
  components: { TemplateTestingDataTypeField },
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
  emits: ['edit', 'remove', 'update:options'],
  setup(props, { emit }) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
      },
      {
        text: t('common.type'),
        value: 'type',
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

    const edit = item => emit('edit', item);
    const remove = item => emit('remove', item);
    const updateOptions = options => emit('update:options', options);
    const updateType = type => updateOptions({ ...props.options, type });

    return {
      t,
      headers,
      edit,
      remove,
      updateOptions,
      updateType,
    };
  },
};
</script>
