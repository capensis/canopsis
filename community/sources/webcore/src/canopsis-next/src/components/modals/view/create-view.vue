<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        c-alert(:value="duplicate", type="info") {{ $t('modals.view.duplicate.infoMessage') }}
        view-form(v-model="form", :groups="groups", :duplicate-private="isInitialViewPrivate && duplicate")
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          v-if="submittable",
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
        v-btn(
          v-if="deletable",
          :disabled="submitting",
          :outline="$system.dark",
          color="error",
          @click="remove"
        ) {{ $t('common.delete') }}
</template>

<script>
import { find, isString } from 'lodash';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { viewToForm, viewToRequest } from '@/helpers/entities/view/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { viewRouterMixin } from '@/mixins/view/router';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ViewForm from '@/components/other/view/form/view-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createView,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: { ViewForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    viewRouterMixin,
    entitiesViewGroupMixin,
    permissionsTechnicalViewMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: viewToForm(this.modal.config.view),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.view.create.title');
    },

    view() {
      return this.config.view;
    },

    duplicate() {
      return !!this.config.duplicate;
    },

    deletable() {
      return this.config.deletable;
    },

    submittable() {
      return this.config.submittable ?? true;
    },

    isInitialViewPrivate() {
      return this.modal.config.view?.is_private ?? false;
    },
  },
  async mounted() {
    this.fetchAllGroupsListWithWidgetsWithCurrentUser();
  },
  methods: {
    /**
     * Remove view
     */
    async remove() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.config.remove?.();

              this.$modals.hide();
            } catch (err) {
              this.$popups.error({ text: this.$t('modals.view.fail.delete') });
            }
          },
        },
      });
    },

    /**
     * Try to find view group by title or create a new one with special title
     *
     * @param {string} title
     * @return {ViewGroup | Promise<ViewGroup>}
     */
    prepareGroup(title) {
      const group = find(this.groups, { title });

      if (group) {
        return group;
      }

      const createFunc = this.form.is_private ? this.createPrivateGroup : this.createGroup;

      return createFunc({ data: { title } });
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (!isFormValid) {
        return;
      }

      const group = isString(this.form.group)
        ? await this.prepareGroup(this.form.group)
        : this.form.group;

      if (this.config.action) {
        await this.config.action(viewToRequest({
          ...this.view,
          ...this.form,

          group,
        }));
      }

      this.$modals.hide();
    },
  },
};
</script>
