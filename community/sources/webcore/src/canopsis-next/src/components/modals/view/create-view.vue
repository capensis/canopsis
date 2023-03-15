<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        v-fade-transition
          v-layout(v-if="pending", justify-center)
            v-progress-circular(color="primary", indeterminate)
          v-layout(v-else)
            v-flex(xs12)
              v-alert(:value="duplicate", type="info") {{ $t('modals.view.duplicate.infoMessage') }}
              view-form(
                v-model="form",
                :groups="groups"
              )
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          v-if="hasUpdateViewAccess",
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
        v-btn(
          v-if="view && hasDeleteViewAccess && !duplicate",
          :disabled="submitting",
          :outline="$system.dark",
          color="error",
          @click="remove"
        ) {{ $t('common.delete') }}
</template>

<script>
import { find, isString } from 'lodash';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { viewToForm, viewToRequest } from '@/helpers/forms/view';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { viewRouterMixin } from '@/mixins/view/router';
import { entitiesViewMixin } from '@/mixins/entities/view';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { permissionsTechnicalViewMixin } from '@/mixins/permissions/technical/view';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ViewForm from '@/components/other/view/view-form.vue';

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
    entitiesViewMixin,
    entitiesViewGroupMixin,
    permissionsTechnicalViewMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      pending: true,
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
      return this.config.duplicate;
    },

    hasUpdateViewAccess() {
      if (this.view && !this.duplicate) {
        return this.checkUpdateAccess(this.view._id) && this.hasUpdateAnyViewAccess;
      }

      return this.hasUpdateAnyViewAccess;
    },

    hasDeleteViewAccess() {
      if (this.view && !this.duplicate) {
        return this.checkDeleteAccess(this.view._id) && this.hasDeleteAnyViewAccess;
      }

      return this.hasDeleteAnyViewAccess;
    },
  },
  async mounted() {
    this.pending = true;

    await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

    this.pending = false;
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
              await this.removeViewWithPopup({ id: this.view._id });
              await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

              this.redirectToHomeIfCurrentRoute();

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
      return find(this.groups, { title }) ?? this.createGroup({ data: { title } });
    },

    /**
     * Convert view form to request object with group and user preference creation if needed
     *
     * @return {Promise}
     */
    async formToRequest() {
      const group = isString(this.form.group)
        ? await this.prepareGroup(this.form.group)
        : this.form.group;

      return viewToRequest({
        ...this.view,
        ...this.form,

        group,
      });
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (!isFormValid) {
        return;
      }

      const data = await this.formToRequest();

      if (this.config.action) {
        await this.config.action(data);
      }

      if (this.duplicate) {
        await this.copyViewWithPopup({ id: this.view._id, data });
      } else if (this.view?._id) {
        await this.updateViewWithPopup({ id: this.view._id, data });
      } else {
        await this.createViewWithPopup({ data });
      }

      await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

      this.$modals.hide();
    },
  },
};
</script>
