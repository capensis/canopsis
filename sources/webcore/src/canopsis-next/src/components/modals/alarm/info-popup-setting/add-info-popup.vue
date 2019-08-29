<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.infoPopupSetting.addInfoPopup.title') }}
    v-card-text
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
      name="template",
      v-validate="'required'",
      :error-messages="errors.collect('template')"
      )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { find } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import TextEditor from '@/components/other/text-editor/text-editor.vue';

export default {
  name: MODALS.addInfoPopup,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    TextEditor,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
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
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action({ column: this.form.selectedColumn.value, template: this.form.template });
        }
        this.hideModal();
      }
    },
  },
};
</script>
