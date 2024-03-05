<template>
  <div>
    <c-payload-text-field
      v-field="form.resource"
      :label="$t('eventFilter.resource')"
      :name="`${name}.resource`"
      :variables="variables"
    />
    <c-payload-text-field
      v-field="form.component"
      :label="$t('eventFilter.component')"
      :name="`${name}.component`"
      :variables="variables"
    />
    <c-payload-text-field
      v-field="form.connector"
      :label="$t('eventFilter.connector')"
      :name="`${name}.connector`"
      :variables="variables"
    />
    <c-payload-text-field
      v-field="form.connector_name"
      :label="$t('eventFilter.connectorName')"
      :name="`${name}.connector_name`"
      :variables="variables"
    />
    <v-alert
      :value="errors.has(name)"
      type="error"
    >
      {{ $t('eventFilter.configRequired') }}
    </v-alert>
  </div>
</template>

<script>
import { formMixin } from '@/mixins/form';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

export default {
  inject: ['$validator'],
  mixins: [formMixin, validationAttachRequiredMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'config',
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  watch: {
    form() {
      this.validateRequiredRule();
    },
  },
  created() {
    this.attachRequiredRule(
      () => this.form.resource || this.form.component || this.form.connector || this.form.connector_name,
    );
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
};
</script>
