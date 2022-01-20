<template lang="pug">
  v-form(data-test="createFilterModal", @submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        v-text-field(
          data-test="filterTitle",
          v-if="!hiddenFields.includes('title')",
          v-model="form.title",
          v-validate="'required|unique-title'",
          :label="$t('modals.filter.fields.title')",
          :error-messages="errors.collect('title')",
          name="title",
          required
        )
        filter-editor(
          v-if="!hiddenFields.includes('filter')",
          v-model="form.filter",
          :entities-type="entitiesType",
          required
        )
      template(slot="actions")
        v-btn(
          data-test="createFilterCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled || advancedJsonWasChanged",
          :loading="submitting",
          type="submit",
          data-test="createFilterSubmitButton"
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
  },
  created() {
    this.$validator.extend('unique-title', {
      getMessage: () => this.$t('validator.unique'),
      validate: value => (this.initialTitle && this.initialTitle === value)
        || !this.existingTitles.find(title => title === value),
    });
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
