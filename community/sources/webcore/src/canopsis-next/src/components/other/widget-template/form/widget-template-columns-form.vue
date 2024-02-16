<template>
  <v-layout column>
    <v-text-field
      v-field="form.title"
      v-validate="'required'"
      :label="$t('common.name')"
      :error-messages="errors.collect('title')"
      name="title"
    />
    <span class="text-body-2 my-2">{{ $tc('common.column', 2) }}</span>
    <v-flex xs12>
      <v-alert
        :value="!form.columns.length"
        color="info"
      >
        {{ $t('widgetTemplate.errors.columnsRequired') }}
      </v-alert>
    </v-flex>
    <c-columns-field
      v-field="form.columns"
      :type="entityType"
      with-color-indicator
      with-template
      with-html
      @input="validateRequiredRule"
    />
  </v-layout>
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
    this.attachRequiredRule(() => this.form.columns.length > 0);
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
};
</script>
