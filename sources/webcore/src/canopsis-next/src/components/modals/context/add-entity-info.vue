<template lang="pug">
  v-form(@submit.stop.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        entity-info-form(
          v-model="form",
          :entityInfo="config.editingInfo",
          :infos="config.infos"
        )
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.add') }}
</template>

<script>
import { MODALS } from '@/constants';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import EntityInfoForm from '@/components/widgets/context/form/entity-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.addEntityInfo,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EntityInfoForm, ModalWrapper },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: {
        name: '',
        description: '',
        value: '',
      },
    };
  },
  mounted() {
    if (this.config.editingInfo) {
      this.form = { ...this.config.editingInfo };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(this.form);
        this.$modals.hide();
      }
    },
  },
};
</script>
