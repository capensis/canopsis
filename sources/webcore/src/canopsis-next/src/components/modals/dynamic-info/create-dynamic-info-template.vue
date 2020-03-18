<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span Edit dynamic info template
      template(slot="text")
        div
          v-text-field(
            v-model="form.title",
            v-validate="'required'",
            :error-messages="errors.collect('title')",
            label="Title",
            name="title"
          )
          h3 Names
          v-layout(v-for="(name, index) in form.names", :key="name.key", row, justify-space-between)
            v-flex(xs11)
              v-text-field(
                v-model="name.value",
                v-validate="'required'",
                :error-messages="errors.collect(`name[${name.key}]`)",
                :name="`name[${name.key}]`",
                placeholder="Name"
              )
            v-flex(xs1)
              v-btn(
                color="error",
                icon,
                @click="deleteValue(index)"
              )
                v-icon delete
          v-btn.primary.mx-0(@click="showAddValueModal") Add new value
          v-alert(:value="errors.has('names')", type="error")
            span {{ errors.first('names') }}
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

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
      form: this.templateToForm(template),
    };
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
      this.form.names.push({ key: uid(), value: '' });

      this.$nextTick(() => this.$validator.validate('names'));
    },
    deleteValue(index) {
      this.form.values = this.form.values.filter((v, i) => i !== index);
    },
    formToTemplate({ title = '', names = [] } = {}) {
      return {
        title,
        names: names.map(name => ({ key: uid(), value: name })),
      };
    },
    templateToForm({ title = '', names = [] } = {}) {
      return {
        title,
        names: names.map(({ value }) => value),
      };
    },
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.formToTemplate(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
