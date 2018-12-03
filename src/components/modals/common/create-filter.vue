<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t(config.title) }}
    v-divider
    v-card-text
      v-text-field(
      :label="$t('modals.filter.fields.title')",
      :error-messages="errors.collect('title')"
      v-model="form.title",
      v-validate="'required'",
      required,
      name="title",
      )
      filter-editor(v-model="form.filter")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import pick from 'lodash/pick';

import { MODALS } from '@/constants';
import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import modalInnerMixin from '@/mixins/modal/inner';

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
    const { filter } = this.modal.config || { title: '', filter: '{}' };

    return {
      form: pick(filter, ['title', 'filter']),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validate();

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
