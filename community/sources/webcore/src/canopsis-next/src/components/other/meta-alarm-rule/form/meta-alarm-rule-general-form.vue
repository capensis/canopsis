<template>
  <v-layout column>
    <c-id-field
      v-field="form._id"
      :disabled="disabledIdField"
      :help-text="$t('metaAlarmRule.idHelp')"
      autofocus
    />
    <c-name-field
      v-field="form.name"
      :autofocus="disabledIdField"
      required
    />
    <c-payload-textarea-field
      v-field="form.output_template"
      :label="$t('metaAlarmRule.outputTemplate')"
      :help-text="$t('metaAlarmRule.outputTemplateHelp')"
      :variables="variables"
    />
    <meta-alarm-rule-tags-form
      v-field="form.tags"
      :variables="variables"
      class="mb-4"
    />
    <meta-alarm-rule-infos-form v-field="form.infos" />
    <v-layout>
      <c-enabled-field
        v-field="form.auto_resolve"
        :label="$t('metaAlarmRule.autoResolve')"
      />
    </v-layout>
  </v-layout>
</template>

<script>
import MetaAlarmRuleTagsForm from './meta-alarm-rule-tags-form.vue';
import MetaAlarmRuleInfosForm from './meta-alarm-rule-infos-form.vue';

export default {
  inject: ['$validator'],
  components: { MetaAlarmRuleTagsForm, MetaAlarmRuleInfosForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    disabledIdField: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
