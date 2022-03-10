<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        v-text-field(
          v-if="!hiddenFields.includes('title')",
          v-model="form.title",
          v-validate="titleRules",
          :label="$t('modals.filter.fields.title')",
          :error-messages="errors.collect('title')",
          name="title",
          required
        )
        div.filter-form
          v-expansion-panel
            v-expansion-panel-content(lazy)
              span(slot="header") Add alarm pattern
              div.pa-3 CONTENT
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled || advancedJsonWasChanged",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { get, isString } from 'lodash';

import { ENTITIES_TYPES, MODALS } from '@/constants';

import { filterToForm, formToFilter, filterToObject } from '@/helpers/forms/filter';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createFilter,
  $_veeValidate: {
    validator: 'new',
  },
  components: { FilterEditor, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { title = '', filter = '{}' } = this.modal.config.filter || {};
    const preparedFilter = filterToObject(filter);

    return {
      form: {
        title,
        filter: filterToForm(preparedFilter),
      },
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.filter.create.title');
    },

    entitiesType() {
      return this.config.entitiesType || ENTITIES_TYPES.alarm;
    },

    hiddenFields() {
      return this.config.hiddenFields || [];
    },

    existingTitles() {
      return this.config.existingTitles || [];
    },

    initialTitle() {
      return this.config.filter && this.config.filter.title;
    },

    advancedJsonWasChanged() {
      return get(this.fields, ['advancedJson', 'changed']);
    },

    titleRules() {
      return {
        required: true,
        unique: {
          values: this.existingTitles,
          initialValue: this.initialTitle,
        },
      };
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        if (this.config.action) {
          const preparedFilter = formToFilter(this.form.filter);
          const newFilter = {
            title: this.form.title,
            filter: isString(get(this.config.filter, 'filter', '{}'))
              ? JSON.stringify(preparedFilter)
              : preparedFilter,
          };

          await this.config.action(newFilter);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

<style lang="scss" scoped>
// TODO: move to main with test-suites
.filter-form {
  & /deep/ .v-expansion-panel__header {
    background-color: #979797;
  }

  & /deep/ .v-expansion-panel {
    border-radius: 5px;
    overflow: hidden;
  }
}
</style>
