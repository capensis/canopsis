<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
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
          :entitiesType="entitiesType",
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
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="createFilterSubmitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { ENTITIES_TYPES, MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createFilter,
  $_veeValidate: {
    validator: 'new',
  },
  components: { FilterEditor, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { hiddenFields = [], filter = {}, entitiesType = ENTITIES_TYPES.alarm } = this.modal.config;

    return {
      hiddenFields,
      entitiesType,

      form: pick(filter, ['title', 'filter']),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.filter.create.title');
    },
    existingTitles() {
      return this.config.existingTitles || [];
    },
    initialTitle() {
      return this.config.filter && this.config.filter.title;
    },
  },
  created() {
    this.$validator.extend('unique-title', {
      getMessage: () => this.$t('validator.unique'),
      validate: value => (this.initialTitle && this.initialTitle === value) ||
        !this.existingTitles.find(title => title === value),
    });
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
