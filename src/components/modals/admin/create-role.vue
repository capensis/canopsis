<template lang="pug">
v-card
  v-card-title.green.darken-4.white--text
    v-layout(justify-space-between, align-center)
      h2 {{ $t(config.title) }}
      v-btn(@click="hideModal", icon, small)
        v-icon.white--text close
  v-card-text
    v-container
      v-form
        v-text-field(
        v-model="form.name",
        :label="$t('common.name')",
        name="name",
        v-validate="'required'",
        :error-messages="errors.collect('name')"
        )
        v-text-field(v-model="form.description", :label="$t('common.description')")
    v-btn(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import pick from 'lodash/pick';
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.createRole,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    const group = this.modal.config.group || { name: '', description: '' };

    return {
      form: pick(group, ['name', 'description']),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action({ ...this.form });
        }
        this.hideModal();
      }
    },
  },
};
</script>

