<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        v-text-field(
          v-model="form.title",
          v-validate="'required'",
          :label="$t('common.title')",
          :error-messages="errors.collect('title')",
          name="title"
        )
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
        v-tooltip.ml-2(
          v-if="group && hasDeleteAnyViewAccess",
          :disabled="group.deletable",
          top
        )
          v-btn.error(
            slot="activator",
            :disabled="submitting || !group.deletable",
            @click="remove"
          ) {{ $t('common.delete') }}
          span {{ $t('modals.group.errors.isNotEmpty') }}
</template>

<script>
import { MODALS } from '@/constants';

import { groupToRequest } from '@/helpers/forms/view';

import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createGroup,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
    entitiesViewGroupMixin,
    permissionsTechnicalViewMixin,
  ],
  data() {
    return {
      form: {
        title: this.modal.config.group.title || '',
      },
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.group.create.title');
    },

    group() {
      return this.config.group;
    },
  },
  methods: {
    remove() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeGroup({ id: this.group._id });
              await this.fetchAllGroupsListWithViewsWithCurrentUser();

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
        const data = groupToRequest({ ...this.group, ...this.form });

        if (this.config.group) {
          await this.updateGroup({ id: this.config.group._id, data });
        } else {
          await this.createGroup({ data });
        }

        await this.fetchCurrentUser();
        await this.fetchAllGroupsListWithViewsWithCurrentUser();

        this.$modals.hide();
      }
    },
  },
};
</script>