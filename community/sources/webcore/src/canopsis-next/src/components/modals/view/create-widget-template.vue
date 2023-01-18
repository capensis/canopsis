<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span TITLE
      template(#text="")
        widget-template-form(v-model="form")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          color="primary"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { widgetTemplateToForm, formToWidgetTemplate } from '@/helpers/forms/widget-template';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { validationErrorsMixinCreator } from '@/mixins/form';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import WidgetTemplateForm from '@/components/other/widget-template/widget-template-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createWidgetTemplate,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: {
    WidgetTemplateForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    validationErrorsMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: widgetTemplateToForm(this.modal.config.widgetTemplate),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (!isFormValid) {
        return;
      }

      if (this.config.action) {
        this.this.config.action(formToWidgetTemplate(this.form));
      }

      this.$modals.hide();
    },
  },
};
</script>
