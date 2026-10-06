import { CONTEXT_ACTIONS_TYPES, ENTITY_TYPES, MODALS } from '@/constants';

import MassActionsPanel from '@/components/widgets/context/actions/mass-actions-panel.vue';

const $popups = {
  success: vi.fn(),
  info: vi.fn(),
  error: vi.fn(),
};

describe('context mass-actions-panel', () => {
  test('disables Unpause when every selected pause is inherited', () => {
    const resumeAction = {
      type: CONTEXT_ACTIONS_TYPES.resume,
      icon: 'play_arrow',
    };
    const context = {
      pausedEntities: [{ pause: { source: 'parent-id' } }],
      entitiesToResume: [],
      actionsMap: { resume: resumeAction },
      $t: key => key,
    };

    expect(MassActionsPanel.computed.resumeAction.call(context)).toEqual({
      ...resumeAction,
      title: 'context.actions.titles.inheritedPauses',
      disabledButton: true,
    });
  });

  test('bulk pauses only unpaused entities and enables propagation for services and components', async () => {
    const bulkPauseContextEntities = vi.fn();
    const afterSubmit = vi.fn();
    const show = vi.fn();
    const context = {
      entitiesToPause: [
        { _id: 'service-id', type: ENTITY_TYPES.service },
        { _id: 'component-id', type: ENTITY_TYPES.component },
      ],
      $modals: { show },
      $popups,
      $t: key => key,
      bulkPauseContextEntities,
      afterSubmit,
    };

    MassActionsPanel.methods.showPauseModal.call(context);

    const [modal] = show.mock.calls[0];
    expect(modal.name).toBe(MODALS.createEntityPauseEvent);
    expect(modal.config.withInherited).toBe(true);

    await modal.config.action({ comment: 'Comment', type: 'maintenance', inherited: true });

    expect(bulkPauseContextEntities).toHaveBeenCalledWith({
      data: [
        {
          comment: 'Comment',
          type: 'maintenance',
          inherited: true,
          _id: 'service-id',
        },
        {
          comment: 'Comment',
          type: 'maintenance',
          inherited: true,
          _id: 'component-id',
        },
      ],
    });
    expect(afterSubmit).toHaveBeenCalledOnce();
    expect($popups.success).toHaveBeenCalledWith({
      text: 'entity.popups.pauseOperation',
    });
  });

  test('hides propagation for a mixed selection containing a resource', () => {
    const show = vi.fn();
    const context = {
      entitiesToPause: [
        { _id: 'service-id', type: ENTITY_TYPES.service },
        { _id: 'resource-id', type: ENTITY_TYPES.resource },
      ],
      $modals: { show },
    };

    MassActionsPanel.methods.showPauseModal.call(context);

    expect(show.mock.calls[0][0].config.withInherited).toBe(false);
  });

  test('bulk resumes only entities with their own pause', async () => {
    const bulkResumeContextEntities = vi.fn();
    const afterSubmit = vi.fn();
    const context = {
      entitiesToResume: [
        { _id: 'resource-id', pause: { timestamp: 123 } },
        { _id: 'component-id', pause: { timestamp: 456 } },
      ],
      $popups,
      $t: key => key,
      bulkResumeContextEntities,
      afterSubmit,
    };

    await MassActionsPanel.methods.resumeEntityPauses.call(context);

    expect(bulkResumeContextEntities).toHaveBeenCalledWith({
      data: [
        { comment: '', _id: 'resource-id' },
        { comment: '', _id: 'component-id' },
      ],
    });
    expect(afterSubmit).toHaveBeenCalledOnce();
    expect($popups.success).toHaveBeenCalledWith({
      text: 'entity.popups.resumeOperation',
    });
  });
});
