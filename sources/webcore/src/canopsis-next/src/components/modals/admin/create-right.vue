<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createRight.title') }}
    v-card-text
      v-form
        v-layout(row)
          v-text-field(
            :label="$t('modals.createRight.fields.id')",
            v-model="form._id",
            data-vv-name="id",
            v-validate="'required'",
            :error-messages="errors.collect('id')"
          )
        v-layout(row)
          v-text-field(
            :label="$t('modals.createRight.fields.description')",
            v-model="form.desc"
          )
        v-layout(row)
          v-select(
            :label="$t('modals.createRight.fields.type')",
            v-model="form.type",
            :items="types"
          )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, USERS_RIGHTS_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesRightMixin from '@/mixins/entities/right';
import { generateRight } from '@/helpers/entities';

export default {
  name: MODALS.createRight,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin, entitiesRightMixin],
  data() {
    return {
      types: [
        { value: '', text: 'Default' },
        { value: USERS_RIGHTS_TYPES.rw, text: USERS_RIGHTS_TYPES.rw },
        { value: USERS_RIGHTS_TYPES.crud, text: USERS_RIGHTS_TYPES.crud },
      ],
      form: {
        _id: '',
        desc: '',
        type: '',
      },
    };
  },
  methods: {
    async submit() {
      try {
        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          const data = { ...generateRight(), ...this.form };

          await this.createRight({ data });

          this.$popups.addSuccess({ text: this.$t('success.default') });
          this.$modals.hide();
        }
        if (this.config.action) {
          await this.config.action();
        }
      } catch (err) {
        this.$popups.addError({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
