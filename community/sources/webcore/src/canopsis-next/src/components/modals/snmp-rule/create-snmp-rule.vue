<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <snmp-rule-form v-model="form" />
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
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
        >
          {{ $t('common.saveChanges') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { snmpRuleToForm, formToSnmpRule } from '@/helpers/entities/snmp-rule/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import SnmpRuleForm from '@/components/other/snmp-rule/form/snmp-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createSnmpRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { SnmpRuleForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { snmpRule } = this.modal.config;

    return {
      form: snmpRuleToForm(snmpRule),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createSnmpRule.create.title');
    },
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(formToSnmpRule(this.form));
      }

      this.$modals.hide();
    },
  },
};
</script>
