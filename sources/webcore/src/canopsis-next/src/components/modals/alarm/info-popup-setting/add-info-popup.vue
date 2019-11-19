<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(data-test="addInfoPopupModal")
      template(slot="title")
        span {{ $t('modals.infoPopupSetting.addInfoPopup.title') }}
      template(slot="text")
        div(data-test="addInfoPopupLayout")
          v-select(
            v-model="form.selectedColumn",
            :items="config.columns",
            item-text="label",
            return-object,
            name="column",
            v-validate="'required'",
            :error-messages="errors.collect('column')"
          )
          text-editor(
            v-model="form.template",
            v-validate="'required'",
            :error-messages="errors.collect('template')",
            name="template"
          )
      template(slot="actions")
        v-btn(
          flat,
          depressed,
          data-test="addInfoCancelButton",
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="errors.any() || submitting",
          type="submit",
          data-test="addInfoSubmitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { find } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import TextEditor from '@/components/other/text-editor/text-editor.vue';

import ModalWrapper from '../../modal-wrapper.vue';

export default {
  name: MODALS.addInfoPopup,
  $_veeValidate: {
    validator: 'new',
  },
  components: { TextEditor, ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      submitting: false,
      form: {
        selectedColumn: {},
        template: '',
      },
    };
  },
  mounted() {
    if (this.config) {
      [this.form.selectedColumn] = this.config.columns;

      if (this.config.popup) {
        const { template, column } = this.config.popup;

        this.form.template = template;
        this.form.selectedColumn = find(this.config.columns, { value: column });
      }
    }
  },
  methods: {
    async submit() {
      try {
        this.submitting = true;

        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          if (this.config.action) {
            await this.config.action({ column: this.form.selectedColumn.value, template: this.form.template });
          }

          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.error({ text: err.description || this.$t('error.default') });
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
