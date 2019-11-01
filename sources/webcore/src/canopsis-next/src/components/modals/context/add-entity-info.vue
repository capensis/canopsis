<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-card-text
      entity-info-form(
        v-model="form",
        :entityInfo="config.editingInfo",
        :infos="config.infos"
      )
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn(@click="submit", color="primary") {{ $t('common.add') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import EntityInfoForm from '@/components/other/context/entity-info-form.vue';

export default {
  name: MODALS.addEntityInfo,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EntityInfoForm },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        name: '',
        description: '',
        value: '',
      },
    };
  },
  mounted() {
    if (this.config.editingInfo) {
      this.form = { ...this.config.editingInfo };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(this.form);
        this.$modals.hide();
      }
    },
  },
};
</script>
