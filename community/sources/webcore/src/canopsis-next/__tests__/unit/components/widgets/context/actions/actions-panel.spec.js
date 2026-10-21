import { CONTEXT_ACTIONS_TYPES, ENTITY_TYPES, MODALS } from '@/constants';

import ActionsPanel from '@/components/widgets/context/actions/actions-panel.vue';

const translate = (key, parameters = {}) => {
  if (key.startsWith('entity.types.')) {
    return key.split('.').at(-1).replace(/^./u, character => character.toUpperCase());
  }

  if (parameters.message) {
    return `${key}:${parameters.message}`;
  }

  return `${key}:`;
};
const $popups = {
  success: vi.fn(),
  info: vi.fn(),
  error: vi.fn(),
};

const createActionsContext = item => ({
  item,
  $t: translate,
  $popups,
  showEditEntityModal: vi.fn(),
  showDuplicateServiceModal: vi.fn(),
  showDeleteEntityModal: vi.fn(),
  showAddPbehaviorModal: vi.fn(),
  showPauseModal: vi.fn(),
  resumeEntityPause: vi.fn(),
  showVariablesHelpModal: vi.fn(),
});

const prepareActionsContext = (item) => {
  const context = createActionsContext(item);
  context.pauseAction = ActionsPanel.computed.pauseAction.call(context);
  context.actionsMap = ActionsPanel.computed.actionsMap.call(context);
  context.filteredActionsMap = context.actionsMap;

  return context;
};

describe('context actions-panel', () => {
  test('uses the order from the specification and shows Pause for an active entity', () => {
    const context = prepareActionsContext({
      _id: 'component-id',
      type: ENTITY_TYPES.component,
      deletable: true,
    });

    expect(ActionsPanel.computed.actions.call(context).map(({ type }) => type)).toEqual([
      CONTEXT_ACTIONS_TYPES.editEntity,
      CONTEXT_ACTIONS_TYPES.pause,
      CONTEXT_ACTIONS_TYPES.pbehaviorAdd,
      CONTEXT_ACTIONS_TYPES.variablesHelp,
      CONTEXT_ACTIONS_TYPES.deleteEntity,
    ]);
  });

  test('shows a disabled Unpause with the inherited pause source', () => {
    const context = createActionsContext({
      _id: 'resource-id',
      type: ENTITY_TYPES.resource,
      pause: {
        inherited: true,
        source: 'parent-component-id',
        message: '\nParent component',
      },
    });

    expect(ActionsPanel.computed.pauseAction.call(context)).toEqual(
      expect.objectContaining({
        type: CONTEXT_ACTIONS_TYPES.resume,
        icon: 'play_arrow',
        title: 'context.actions.titles.inheritedPause:Parent component',
        disabledButton: true,
      }),
    );
  });

  test('keeps Unpause enabled for an own pause propagated to dependencies', () => {
    const context = createActionsContext({
      _id: 'component-id',
      type: ENTITY_TYPES.component,
      pause: {
        inherited: true,
        message: 'Comment',
      },
    });

    expect(ActionsPanel.computed.pauseAction.call(context)).toEqual(
      expect.objectContaining({
        type: CONTEXT_ACTIONS_TYPES.resume,
        title: 'context.actions.titles.unpause:',
        disabledButton: false,
      }),
    );
  });

  test('opens Pause modal with propagation only for supported entity types', async () => {
    const pauseContextEntity = vi.fn();
    const fetchContextEntitiesListWithPreviousParams = vi.fn();
    const show = vi.fn();
    const context = {
      item: { _id: 'service-id', type: ENTITY_TYPES.service },
      $modals: { show },
      $popups,
      $t: translate,
      pauseContextEntity,
      fetchContextEntitiesListWithPreviousParams,
    };

    ActionsPanel.methods.showPauseModal.call(context);

    const [modal] = show.mock.calls[0];
    expect(modal.name).toBe(MODALS.createEntityPauseEvent);
    expect(modal.config.withInherited).toBe(true);

    await modal.config.action({ comment: 'Comment', type: 'maintenance', inherited: true });

    expect(pauseContextEntity).toHaveBeenCalledWith({
      data: {
        comment: 'Comment',
        type: 'maintenance',
        inherited: true,
        _id: 'service-id',
      },
    });
    expect(fetchContextEntitiesListWithPreviousParams).toHaveBeenCalledOnce();
    expect($popups.success).toHaveBeenCalledWith({
      text: 'entity.popups.pauseOperation:',
    });
  });

  test('resumes an own pause directly', async () => {
    const resumeContextEntity = vi.fn();
    const fetchContextEntitiesListWithPreviousParams = vi.fn();
    const context = {
      item: { _id: 'resource-id', pause: { timestamp: 123 } },
      $popups,
      $t: translate,
      resumeContextEntity,
      fetchContextEntitiesListWithPreviousParams,
    };

    await ActionsPanel.methods.resumeEntityPause.call(context);

    expect(resumeContextEntity).toHaveBeenCalledWith({
      data: { comment: '', _id: 'resource-id' },
    });
    expect(fetchContextEntitiesListWithPreviousParams).toHaveBeenCalledOnce();
    expect($popups.success).toHaveBeenCalledWith({
      text: 'entity.popups.resumeOperation:',
    });
  });

  test('handles a resume 404 and displays an error popup', async () => {
    const resumeContextEntity = vi.fn().mockRejectedValue({ error: 'Not Found' });
    const context = {
      item: { _id: 'resource-id', pause: { timestamp: 123 } },
      $popups,
      $t: translate,
      resumeContextEntity,
      fetchContextEntitiesListWithPreviousParams: vi.fn(),
    };

    await ActionsPanel.methods.resumeEntityPause.call(context);

    expect($popups.error).toHaveBeenCalledWith({
      text: 'alarm.popups.unpauseFailed:',
    });
  });
});
