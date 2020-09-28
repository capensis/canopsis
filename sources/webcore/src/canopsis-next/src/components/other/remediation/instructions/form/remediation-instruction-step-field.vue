<template lang="pug">
  v-layout(row)
    v-flex(xs10)
      v-layout
        v-text-field(
          v-field="step.name",
          v-validate="'required'",
          :label="$t('common.name')",
          :disabled="step.saved",
          :error-messages="nameErrorMessages",
          :name="name",
          box,
          @keyup.prevent.enter="saveName"
        )
      v-layout(v-if="step.saved")
        remediation-instruction-steps-workflow-field(v-field="step.workflow")
      v-layout(v-if="!step.saved", justify-end)
        v-btn.mt-0(depressed, flat, @click="cancelChangeName") {{ $t('common.cancel') }}
        v-btn.mt-0.mr-0.primary(@click="saveName") {{ $t('common.save') }}
    v-flex.mt-3(v-if="step.saved && !hideActions", xs2)
      v-layout(justify-start)
        v-btn.ma-0.ml-2(icon, small, @click="editName")
          v-icon edit
        v-btn.ma-0.ml-1(icon, small, @click.prevent="$emit('remove')")
          v-icon(color="error") delete
</template>

<script>
import formMixin from '@/mixins/form';

import RemediationInstructionStepsWorkflowField from './remediation-instruction-steps-workflow-field.vue';

export default {
  components: { RemediationInstructionStepsWorkflowField },
  mixins: [formMixin],
  inject: ['$validator'],
  model: {
    prop: 'step',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      oldName: null,
    };
  },
  computed: {
    fieldSuffix() {
      return this.step.key ? `-${this.step.key}` : '';
    },

    name() {
      return `name${this.fieldSuffix}`;
    },

    nameErrorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.fieldSuffix, ''));
    },
  },
  methods: {
    async editName() {
      this.oldName = this.step.name;

      this.updateField('saved', false);
    },

    async cancelChangeName() {
      if (this.oldName) {
        this.updateModel({
          ...this.step,
          name: this.oldName,
          saved: true,
        });
      } else {
        this.$emit('remove');
      }
    },

    async saveName() {
      const isValid = await this.$validator.validate(this.name);

      if (isValid) {
        this.oldName = null;

        this.updateField('saved', true);
      }
    },
  },
};
</script>
