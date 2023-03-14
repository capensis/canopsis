<template lang="pug">
  v-form(@submit.stop.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        entity-info-form(
          v-model="form",
          :entity-info="config.entityInfo",
          :infos="config.infos"
        )
      template(#actions="")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.add') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { entityInfoToForm } from '@/helpers/forms/entity-info';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import EntityInfoForm from '@/components/other/entity/form/entity-info-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEntityInfo,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { EntityInfoForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: entityInfoToForm(this.modal.config.entityInfo),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createEntityInfo.create.title');
    },
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
