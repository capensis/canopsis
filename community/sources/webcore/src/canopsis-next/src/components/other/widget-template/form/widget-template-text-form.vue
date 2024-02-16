<template>
  <v-layout column>
    <v-text-field
      v-field="form.title"
      v-validate="'required'"
      :label="$t('common.name')"
      :error-messages="errors.collect('title')"
      class="mb-2"
      name="title"
    />
    <text-editor-field
      v-field="form.content"
      v-validate="'required'"
      :error-messages="errors.collect('content')"
      :variables="variables"
      :dark="$system.dark"
      name="content"
    />
  </v-layout>
</template>

<script>
import { WIDGET_TEMPLATES_TYPES } from '@/constants';

import { alarmVariablesMixin } from '@/mixins/widget/variables/alarm';
import { entityVariablesMixin } from '@/mixins/widget/variables/entity';

import TextEditorField from '@/components/common/text-editor/text-editor.vue';

export default {
  inject: ['$validator', '$system'],
  components: { TextEditorField },
  mixins: [alarmVariablesMixin, entityVariablesMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    entityInfos: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    variables() {
      return this.form.type === WIDGET_TEMPLATES_TYPES.alarmMoreInfos
        ? this.alarmVariables
        : this.entityVariables;
    },
  },
};
</script>
