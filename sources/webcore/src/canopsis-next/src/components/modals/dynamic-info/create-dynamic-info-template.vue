<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        div
          v-text-field(
            v-model="form.title",
            v-validate="'required'",
            :error-messages="errors.collect('title')",
            :label="$t('common.title')",
            name="title"
          )
          h3 {{ $t('modals.createDynamicInfoTemplate.fields.names') }}
          v-layout(
            v-for="(name, index) in form.names",
            :key="name.key",
            row,
            justify-space-between
          )
            v-flex(xs11)
              v-text-field(
                v-model="name.value",
                v-validate="'required'",
                :error-messages="errors.collect(`name[${name.key}]`)",
                :name="`name[${name.key}]`",
                :placeholder="$t('common.name')"
              )
            v-flex(xs1)
              v-btn(
                color="error",
                icon,
                @click="deleteValue(index)"
              )
                v-icon delete
          v-btn.primary.mx-0(@click="showAddValueModal") {{ $t('modals.createDynamicInfoTemplate.buttons.addName') }}
          v-alert(:value="errors.has('names')", type="error")
            span {{ $t('modals.createDynamicInfoTemplate.errors.noNames') }}
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import {
  templateToForm,
  formToTemplate,
  generateTemplateFormName,
} from '@/helpers/forms/dynamic-info-template';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.createDynamicInfoTemplate,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { template = {} } = this.modal.config;

    return {
      form: templateToForm(template),
    };
  },
  computed: {
    title() {
      if (this.config.template) {
        return this.$t('modals.createDynamicInfoTemplate.edit.title');
      }

      return this.$t('modals.createDynamicInfoTemplate.create.title');
    },
  },
  created() {
    this.$validator.attach({
      name: 'names',
      rules: 'required:true',
      getter: () => this.form.names.length > 0,
      context: () => this,
      vm: this,
    });
  },
  methods: {
    showAddValueModal() {
      this.form.names.push(generateTemplateFormName());

      this.$nextTick(() => this.$validator.validate('names'));
    },

    deleteValue(index) {
      this.form.names = this.form.names.filter((v, i) => i !== index);
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToTemplate(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
