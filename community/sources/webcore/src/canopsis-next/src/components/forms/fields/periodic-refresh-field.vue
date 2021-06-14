<template lang="pug">
  v-layout.mb-3(align-top)
    v-flex(xs5)
      c-enabled-field(
        v-field="periodicRefresh.enabled",
        :label="label",
        hide-details,
        @input="validateDuration"
      )
    v-flex(xs7)
      c-duration-field(
        v-field="periodicRefresh",
        :disabled="!periodicRefresh.enabled",
        :required="periodicRefresh.enabled",
        :name="durationFieldName",
        :min="0"
      )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'periodicRefresh',
    event: 'input',
  },
  props: {
    periodicRefresh: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'periodic_refresh',
    },
  },
  computed: {
    durationFieldName() {
      return `${this.name}.duration`;
    },
  },
  methods: {
    validateDuration() {
      this.$nextTick(() => this.$validator.validate(`${this.durationFieldName}.value`));
    },
  },
};
</script>
