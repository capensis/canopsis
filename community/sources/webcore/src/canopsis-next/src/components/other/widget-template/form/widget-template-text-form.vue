<template lang="pug">
  v-layout(column)
    v-text-field.mb-2(
      v-field="form.title",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('title')",
      name="title"
    )
    text-editor-field(
      v-model="form.content",
      v-validate="'required'",
      :error-messages="errors.collect('content')",
      :variables="variables",
      :dark="$system.dark",
      name="content"
    )
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
