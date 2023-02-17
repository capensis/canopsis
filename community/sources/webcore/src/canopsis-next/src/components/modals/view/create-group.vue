<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        v-text-field(
          v-model="form.title",
          v-validate="'required'",
          :label="$t('common.title')",
          :error-messages="errors.collect('title')",
          name="title"
        )
      template(#actions="")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
        v-tooltip(
          v-if="group && hasDeleteAnyViewAccess",
          :disabled="group.deletable",
          top
        )
          template(#activator="{ on }")
            div.ml-2(v-on="on")
              v-btn.error(
                :disabled="submitting || !group.deletable",
                :outline="$system.dark",
                color="error",
                @click="remove"
              ) {{ $t('common.delete') }}
          span {{ $t('modals.group.errors.isNotEmpty') }}
</template>

<script>
import { MODALS } from '@/constants';

import { groupToRequest } from '@/helpers/forms/view';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createGroup,
  $_veeValidate: {
    validator: 'new',
  },
  inject: ['$system'],
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
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
              await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

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
        await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

        this.$modals.hide();
      }
    },
  },
};
</script>
