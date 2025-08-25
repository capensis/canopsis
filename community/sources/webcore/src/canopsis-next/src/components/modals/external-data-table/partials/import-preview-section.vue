<template>
  <div>
    <span class="text-subtitle-2">{{ $t('common.preview') }}</span>
    <external-data-table-records-list
      v-field="form"
      :records="records"
      :pending="pending"
      :options="options"
      :total-items="totalItems"
      :separator="separator"
      has-structure
      @update:options="updateOptions"
    />
  </div>
</template>

<script>
import ExternalDataTableRecordsList
  from '@/components/other/external-data-table/partials/external-data-table-records-list.vue';

export default {
  name: 'ImportPreviewSection',
  components: {
    ExternalDataTableRecordsList,
  },
  model: {
    prop: 'form',
    event: 'input',
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
  setup(props, { emit }) {
    const updateOptions = (newOptions) => {
      emit('update:options', newOptions);
    };

    return {
      updateOptions,
    };
  },
};
</script>
