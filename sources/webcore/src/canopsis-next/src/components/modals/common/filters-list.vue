<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(data-test="filtersListModal", close)
      template(slot="title")
        span {{ $t('common.filters') }}
      template(slot="text")
        filters-form(
          v-model="form.filters",
          :entities-type="config.entitiesType",
          :has-access-to-add-filter="config.hasAccessToAddFilter",
          :has-access-to-edit-filter="config.hasAccessToEditFilter"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import { filtersToForm, formToFilters } from '@/helpers/forms/filter';

import FiltersForm from '@/components/other/filter/form/filters-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.filtersList,
  components: { FiltersForm, ModalWrapper },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
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
