<template lang="pug">
  v-card(data-test="createGroupViewModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      v-text-field(
      data-test="modalGroupNameField",
      :label="$t('modals.group.fields.name')",
      :error-messages="errors.collect('name')"
      v-model="form.name",
      v-validate="'required'",
      name="name",
      )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(
      @click="submit",
      data-test="createGroupSubmitButton"
      ) {{ $t('common.submit') }}
      v-btn.error(
      v-if="config.group && hasDeleteAnyViewAccess",
      @click="remove",
      data-test="createGroupDeleteButton"
      ) {{ $t('common.delete') }}
</template>

<script>
import { get } from 'lodash';

import { MODALS } from '@/constants';

import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  name: MODALS.createGroup,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    popupMixin,
    modalInnerMixin,
    entitiesViewGroupMixin,
    rightsTechnicalViewMixin,
  ],
  data() {
    return {
      form: {
        name: '',
      },
    };
  },
  computed: {
    title() {
      if (this.config.group) {
        return this.$t('modals.group.edit.title');
      }

      return this.$t('modals.group.create.title');
    },
  },
  mounted() {
    this.form.name = get(this.config, 'group.name', '');
  },
  methods: {
    remove() {
      this.showModal({
        name: this.$constants.MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeGroup({ id: this.config.group._id });
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
        if (this.config.group) {
          await this.updateGroup({
            id: this.config.group._id,
            data: { ...this.form },
          });
        } else {
          await this.createGroup({ data: { ...this.form } });
        }

        await this.fetchGroupsList();

        this.hideModal();
      }
    },
  },
};
</script>
