<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.title",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('title')",
      name="title"
    )
    span.body-2.my-2 {{ $tc('common.column', 2) }}
    v-flex(xs12)
      v-alert(:value="!form.columns.length", color="info") You should add at least one column.
    c-columns-field(
      v-field="form.columns",
      :type="entityType",
      with-color-indicator,
      with-template,
      with-html,
      @input="validateRequiredRule"
    )
</template>

<script>
import { ENTITIES_TYPES, WIDGET_TEMPLATES_TYPES } from '@/constants';

import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

export default {
  inject: ['$validator'],
  mixins: [validationAttachRequiredMixin],
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
    entityType() {
      return this.form.type === WIDGET_TEMPLATES_TYPES.alarmColumns
        ? ENTITIES_TYPES.alarm
        : ENTITIES_TYPES.entity;
    },

    name() {
      return 'columns';
    },
  },
  mounted() {
    this.attachRequiredRule(() => !this.form.columns.length);
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
};
</script>
