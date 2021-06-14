<template lang="pug">
  v-form(@submit.stop.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        entity-info-form(
          v-model="form",
          :entity-info="config.entityInfo",
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

import { entityInfoToForm, formToEntityInfo } from '@/helpers/forms/entity-info';

import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';

import EntityInfoForm from '@/components/other/entity/form/entity-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEntityInfo,
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
      form: entityInfoToForm(this.modal.config.entityInfo),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createEntityInfo.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(formToEntityInfo(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>
