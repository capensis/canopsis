<template lang="pug">
  v-card
    v-card-title(primary-title)
      h2 {{ $t('modals.createFilter.title') }}
    v-divider
    v-card-text
      v-text-field(
      :label="$t('modals.createFilter.fields.title')",
      :error-messages="errors.collect('title')"
      v-model="form.title",
      v-validate="'required'",
      required,
      name="title",
      )
      filter-editor(:filter.sync="form.filter", v-model="form.filter")
      v-btn(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';
import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.createFilter,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FilterEditor,
  },
  mixins: [
    modalInnerMixin,
  ],
  data() {
    return {
      form: {
        title: this.modal.config.filter.title,
        filter: this.modal.config.filter.filter,
      },
    };
  },
  methods: {
    async submit() {
      const validationResult = await this.$validator.validate();

      if (validationResult) {
        if (this.config.createFilter) {
          await this.config.createFilter(this.name, this.filterValue);
        }

        this.hideModal();
      }
    },
  },
};
</script>
