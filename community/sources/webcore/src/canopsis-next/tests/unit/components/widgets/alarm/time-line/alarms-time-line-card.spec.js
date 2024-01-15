import flushPromises from 'flush-promises';

import { generateRenderer } from '@unit/utils/vue';

import { ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR, ENTITIES_STATES } from '@/constants';

import AlarmsTimeLineCard from '@/components/widgets/alarm/time-line/alarms-time-line-card.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

const stubs = {
  'c-alarm-chip': true,
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
};

describe('alarms-time-line-card', () => {
  const stateCounterStep = {
    _t: 'statecounter',
    t: 1626159262,
    a: 'root',
    m: 'Idle rule Test all resource',
    val: 3,
    initiator: 'system',
  };
  const stateIncStep = {
    _t: 'stateinc',
    t: 1626159262,
    a: 'root',
    m: 'Idle rule Test all resource',
    val: 3,
    initiator: 'system',
  };
  const stateStepWithStates = {
    _t: 'statecounter',
    t: 1626159262,
    a: 'root',
    m: 'State counter',
    val: Object.entries(ENTITIES_STATES)
      .reduce((acc, [key, state]) => ({
        [`state:${state}`]: key,
        ...acc,
      }), {
        status: 'status',
      }),
    initiator: 'system',
  };
  const statusDecStep = {
    _t: 'statusdec',
    t: 1626159262,
    a: ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR,
    m: 'Status dec',
    role: 'role',
    val: 3,
    initiator: 'system',
  };
  const pbehaviorEnterStep = {
    _t: 'pbhenter',
    t: 1627641985,
    a: 'system',
    m: 'Pbehavior Name pbh. Type: Default pause. Reason: Test reason',
    val: 0,
    initiator: 'external',
  };
  const pbehaviorEnterStepWithHtml = {
    ...pbehaviorEnterStep,
    m: `<p>${pbehaviorEnterStep.m}</p>`,
  };

  const snapshotFactory = generateRenderer(AlarmsTimeLineCard, { stubs });

  it('Renders `alarms-time-line-card` with state counter type', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: stateCounterStep,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-time-line-card` with html as message', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: pbehaviorEnterStepWithHtml,
        isHtmlEnabled: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-time-line-card` without translate', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: pbehaviorEnterStep,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-time-line-card` with state but without translate', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: stateIncStep,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-time-line-card` with state but without translate', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: statusDecStep,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-time-line-card` with states', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: stateStepWithStates,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
