<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(data-test="filtersListModal", close)
      template(slot="title")
        span {{ $t('common.filters') }}
      template(slot="text")
        filters-list-form(
          v-model="form.filters",
          :with-pbehavior="config.withPbehavior",
          :with-alarm="config.withAlarm",
          :with-event="config.withEvent",
          :addable="config.hasAccessToAddFilter",
          :editable="config.hasAccessToEditFilter"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import { filtersToForm, formToFilters } from '@/helpers/forms/filter';

import FiltersListForm from '@/components/forms/filters/filters-list-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.filtersList,
  components: { FiltersListForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { filters = [] } = this.modal.config;

    return {
      form: {
        filters: filtersToForm(filters),
      },
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(formToFilters(this.form.filters));
      }

      this.$modals.hide();
    },
  },
};
</script>
