<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t(config.title) }}
    v-container
      v-form(ref="form")
        v-text-field(
        :placeholder="$t('common.title')",
        v-model="form.title", name="title",
        v-validate="'required'",
        required,
        :error-messages="errors.collect('title')"
        )
        v-btn(@click="showFilterModal") {{ $t('settings.filterEditor') }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="save") {{ $t('common.save') }}
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
