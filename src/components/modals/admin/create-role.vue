<template lang="pug">
v-card
  v-card-title.green.darken-4.white--text
    v-layout(justify-space-between, align-center)
      h2 {{ config.title }}
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
        v-layout
          v-btn.mx-0(@click.stop="showViewSelectModal", depressed) Select default view
          div {{ form.defaultView }}
    v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import pick from 'lodash/pick';
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

export default {
  name: MODALS.createRole,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin, entitiesViewGroupMixin],
  data() {
    const group = this.modal.config.group || { name: '', description: '', defaultView: '' };

    return {
      form: pick(group, ['name', 'description', 'defaultview']),
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
    showViewSelectModal() {
      this.showModal({
        name: MODALS.viewSelect,
        config: {
          action: view => this.form.defaultView = view,
        },
      });
    },
  },
};
</script>

