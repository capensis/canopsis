<template lang="pug">
  v-layout(column)
    v-layout
      v-text-field(
        v-field="operation.name",
        v-validate="'required'",
        :disabled="operation.saved",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name",
        box
      )
    remediation-instruction-time-to-complete-field(
      v-field="operation.time_to_complete",
      :disabled="operation.saved"
    )
    v-layout
      v-textarea(
        v-field="operation.description",
        v-validate="'required'",
        :label="$t('common.description')",
        :error-messages="errors.collect('description')",
        :disabled="operation.saved",
        name="description",
        box
      )
    v-layout(v-if="!operation.saved", justify-end)
      v-btn.mt-0(depressed, flat, @click="cancelChangeOperation") {{ $t('common.cancel') }}
      v-btn.mt-0.mr-0.primary(@click="saveOperation") {{ $t('common.save') }}
</template>

<script>
import formMixin from '@/mixins/form';

import RemediationInstructionTimeToCompleteField from './remediation-instruction-time-to-complete-field.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { RemediationInstructionTimeToCompleteField },
  mixins: [formMixin],
  model: {
    prop: 'operation',
    event: 'input',
  },
  props: {
    operation: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      oldOperation: null,
    };
  },
  methods: {
    editName() {
      this.oldOperation = this.operation;

      this.updateField('saved', false);
    },

    cancelChangeOperation() {
      if (this.oldOperation) {
        this.updateModel({
          ...this.oldOperation,
          saved: true,
        });
      } else {
        this.$emit('remove');
      }
    },

    async saveOperation() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        this.oldOperation = null;

        this.updateField('saved', true);
      }
    },
  },
};
</script>
