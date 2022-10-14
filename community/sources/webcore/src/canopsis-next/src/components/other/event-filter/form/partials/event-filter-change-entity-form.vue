<template lang="pug">
  div
    v-text-field(
      v-field="form.resource",
      :label="$t('eventFilter.resource')"
    )
    v-text-field(
      v-field="form.component",
      :label="$t('eventFilter.component')"
    )
    v-text-field(
      v-field="form.connector",
      :label="$t('eventFilter.connector')"
    )
    v-text-field(
      v-field="form.connector_name",
      :label="$t('eventFilter.connectorName')"
    )
    v-alert(:value="errors.has(name)", type="error") {{ $t('eventFilter.configRequired') }}
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
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
  },
  watch: {
    form() {
      this.$validator.validate(this.name);
    },
  },
  created() {
    this.attachConfigRequiredRule();
  },
  beforeDestroy() {
    this.detachConfigRequiredRule();
  },
  methods: {
    attachConfigRequiredRule() {
      this.$validator.attach({
        name: this.name,
        rules: 'required:true',
        getter: () => this.form.resource
          || this.form.component
          || this.form.connector
          || this.form.connector_name,
      });
    },

    detachConfigRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
