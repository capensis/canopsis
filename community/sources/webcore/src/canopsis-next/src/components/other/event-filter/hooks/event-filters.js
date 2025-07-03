import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useEventFilterStore } from '@/hooks/store/modules/event-filter';

/**
 * Provides modal actions for event filter CRUD operations using the event filter store and popups.
 *
 * @param {Function} [refresh=() => {}] - Callback to refresh the event filter list after an action.
 * @returns {Object} Modal action handlers for event filter CRUD operations.
 * @property {Function} showCreateRuleModal - Show modal to create a new event filter.
 * @property {Function} showDuplicateRuleModal - Show modal to duplicate an event filter.
 * @property {Function} showEditRuleModal - Show modal to edit an event filter.
 * @property {Function} showDeleteRuleModal - Show modal to delete an event filter.
 * @property {Function} showDeleteSelectedRulesModal - Show modal to delete multiple event filters.
 */
export const useEventFilterActions = (refresh = () => {}) => {
  const { t } = useI18n();
  const modals = useModals();
  const popups = usePopups();

  const {
    createEventFilter,
    updateEventFilter,
    removeEventFilter,
  } = useEventFilterStore();

  /**
   * Show modal to create a new event filter.
   *
   * @returns {void}
   */
  const showCreateRuleModal = () => {
    modals.show({
      name: MODALS.createEventFilter,
      config: {
        action: async (data) => {
          await createEventFilter({ data });
          popups.success({ text: t('modals.createEventFilter.create.success') });
          refresh();
        },
      },
    });
  };

  /**
   * Show modal to duplicate an event filter.
   *
   * @param {Object} [rule={}] - The event filter rule to duplicate.
   * @returns {void}
   */
  const showDuplicateRuleModal = (rule = {}) => {
    modals.show({
      name: MODALS.createEventFilter,
      config: {
        title: t('modals.createEventFilter.duplicate.title'),
        rule: omit(rule, ['_id']),
        action: async (data) => {
          await createEventFilter({ data });
          popups.success({ text: t('modals.createEventFilter.duplicate.success') });
          refresh();
        },
      },
    });
  };

  /**
   * Show modal to edit an event filter.
   *
   * @param {Object} [rule={}] - The event filter rule to edit.
   * @returns {void}
   */
  const showEditRuleModal = (rule = {}) => {
    modals.show({
      name: MODALS.createEventFilter,
      config: {
        rule,
        title: t('modals.createEventFilter.edit.title'),
        isDisabledIdField: true,
        action: async (data) => {
          await updateEventFilter({ id: rule._id, data });
          popups.success({ text: t('modals.createEventFilter.edit.success') });
          refresh();
        },
      },
    });
  };

  /**
   * Show modal to delete an event filter.
   *
   * @param {string} id - The ID of the event filter to delete.
   * @returns {void}
   */
  const showDeleteRuleModal = (id) => {
    modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeEventFilter({ id });
          popups.success({ text: t('modals.createEventFilter.remove.success') });
          refresh();
        },
      },
    });
  };

  /**
   * Show modal to delete multiple event filters.
   *
   * @param {Array<Object>} [selected=[]] - Array of selected event filter rules to delete.
   * @returns {void}
   */
  const showDeleteSelectedRulesModal = (selected = []) => {
    modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await Promise.all(selected.map(({ _id: id }) => removeEventFilter({ id })));
          popups.success({ text: t('modals.createEventFilter.remove.success') });
          refresh();
        },
      },
    });
  };

  return {
    showCreateRuleModal,
    showDuplicateRuleModal,
    showEditRuleModal,
    showDeleteRuleModal,
    showDeleteSelectedRulesModal,
  };
};
