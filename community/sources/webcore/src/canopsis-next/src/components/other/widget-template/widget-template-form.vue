<template lang="pug">
  v-layout(column)
    v-select(
      v-field="form.type",
      v-validate="'required'",
      :label="$t('common.type')",
      :items="availableTypes",
      :error-messages="errors.collect('type')",
      name="type"
    )
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.templateName')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    span.body-2.mb-2 {{ $tc('common.column', 2) }}
    c-columns-field(
      v-field="form.columns",
      :type="form.type"
    )
</template>

<script>
import { WIDGET_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    availableTypes() {
      return [
        {
          value: WIDGET_TYPES.alarmList,
          text: this.$t(`modals.createWidget.types.${WIDGET_TYPES.alarmList}.title`),
        },
        {
          value: WIDGET_TYPES.context,
          text: this.$t(`modals.createWidget.types.${WIDGET_TYPES.context}.title`),
        },
      ];
    },
  },
};
</script>
