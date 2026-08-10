<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper :close="closeModal">
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-alert v-if="duplicate" type="info">
          {{ t('modals.view.duplicate.infoMessage') }}
        </c-alert>

        <view-form
          :form="form"
          :groups="availableGroups"
          :duplicate-private="isInitialViewPrivate && config.duplicableToAll"
          @update:form="updateForm"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          v-if="submittable"
          :disabled="isDisabled"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ t('common.submit') }}
        </v-btn>

        <v-btn
          v-if="deletable"
          :disabled="submitting"
          variant="outlined"
          color="error"
          @click="remove"
        >
          {{ t('common.delete') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { cloneDeep, find, isEqual, isObject, isString } from 'lodash';
import { computed, onMounted, provide, ref, watch } from 'vue';

import ViewForm from '@/components/other/view/form/view-form.vue';
import { MAX_LIMIT, MODALS } from '@/constants';
import { viewToForm, viewToRequest } from '@/helpers/entities/view/form';
import { uid } from '@/helpers/uid';
import { useAuth } from '@/hooks/auth';
import { useRegisterClickOutsideHandler } from '@/hooks/click-outside';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useView } from '@/hooks/store/modules/view';

/**
 * Modal to create view
 */
export default {
  name: MODALS.createView,

  components: {
    ViewForm,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { fetchCurrentUser } = useAuth();
    const { viewGroups, fetchViewGroupsList, createViewGroup, createPrivateViewGroup } = useView();

    const formContext = useForm();
    const initialForm = viewToForm(props.modal.config?.view);
    const form = ref(initialForm);
    const originalForm = cloneDeep(initialForm);
    const confirmationModalId = uid('modal');
    const submitting = ref(false);

    const config = computed(() => props.modal.config ?? {});
    const title = computed(() => config.value.title || t('modals.view.create.title'));
    const view = computed(() => config.value.view);
    const duplicate = computed(() => !!config.value.duplicate);
    const deletable = computed(() => config.value.deletable);
    const submittable = computed(() => config.value.submittable ?? true);
    const isDisabled = computed(
      () => submitting.value || Object.keys(formContext.errors.value ?? {}).length > 0,
    );
    const isInitialViewPrivate = computed(() => props.modal.config?.view?.is_private ?? false);
    const availableGroups = computed(() => {
      if (config.value.duplicableToAll) {
        return viewGroups.value;
      }

      return viewGroups.value.filter(group => group.is_private === form.value.is_private);
    });

    const hideModal = () => modals.hide({ id: props.modal.id });

    /**
     * Try to find view group by title or create a new one with special title.
     *
     * @param {string} groupTitle
     * @return {ViewGroup | Promise<ViewGroup>}
     */
    const prepareGroup = (groupTitle) => {
      const group = find(viewGroups.value, { title: groupTitle });

      if (group) {
        return group;
      }

      const createFunc = form.value.is_private ? createPrivateViewGroup : createViewGroup;

      return createFunc({ data: { title: groupTitle } });
    };

    const submit = async () => {
      const { valid } = await formContext.validate();

      if (!valid) {
        return;
      }

      try {
        submitting.value = true;

        const group = isString(form.value.group)
          ? await prepareGroup(form.value.group)
          : form.value.group;

        if (config.value.action) {
          await config.value.action(
            viewToRequest({
              ...view.value,
              ...form.value,

              group,
            }),
          );
        }

        hideModal();
      } finally {
        submitting.value = false;
      }
    };

    const confirmClose = () => {
      if (isEqual(form.value, originalForm)) {
        return true;
      }

      if (modals.isModalOpenedById(confirmationModalId)) {
        return false;
      }

      modals.show({
        id: confirmationModalId,
        name: MODALS.clickOutsideConfirmation,
        dialogProps: {
          persistent: true,
        },
        config: {
          action: confirmed => (confirmed ? submit() : hideModal()),
        },
      });

      return false;
    };

    const closeModal = () => {
      if (confirmClose()) {
        hideModal();
      }
    };

    /**
     * Remove view
     */
    const remove = async () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await config.value.remove?.();

              hideModal();
            } catch {
              popups.error({ text: t('modals.view.fail.delete') });
            }
          },
        },
      });
    };

    const updateForm = (nextForm) => {
      form.value = nextForm;
    };

    const fetchAllGroupsListWithWidgets = () =>
      fetchViewGroupsList({
        params: {
          limit: MAX_LIMIT,
          page: 1,
          with_views: true,
          with_tabs: true,
          with_widgets: true,
          with_flags: true,
          with_private: true,
        },
      });

    provide('$closeModal', closeModal);
    useRegisterClickOutsideHandler(confirmClose);

    watch(
      () => form.value.is_private,
      (isPrivate) => {
        if (isObject(form.value.group) && form.value.group.is_private !== isPrivate) {
          form.value.group = undefined;
        }
      },
    );

    watch(
      () => form.value.group,
      (group) => {
        if (isObject(group)) {
          form.value.is_private = group.is_private;
        }
      },
    );

    onMounted(() => Promise.all([fetchAllGroupsListWithWidgets(), fetchCurrentUser()]));

    return {
      t,
      form,
      submitting,
      config,
      title,
      duplicate,
      deletable,
      submittable,
      isDisabled,
      isInitialViewPrivate,
      availableGroups,
      updateForm,
      closeModal,
      remove,
      submit,
    };
  },
};
</script>
