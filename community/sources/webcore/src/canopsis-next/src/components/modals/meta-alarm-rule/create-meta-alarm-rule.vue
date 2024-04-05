<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <meta-alarm-rule-form
          v-model="form"
          ref="formElement"
          :active-step.sync="activeStep"
          :disabled-id-field="config.isDisabledIdField"
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          v-if="isLastStep"
          key="submit"
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
        <v-btn
          v-else
          key="next"
          :disabled="!isStepValid"
          type="button"
          class="primary"
          @click="next"
        >
          {{ $t('common.next') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY, META_ALARMS_FORM_STEPS } from '@/constants';

import { formToMetaAlarmRule, metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMetaAlarmRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    MetaAlarmRuleForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      activeStep: META_ALARMS_FORM_STEPS.general,
      isStepValid: false,
      form: metaAlarmRuleToForm(this.modal.config.rule),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.metaAlarmRule.create.title');
    },

    isLastStep() {
      return this.activeStep === META_ALARMS_FORM_STEPS.parameters;
    },
  },
  mounted() {
    this.$watch(() => this.$refs.formElement.isStepValid, (value) => {
      this.isStepValid = value;
    }, { immediate: true });
  },
  methods: {
    async next() {
      const isValid = await this.$refs.formElement.validateStep();

      if (isValid) {
        this.activeStep += 1;
      }
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(formToMetaAlarmRule(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>
