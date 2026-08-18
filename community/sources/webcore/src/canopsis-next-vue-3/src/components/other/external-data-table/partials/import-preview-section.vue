<template>
  <div>
    <span class="text-subtitle-2">
      {{ t('common.preview') }}
    </span>
    <external-data-table-records-list
      :columns="form"
      :original-columns="originalForm"
      :records="records"
      :pending="pending"
      :options="options"
      :total-items="totalItems"
      :separator="separator"
      has-structure
      @update:columns="updateForm"
      @update:options="updateOptions"
    />
  </div>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

import ExternalDataTableRecordsList from '@/components/other/external-data-table/partials/external-data-table-records-list.vue';

export default {
  components: {
    ExternalDataTableRecordsList,
  },

  props: {
    activeImportFileId: {
      type: String,
      default: null,
    },
    form: {
      type: Object,
      required: true,
    },
    originalForm: {
      type: Object,
      required: false,
    },
    records: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    separator: {
      type: String,
      required: true,
    },
  },

  emits: ['input', 'update:form', 'update:options'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const updateForm = (newForm) => {
      emit('update:form', newForm);
      emit('input', newForm);
    };

    const updateOptions = newOptions => emit('update:options', newOptions);

    return {
      t,
      updateForm,
      updateOptions,
    };
  },
};
</script>
