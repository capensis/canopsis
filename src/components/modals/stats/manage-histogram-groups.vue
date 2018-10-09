<template lang="pug">
  v-card
    v-card-title.blue.darken-4.white--text
      h2 {{ $t(config.title) }}
    v-container
      v-form(ref="form")
        v-text-field(
        :placeholder="$t('common.title')",
        v-model="form.title", name="title",
        v-validate="'required'",
        required,
        :error-messages="errors.collect('title')"
        )
        v-btn(@click="showFilterModal", small) {{ $t('settings.filterEditor') }}
    v-layout(justify-end)
      v-btn(@click="save").green.darken-4.white--text.mt-3 {{ $t('common.save') }}
</template>

<script>
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

export default {
  name: MODALS.manageHistogramGroups,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        title: '',
        filter: '',
      },
    };
  },
  mounted() {
    this.form = { ...this.config.group };
  },
  methods: {
    showFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: 'modals.filter.create.title',
          filter: this.form.filter,
          action: filter => this.form.filter = filter,
        },
      });
    },
    async save() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }
        this.hideModal();
      }
    },
  },
};
</script>
