<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        snmp-rule-form(v-model="form")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import { snmpRuleToForm } from '@/helpers/forms/snmp-rule';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import SnmpRuleForm from '@/components/other/snmp-rule/form/snmp-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createSnmpRule,
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
      if (this.config.snmpRule) {
        return this.$t('modals.createSnmpRule.edit.title');
      }

      return this.$t('modals.createSnmpRule.create.title');
    },
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(this.form);
      }

      this.$modals.hide();
    },
  },
};
</script>
