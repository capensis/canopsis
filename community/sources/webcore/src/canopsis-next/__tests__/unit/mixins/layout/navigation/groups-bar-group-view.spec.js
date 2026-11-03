import groupsBarGroupViewMixin from '@/mixins/layout/navigation/groups-bar-group-view';

describe('groups-bar-group-view', () => {
  test.each([
    ['view-id', true],
    ['another-view-id', false],
  ])('Should refresh saved settings only for the active view (%s)', async (activeViewId, shouldRefresh) => {
    const context = {
      view: { _id: 'view-id' },
      activeViewId,
      $t: key => key,
      $modals: { show: vi.fn() },
      updateViewWithPopup: vi.fn().mockResolvedValue(),
      fetchActiveView: vi.fn().mockResolvedValue(),
      fetchAllGroupsListWithWidgetsWithCurrentUser: vi.fn().mockResolvedValue(),
    };

    groupsBarGroupViewMixin.methods.showEditViewModal.call(context);

    const { action } = context.$modals.show.mock.calls[0][0].config;
    const data = { periodic_refresh: { enabled: true, value: 5, unit: 'm' } };

    await action(data);

    expect(context.updateViewWithPopup).toHaveBeenCalledWith({ id: 'view-id', data });
    expect(context.fetchActiveView).toHaveBeenCalledTimes(shouldRefresh ? 1 : 0);
    expect(context.fetchAllGroupsListWithWidgetsWithCurrentUser).toHaveBeenCalledTimes(1);
  });

  test('Should keep the current settings when saving fails', async () => {
    const error = new Error('Saving failed');
    const context = {
      view: { _id: 'view-id' },
      activeViewId: 'view-id',
      $t: key => key,
      $modals: { show: vi.fn() },
      updateViewWithPopup: vi.fn().mockRejectedValue(error),
      fetchActiveView: vi.fn(),
      fetchAllGroupsListWithWidgetsWithCurrentUser: vi.fn(),
    };

    groupsBarGroupViewMixin.methods.showEditViewModal.call(context);

    const { action } = context.$modals.show.mock.calls[0][0].config;

    await expect(action({})).rejects.toThrow(error);

    expect(context.fetchActiveView).not.toHaveBeenCalled();
    expect(context.fetchAllGroupsListWithWidgetsWithCurrentUser).not.toHaveBeenCalled();
  });
});
