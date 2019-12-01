<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        v-text-field(
          v-model="form.name",
          v-validate="'required'",
          :label="$t('modals.group.fields.name')",
          :error-messages="errors.collect('name')",
          name="name",
          data-test="modalGroupNameField"
        )
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="createGroupSubmitButton"
        ) {{ $t('common.submit') }}
        v-btn.error(
          v-if="config.group && hasDeleteAnyViewAccess",
          :disabled="submitting",
          data-test="createGroupDeleteButton",
          @click="remove"
        ) {{ $t('common.delete') }}
</template>

<script>
import { get } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createGroup,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixin(),
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
      this.$modals.show({
        name: this.$constants.MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeGroup({ id: this.config.group._id });
              await this.fetchGroupsList();

              this.$modals.hide();
            } catch (err) {
              this.$popups.error({ text: this.$t('modals.group.errors.isNotEmpty') });
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

        this.$modals.hide();
      }
    },
  },
};
</script>
