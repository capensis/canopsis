<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <alarm-status-rule-form
          v-model="form"
          :flapping="config.flapping"
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
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { alarmStatusRuleToForm, formToAlarmStatusRule } from '@/helpers/entities/alarm-status-rule/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmStatusRuleForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createAlarmStatusRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { AlarmStatusRuleForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { rule, flapping } = this.modal.config;

    return {
      form: alarmStatusRuleToForm(rule, flapping),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToAlarmStatusRule(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
