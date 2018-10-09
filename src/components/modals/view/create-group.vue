<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.group.title') }}
    v-card-text
      v-text-field(
      :label="$t('modals.group.fields.name')",
      :error-messages="errors.collect('name')"
      v-model="form.name",
      v-validate="'required'",
      name="name",
      )
    v-card-actions
      v-flex(xs6)
        v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
      v-flex.text-xs-right(xs6)
        v-btn.red.darken-4.white--text(@click="remove") {{ $t('common.delete') }}
</template>

<script>
import get from 'lodash/get';

import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import popupMixin from '@/mixins/popup';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

export default {
  name: MODALS.createGroup,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin, popupMixin, entitiesViewGroupMixin],
  data() {
    return {
      form: {
        name: '',
      },
    };
  },
  mounted() {
    this.form.name = get(this.config, 'group.name', '');
  },
  methods: {
    remove() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeGroup({ id: this.modal.config.group._id });
              await this.fetchGroupsList();

              this.hideModal();
            } catch (err) {
              this.addErrorPopup({ text: this.$t('modals.group.errors.isNotEmpty') });
            }
          },
        },
      });
    },
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.updateGroup({
          id: this.modal.config.group._id,
          data: { ...this.form },
        });

        await this.fetchGroupsList();

        this.hideModal();
      }
    },
  },
};
</script>
