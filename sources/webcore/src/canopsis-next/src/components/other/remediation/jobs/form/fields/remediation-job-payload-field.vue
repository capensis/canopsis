<template lang="pug">
  v-textarea(
    v-model="localValue",
    v-validate="'required|json'",
    :label="$t('modals.createRemediationJob.fields.payload')",
    :error-messages="errors.collect('payload')",
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
  watch: {
    value() {
      this.localValue = this.convertPayloadValueToString(this.value);
    },
  },
  created() {
    this.$validator.extend('json', {
      getMessage: () => this.$t('modals.createRemediationJob.errors.invalidJSON'),
      validate: (value) => {
        try {
          return !!JSON.parse(value);
        } catch (err) {
          return false;
        }
      },
    });
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

