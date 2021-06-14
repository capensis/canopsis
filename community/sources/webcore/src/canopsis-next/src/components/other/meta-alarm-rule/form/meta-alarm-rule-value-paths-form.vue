<template lang="pug">
  div
    text-fields(
      v-field="form.value_paths",
      :title="$tc('metaAlarmRule.fields.valuePath', 2)",
      :label="$tc('metaAlarmRule.fields.valuePath', 1)",
      validation-rules="required",
      @input="validateValuePaths"
    )
    v-alert(:value="errors.has('value_paths')", type="error")
      span {{ $t('metaAlarmRule.errors.noValuePaths') }}
</template>

<script>
import TextFields from '@/components/forms/fields/text-fields.vue';

export default {
  inject: ['$validator'],
  components: { TextFields },
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
  created() {
    this.$validator.attach({
      name: 'value_paths',
      rules: 'required:true',
      getter: () => this.form.value_paths && this.form.value_paths.length > 0,
      context: () => this,
    });
  },
  beforeDestroy() {
    this.$validator.attach('value_paths');
  },
  methods: {
    validateValuePaths() {
      this.$nextTick(() => this.$validator.validate('value_paths'));
    },
  },
};
</script>
