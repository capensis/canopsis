<template lang="pug">
  v-card
    v-card-title(primary-title)
      h2 {{ $t('modals.createFilter.title') }}
    v-divider
    v-card-text
      v-text-field(
      :label="$t('modals.createFilter.nameLabel')",
      v-model="name",
      v-validate="'required'",
      required,
      name="name",
      :error-messages="errors.collect('name')"
      )
      filter-editor(:filter.sync="filterValue")
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
      name: '',
      filterValue: '{}',
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
