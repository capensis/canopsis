<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <widget-template-columns-form
          v-if="isColumnsType"
          v-model="form"
        />
        <widget-template-text-form
          v-else
          v-model="form"
          :entity-infos="entityInfos"
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
          type="submit"
          color="primary"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { COLUMNS_WIDGET_TEMPLATES_TYPES, MODALS, VALIDATION_DELAY } from '@/constants';

import { widgetTemplateToForm, formToWidgetTemplate } from '@/helpers/entities/widget/template/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { validationErrorsMixinCreator } from '@/mixins/form';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import WidgetTemplateColumnsForm from '@/components/other/widget-template/form/widget-template-columns-form.vue';
import WidgetTemplateTextForm from '@/components/other/widget-template/form/widget-template-text-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createWidgetTemplate,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    WidgetTemplateColumnsForm,
    WidgetTemplateTextForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    entitiesInfosMixin,
    submittableMixinCreator(),
    validationErrorsMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: widgetTemplateToForm(this.modal.config.widgetTemplate),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createWidgetTemplate.create.title');
    },

    isColumnsType() {
      return COLUMNS_WIDGET_TEMPLATES_TYPES.includes(this.form.type);
    },
  },
  mounted() {
    this.fetchInfos({ withRules: true });
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (!isFormValid) {
        return;
      }

      if (this.config.action) {
        await this.config.action(formToWidgetTemplate(this.form));
      }

      this.$modals.hide();
    },
  },
};
</script>
