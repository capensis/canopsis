import { fakeAlarmDetails } from '@unit/data/alarm';
import { mockPopups } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { API_HOST, API_ROUTES } from '@/config';
import { ALARM_STEPS_TYPES, EXPORT_STATUSES } from '@/constants';

import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmTimeline from '@/components/widgets/alarm/timeline/alarm-timeline.vue';

const alarm = vi.hoisted(() => ({
  createAlarmTimelineExport: vi.fn(),
  fetchAlarmTimelineExport: vi.fn(),
}));
const permissions = vi.hoisted(() => ({ hasExportAccess: true }));
const openUrlInNewTab = vi.hoisted(() => vi.fn());

vi.mock('@/hooks/store/modules/alarm', () => ({
  useAlarm: () => alarm,
}));

vi.mock('@/hooks/auth', () => ({
  useCanPermission: () => ({ hasAccess: permissions.hasExportAccess }),
}));

vi.mock('@/helpers/url', async importOriginal => ({
  ...(await importOriginal()),
  openUrlInNewTab,
}));

const stubs = {
  'c-enabled-field': true,
  'alarm-timeline-days': true,
  'c-pagination': true,
  'v-btn': true,
  'v-icon': true,
  'v-spacer': true,
};

const selectGroupField = wrapper => wrapper.findAll('c-enabled-field-stub').at(0);
const selectOnlyCommentsField = wrapper => wrapper.findAll('c-enabled-field-stub').at(1);
const selectPagination = wrapper => wrapper.find('c-pagination-stub');
const selectExportButton = wrapper => wrapper.find('v-btn-stub');

describe('alarm-timeline', () => {
  const $popups = mockPopups();
  const {
    steps: { data: steps, meta },
  } = fakeAlarmDetails();
  const factory = generateShallowRenderer(AlarmTimeline, {
    stubs,
    parentComponent: {
      provide: { $popups },
    },
  });

  beforeEach(() => {
    permissions.hasExportAccess = true;
    alarm.createAlarmTimelineExport.mockReset();
    alarm.fetchAlarmTimelineExport.mockReset();
    openUrlInNewTab.mockReset();
    $popups.error.mockReset();
  });

  test('User can update group value and emit update event', () => {
    const newGroup = true;
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query: { group: false },
      },
    });

    const groupField = selectGroupField(wrapper);

    groupField.triggerCustomEvent('input', newGroup);

    expect(wrapper).toEmit('update:query', { group: newGroup, page: 1 });
  });

  test('User can enable only comments filter and emit update event', () => {
    const query = { group: true, page: 2 };
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query,
      },
    });

    const onlyCommentsField = selectOnlyCommentsField(wrapper);

    onlyCommentsField.triggerCustomEvent('input', true);

    expect(wrapper).toEmit('update:query', {
      group: true,
      page: 1,
      type: ALARM_STEPS_TYPES.comment,
    });
  });

  test('User can update page value and emit update event', () => {
    const query = { group: true };
    const newPage = 2;
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query,
      },
    });

    const pagination = selectPagination(wrapper);

    pagination.triggerCustomEvent('update:page', newPage);

    expect(wrapper).toEmit('update:query', { ...query, page: newPage });
  });

  test('Days computed property returns expected value', () => {
    const wrapper = factory({
      propsData: {
        steps,
        meta,
        query: {},
      },
    });

    expect(wrapper.vm.days).toEqual(groupAlarmSteps(steps));
  });

  test('Timeline export is hidden without permission', () => {
    permissions.hasExportAccess = false;

    const wrapper = factory({
      propsData: {
        alarmId: 'alarm-id',
        steps,
        meta,
        query: {},
      },
    });

    expect(selectExportButton(wrapper).exists()).toBeFalsy();
  });

  test('Timeline is exported after clicking export button', async () => {
    const alarmId = 'alarm-id';
    const exportId = 'export-id';
    const exportData = { _id: exportId, status: EXPORT_STATUSES.completed };

    alarm.createAlarmTimelineExport.mockResolvedValue(exportData);
    alarm.fetchAlarmTimelineExport.mockResolvedValue(exportData);

    const wrapper = factory({
      propsData: {
        alarmId,
        steps,
        meta,
        query: {},
      },
    });

    selectExportButton(wrapper).triggerCustomEvent('click');

    await flushPromises();

    expect(alarm.createAlarmTimelineExport).toHaveBeenCalledWith({
      data: { alarm: alarmId },
    });
    expect(alarm.fetchAlarmTimelineExport).toHaveBeenCalledWith({ id: exportId });
    expect(openUrlInNewTab).toHaveBeenCalledWith(
      `${API_HOST}${API_ROUTES.alarmTimelineExport}/${exportId}/download`,
    );
  });

  test('Error popup is shown when timeline export fails', async () => {
    const exportData = { _id: 'export-id', status: EXPORT_STATUSES.failed };

    alarm.createAlarmTimelineExport.mockResolvedValue(exportData);
    alarm.fetchAlarmTimelineExport.mockResolvedValue(exportData);

    const wrapper = factory({
      propsData: {
        alarmId: 'alarm-id',
        steps,
        meta,
        query: {},
      },
    });

    selectExportButton(wrapper).triggerCustomEvent('click');

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: 'Failed to export alarm timeline in CSV format',
    });
    expect(openUrlInNewTab).not.toHaveBeenCalled();
  });
});
