import { isEqual } from 'lodash';
import { computed, nextTick, ref, unref, watch } from 'vue';

import { MODALS } from '@/constants';
import { formToWidget, widgetToForm } from '@/helpers/entities/widget/form';
import { uid } from '@/helpers/uid';
import { useForm } from '@/hooks/form';
import { useModals } from '@/hooks/modals';
import { useSidebar } from '@/hooks/sidebar';
import { useActiveView } from '@/hooks/store/modules/active-view';
import { useUserPreference } from '@/hooks/store/modules/user-preference';
import { useWidget } from '@/hooks/store/modules/widget';

/**
 * Hook for widget settings sidebars.
 *
 * Keeps the old widget settings mixin behavior available for Composition API sidebars:
 * local widget form state, dirty-close confirmation, submit, user preference sync and
 * active view refresh.
 *
 * @param {Object} options
 * @param {Object|import('vue').Ref<Object>} options.sidebar Sidebar store item.
 * @param {boolean|import('vue').Ref<boolean>} [options.localWidget=false] Whether preference updates are local only.
 * @returns {Object}
 */
export const useWidgetSettings = ({ sidebar, localWidget = false } = {}) => {
  const modals = useModals();
  const sidebarApi = useSidebar();
  const { fetchActiveView } = useActiveView();
  const { createWidget, updateWidget } = useWidget();
  const {
    getUserPreferenceByWidgetId,
    fetchUserPreference,
    updateUserPreference,
    updateLocalUserPreference,
  } = useUserPreference();

  const sidebarModel = computed(() => unref(sidebar) ?? {});
  const validationForm = useForm({
    initialValues: widgetToForm(sidebarModel.value.config?.widget),
  });

  const form = validationForm.values;
  const originalForm = ref(widgetToForm(sidebarModel.value.config?.widget));
  const hasChanges = ref(false);
  const submitting = ref(false);
  const confirmationModalId = uid('modal');

  const config = computed(() => sidebarModel.value.config ?? {});
  const widget = computed(() => config.value.widget ?? {});
  const duplicate = computed(() => config.value.duplicate);
  const userPreference = computed(() => unref(getUserPreferenceByWidgetId)(widget.value._id));

  const hideSettingsSidebar = () => sidebarApi.hide({ id: sidebarModel.value.id });

  const scrollToFirstError = () => {
    nextTick(() => {
      document.querySelector('.sidebar-wrapper .v-input--error')?.scrollIntoView({
        block: 'center',
        behavior: 'smooth',
      });
    });
  };

  const updateContentInUserPreference = (content = {}) => {
    const method = unref(localWidget) ? updateLocalUserPreference : updateUserPreference;

    return method({
      data: {
        ...userPreference.value,

        content: {
          ...userPreference.value.content,
          ...content,
        },
      },
    });
  };

  const fetchWidgetUserPreference = data =>
    unref(localWidget) ? Promise.resolve() : fetchUserPreference(data);

  /**
   * Submit settings form.
   */
  const submit = async () => {
    await nextTick();

    const { valid: isFormValid } = await validationForm.validate();

    if (!isFormValid) {
      scrollToFirstError();

      return;
    }

    try {
      submitting.value = true;

      const { _id: widgetId, tab: tabId } = widget.value;
      const data = formToWidget(form.value);

      data.tab = tabId;

      if (duplicate.value || !widgetId) {
        await createWidget({ data });
      } else {
        await updateWidget({ id: widgetId, data });
      }

      if (
        data.parameters.mainFilter &&
        userPreference.value.content.mainFilter === data.parameters.mainFilter
      ) {
        await updateContentInUserPreference({ mainFilter: null });
      }

      if (widgetId) {
        await fetchWidgetUserPreference({ id: widgetId });
      }

      await fetchActiveView();

      hideSettingsSidebar();
    } finally {
      submitting.value = false;
    }
  };

  const close = () => {
    if (!hasChanges.value || isEqual(form.value, originalForm.value)) {
      hideSettingsSidebar();

      return;
    }

    if (modals.isModalOpenedById(confirmationModalId)) {
      return;
    }

    modals.show({
      id: confirmationModalId,
      name: MODALS.clickOutsideConfirmation,
      dialogProps: {
        persistent: true,
      },
      config: {
        action: confirmed => (confirmed ? submit() : hideSettingsSidebar()),
      },
    });
  };

  const stopFormWatcher = watch(
    form,
    () => {
      hasChanges.value = true;

      nextTick(() => stopFormWatcher());
    },
    { deep: true },
  );

  return {
    form,
    originalForm,
    hasChanges,
    submitting,
    sidebar: sidebarModel,
    config,
    widget,
    duplicate,
    userPreference,
    validationForm,

    close,
    submit,
    hideSettingsSidebar,
    scrollToFirstError,
    updateContentInUserPreference,
    fetchWidgetUserPreference,
  };
};
