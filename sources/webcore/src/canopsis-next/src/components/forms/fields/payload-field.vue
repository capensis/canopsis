<template lang="pug">
  v-textarea(
    v-model="localValue",
    v-validate="'required|json'",
    :label="$t('common.payload')",
    :error-messages="errorsMessages",
    name="payload",
    @blur="updatePayload($event)"
  )
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      localValue: this.convertPayloadValueToString(this.value),
    };
  },
  computed: {
    errorsMessages() {
      return this.fields.payload && this.fields.payload.touched ? this.errors.collect('payload') : [];
    },
  },
  watch: {
    value() {
      this.localValue = this.convertPayloadValueToString(this.value);
    },
  },
  methods: {
    async updatePayload() {
      const isValid = await this.$validator.validate('payload');

      if (isValid) {
        this.updateModel(JSON.parse(this.localValue));
      }
    },

    convertPayloadValueToString(value) {
      return JSON.stringify(value, null, 2);
    },
  },
};
</script>

